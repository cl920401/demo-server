package kafka

import (
	"demo-server/lib/apm"
	"demo-server/lib/config"
	"demo-server/lib/errors"
	"demo-server/lib/log"
	"fmt"
	"github.com/Shopify/sarama"
	"sort"
	"strings"
)

//TopicOffset kafka 每个分区的offset信息
type TopicOffset struct {
	PartitionIndex int32 `json:"index"`
	Offset         int64 `json:"offset"`
}

func kafkaConf() *sarama.Config {
	conf := sarama.NewConfig()
	conf.Version = Version //kafka版本号
	conf.MetricRegistry = apm.GetRegistry()

	if sasl := config.Get("kafka.sasl.enable").Bool(false); sasl {
		conf.Net.SASL.Enable = true
		conf.Net.SASL.User = config.Get("kafka.sasl.user").String("")
		conf.Net.SASL.Password = config.Get("kafka.sasl.password").String("")
	}
	return conf
}

//NewAdmin 创建一个kafka admin clinet
func NewAdmin(addrs ...string) (sarama.ClusterAdmin, error) {
	if len(addrs) == 0 {
		addrs = strings.Split(config.Get("kafka.addrs").String(""), ",")
	}

	admin, err := sarama.NewClusterAdmin(addrs, kafkaConf())
	if err != nil {
		err := errors.New(ErrKafkaClusterCreate,
			fmt.Sprintf("error in create new cluster addrs %s, %s ", strings.Join(addrs, ","), err.Error()))
		log.Error(err)
		return nil, err
	}
	return admin, nil
}

func NewClient(addrs ...string) (sarama.Client, error) {
	if len(addrs) == 0 {
		addrs = strings.Split(config.Get("kafka.addrs").String(""), ",")
	}

	client, err := sarama.NewClient(addrs, kafkaConf())
	if err != nil {
		err := errors.New(ErrKafkaClusterCreate,
			fmt.Sprintf("error in create new client addrs %s, %s ", strings.Join(addrs, ","), err.Error()))
		log.Error(err)
		return nil, err
	}
	return client, nil
}

//GetOffset 获取指定topic的offset
func GetOffset(topic string, partitions []int32) ([]TopicOffset, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	to := make([]TopicOffset, 0, len(partitions))
	for _, v := range partitions {
		offset, err := client.GetOffset(topic, v, sarama.OffsetNewest)
		if err != nil {
			log.Error(err)
			continue
		}
		to = append(to, TopicOffset{
			PartitionIndex: v,
			Offset:         offset,
		})
	}
	sort.Slice(to, func(i, j int) bool {
		return to[i].PartitionIndex < to[j].PartitionIndex
	})
	return to, nil
}

//CreateTopic 创建topic
func CreateTopic(name string, partitions int32, replication int16) error {
	admin, err := NewAdmin()
	if err != nil {
		return err
	}
	defer admin.Close()

	err = admin.CreateTopic(name, &sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replication,
	}, false)
	if err != nil {
		err := errors.New(ErrKafkaCreateTopic, name+": "+err.Error())
		log.Error(err)
		return err
	}
	log.Debugf("create topic success: %s", name)
	return nil
}

//GetOffsetByGroup 获取某个group 在指定topic下的消费情况
func GetOffsetByGroup(group string, topicPartitions map[string][]int32) (map[string][]TopicOffset, error) {
	admin, err := NewAdmin()
	if err != nil {
		return nil, err
	}
	defer admin.Close()

	offset, err := admin.ListConsumerGroupOffsets(group, topicPartitions)
	if err != nil {
		return nil, err
	}
	var data = map[string][]TopicOffset{}
	for k, partitions := range topicPartitions {

		to := make([]TopicOffset, 0, len(partitions))

		for _, v := range partitions {
			offsetBlock := offset.GetBlock(k, v)
			if offsetBlock == nil {
				continue
			}
			to = append(to, TopicOffset{
				PartitionIndex: v,
				Offset:         offsetBlock.Offset,
			})
		}
		sort.Slice(to, func(i, j int) bool {
			return to[i].PartitionIndex < to[j].PartitionIndex
		})
		data[k] = to
	}
	log.Debug(data)
	return data, nil
}

//DeleteTopic 删除topic
func DeleteTopic(name string) error {
	admin, err := NewAdmin()
	if err != nil {
		return err
	}
	defer admin.Close()

	err = admin.DeleteTopic(name)
	if err != nil {
		err := errors.New(ErrKafkaDeleteTopic, err.Error())
		log.Error(err)
		return err
	}
	log.Debugf("delete topic success: %s", name)
	return nil
}

func DeleteGroup(name string) error {
	admin, err := NewAdmin()
	if err != nil {
		return err
	}
	defer admin.Close()

	err = admin.DeleteConsumerGroup(name)
	if err != nil {
		err := errors.New(ErrKafkaDeleteGroup, err.Error())
		log.Error(err)
		return err
	}
	log.Debugf("delete group success: %s", name)
	return nil
}

func GetTopicsAndGroups() (map[string]sarama.TopicDetail, map[string]string, error) {
	admin, err := NewAdmin()
	if err != nil {
		return nil, nil, err
	}
	topics, err := admin.ListTopics()
	if err != nil {
		return nil, nil, err
	}
	groups, err := admin.ListConsumerGroups()
	if err != nil {
		return nil, nil, err
	}
	return topics, groups, nil
}

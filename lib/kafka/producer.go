package kafka

import (
	"demo-server/lib/errors"
	"demo-server/lib/log"
	"github.com/Shopify/sarama"
	"time"
)

//SyncProducer kafka 同步消息发送
type SyncProducer struct {
	Address []string
	Version sarama.KafkaVersion
	Timeout int
	stop    chan bool
}

//AsyncProducer kafka异步消息发送
type AsyncProducer struct {
	Address []string
	Version sarama.KafkaVersion
	Timeout int
	stop    chan bool
}

//SyncProducer.Stop 停止kafka发送数据
func (v *SyncProducer) Stop() {
	v.stop <- true
}

//AsyncProducer.Stop 停止kafka发送数据
func (v *AsyncProducer) Stop() {
	v.stop <- true
}

//SyncProducer.Start 同步消息模式 发送消息
func (v *SyncProducer) Start(topic string, data <-chan string, f func(topic string, partition int32, offset int64)) error {
	v.stop = make(chan bool)
	conf := kafkaConf()
	if v.Version != (sarama.KafkaVersion{}) {
		conf.Version = v.Version //kafka版本号
	} else {
		conf.Version = Version //kafka版本号
	}
	conf.Producer.Return.Successes = true
	conf.Producer.Timeout = time.Duration(v.Timeout) * time.Second
	producer, err := sarama.NewSyncProducer(v.Address, conf)
	if err != nil {
		err := errors.New(ErrKafkaSyncProducer, err.Error())
		log.Error(err)
		return err
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Error(err)
		}
	}()
	for {
		select {
		case <-v.stop:
			log.Println("kafka SyncProducer is stop")
			return nil
		case value := <-data:
			if value != "" {
				msg := &sarama.ProducerMessage{
					Topic: topic,
					Value: sarama.StringEncoder(value),
				}
				part, offset, err := producer.SendMessage(msg)
				if f != nil {
					f(topic, part, offset)
				}
				if err != nil {
					log.Errorf("send message(%s) err=%s ", value, err)
				} else {
					log.Debugf("[%s] "+value+" send success，partition=%d, offset=%d ", topic, part, offset)
				}
			}
		}
	}
}

//AsyncProducer.Start 异步消息
func (v *AsyncProducer) Start(topic string, data <-chan string, f func(topic string, partition int32, offset int64)) error {
	v.stop = make(chan bool)
	config := kafkaConf()
	config.Producer.Timeout = time.Duration(v.Timeout) * time.Second
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	if v.Version != (sarama.KafkaVersion{}) {
		config.Version = v.Version
	} else {
		config.Version = Version
	}
	//使用配置,新建一个异步生产者
	producer, err := sarama.NewAsyncProducer(v.Address, config)
	if err != nil {
		log.Errorf("sarama.NewAsyncProducer err, message=%s ", err)
		return err
	}
	defer producer.AsyncClose()
	////循环判断哪个通道发送过来数据.
	go func(p sarama.AsyncProducer) {
		for {
			select {
			case suc := <-p.Successes():
				if f != nil {
					f(topic, suc.Partition, suc.Offset)
				}
				log.Debug("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
			case fail := <-p.Errors():
				if f != nil {
					f(topic, 0, 0)
				}
				log.Error("err: ", fail.Err)
			}
		}
	}(producer)
	for {
		select {
		case <-v.stop:
			log.Println("kafka AsyncProducer is stop")
			return nil
		case value := <-data:
			if value != "" {
				msg := &sarama.ProducerMessage{
					Topic: topic,
					Value: sarama.StringEncoder(value),
				}
				producer.Input() <- msg
			}
		}
	}
}

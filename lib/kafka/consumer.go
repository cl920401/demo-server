package kafka

import (
	"context"
	"demo-server/lib/log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

//ClusterConsumer kafka集群消费者
type Consumer struct {
	Address []string
	Version sarama.KafkaVersion
	Timeout int
	Handle  func(topic string, timestamp time.Time, partition int32, offset int64, key, value []byte) error
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		_ = consumer.Handle(message.Topic, message.Timestamp, message.Partition, message.Offset, message.Key, message.Value)
		session.MarkMessage(message, "")
	}
	return nil

}

//ClusterConsumer.Start 支持cluster的消费者
func (consumer *Consumer) Start(topics []string, groupID string) (context.CancelFunc, error) {
	conf := kafkaConf()
	if consumer.Version != (sarama.KafkaVersion{}) {
		conf.Version = consumer.Version //kafka版本号
	} else {
		conf.Version = Version //kafka版本号
	}
	conf.Consumer.Group.Session.Timeout = time.Duration(consumer.Timeout) * time.Second
	conf.Consumer.Group.Heartbeat.Interval = 6 * time.Second
	conf.Consumer.MaxProcessingTime = 500 * time.Millisecond
	conf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest

	if consumer.Handle == nil {
		consumer.Handle = func(topic string, timestamp time.Time, partition int32, offset int64, key, value []byte) error {
			log.Debugf("%s: %d-%d,%s,%s", topic, partition, offset, key, value)
			return nil
		}
	}

	client, err := sarama.NewClient(consumer.Address, conf)
	if err != nil {
		log.Errorf("error creating consumer client: %v", err)
		return nil, err
	}
	group, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		log.Errorf("error creating consumer group : %v", err)
		log.Debug("close consumer client.")
		_ = client.Close()
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
		}()

		for {
			if err := group.Consume(ctx, topics, consumer); err != nil {
				log.Errorf("error from consumer: %v", err)
				break
			}
			if ctx.Err() != nil {
				log.Error(ctx.Err())
				break
			}
		}
	}()
	go func() {
		<-ctx.Done()
		log.Info("terminating: context cancelled")
	}()
	return func() {
		cancel()
		wg.Wait()
		_ = group.Close()
		_ = client.Close()
	}, nil
}

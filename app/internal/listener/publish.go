package listener

import (
	"demo-server/lib/config"
	"demo-server/lib/kafka"
	"demo-server/lib/log"
)

var publish = make(chan string)

//var callback []func(topic string, partition int32, offset int64)
//
func InitProducerConf() {
	go func() {
		// producer init
		address := config.Get("kafka.address").StringSlice([]string{})
		topic := config.Get("kafka.producer.topic").String("")
		if topic == "" {
			return
		}
		producer := kafka.SyncProducer{
			Address: address,
			Timeout: 60,
		}
		if err := producer.Start(topic, publish, Callback); err != nil {
			log.Error("Kafka Producer Start Failed", err)
		}
	}()
}

func Callback(topic string, partition int32, offset int64) {
	//if callback == nil {
	//	return
	//}
	//for _, f := range callback {
	//	f(topic, partition, offset)
	//}
	// TODO 发送消息回调
}

// Publish 向kafka队列发布数据(线程安全)
func Publish(data string, f ...func(topic string, partition int32, offset int64)) {
	//if f != nil {
	//	callback = f
	//}
	publish <- data
	// TODO 这里清空太早，需要等回调执行完，目前没有执行回调的场景，用到时需要处理
	//callback = nil
}

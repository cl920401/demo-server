package listener

import (
	"context"
	"encoding/json"
	"runtime"
	"time"

	"github.com/pkg/errors"

	"demo-server/lib/config"
	"demo-server/lib/kafka"

	"demo-server/lib/log"
)

var isInit bool
var isClosed bool
var handler CommonHandler

var CancelFunc = func() {
	log.Errorf("CancelFunc is empty.")
}

func InitConsumerConf() context.CancelFunc {
	// consumer init
	address := config.Get("kafka.address").StringSlice([]string{})
	topics := config.Get("kafka.consumer.topics").StringSlice([]string{})
	groupID := config.Get("kafka.consumer.group").String("")
	if groupID == "" {
		return CancelFunc
	}
	cc := kafka.Consumer{
		Address: address,
		Timeout: 30,
		Handle:  safeProcessor,
	}
	cancel, err := cc.Start(topics, groupID)
	if err != nil {
		log.Error("Kafka Consumer Start Failed", err)
		return CancelFunc
	}
	CancelFunc = cancel
	return CancelFunc
}

const (
	EVT_NONE      = "none"
	EVT_ALL       = "all"
	EVT_HEARTBEAT = "report_live"
	// 业务event
)

const (
	STATUS_ONLINE = "is_online"
)

type RobotStatus struct {
	Data     string `json:"data" form:"data"`
	Event    string `json:"event" form:"event"`
	AppID    string `json:"app_id" form:"app_id"`
	Version  string `json:"appv" form:"appv"`
	Brand    string `json:"brand" form:"brand"`
	Ch       string `json:"ch" form:"ch"`
	FamilyID string `json:"family_id" form:"corpid"`
	Ctime    string `json:"ctime" form:"ctime"`
	Hwid     string `json:"hwid" form:"hwid"`
	Osv      string `json:"osv" form:"osv"`
	Pf       string `json:"pf" form:"pf"`
	RobotSN  string `json:"robot_sn" form:"rdid"`
	RobotID  string `json:"robot_id" form:"ruid"`
	Token    string `json:"token" form:"token"`
}

func safeProcessor(topic string, timestamp time.Time, partition int32, offset int64, key, value []byte) (err error) {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Error("runtime.Error", err)
			default:
				log.Error("unknown Error", err)
			}
		}
	}()

	return handler.processor(topic, timestamp, partition, offset, key, value)
}

func (v *CommonHandler) processor(topic string, timestamp time.Time, partition int32, offset int64, key, value []byte) error {
	log.Debug("processor data: ", string(value))
	if time.Since(timestamp) > time.Minute {
		return errors.New("processor error : data expire")
	}
	var data RobotStatus
	if err := json.Unmarshal(value, &data); err != nil {
		log.Error("json.Unmarshal data: ", data, "error: ", err)
		return err
	}

	for _, h := range v.listener {
		if data.Event == "" {
			break
		}
		hEvt := h.Event()
		if hEvt != EVT_ALL && data.Event != hEvt {
			continue
		}
		err := h.Processor(&data)
		if err != nil {
			log.Warn("h.Processor(&data) error: ", err, "data: ", data)
			continue
		}
	}
	return nil
}

type IHandler interface {
	Init()
	Close()
	Event() string
	Processor(*RobotStatus) error
}

type CommonHandler struct {
	listener []IHandler
}

func (v *CommonHandler) Init() {
	if isInit {
		return
	}
	InitProducerConf()
	InitConsumerConf()
	handler.AddListener(v)
	isInit = true
}

func (v *CommonHandler) Close() {
	if isClosed {
		return
	}
	log.Error("listener.Close()")
	CancelFunc()
	isClosed = true
}

func (v *CommonHandler) Event() string {
	return EVT_ALL
}

func (v *CommonHandler) Processor(data *RobotStatus) error {
	log.Debug("Processor data: ", *data)
	return nil
}

func (v *CommonHandler) AddListener(listener ...IHandler) {
	v.listener = append(v.listener, listener...)
}

func GetStatusHandler() *CommonHandler {
	return &handler
}

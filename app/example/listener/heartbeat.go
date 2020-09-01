package listener

import (
	"context"
	"demo-server/app/internal/listener"
	"time"

	"demo-server/lib/log"
)

type HeartbeatListener struct{}

var Rate = 5 * time.Minute
var TimerCancel func()

func (h *HeartbeatListener) Init() {
	common := listener.GetStatusHandler()
	common.Init()
	common.AddListener(h)
	TimerCancel = AddTimer(Rate)
}

func (h *HeartbeatListener) Close() {
	TimerCancel()
	common := listener.GetStatusHandler()
	common.Close()
}

func (h *HeartbeatListener) Event() string {
	return listener.EVT_HEARTBEAT
}

func (h *HeartbeatListener) Processor(data *listener.RobotStatus) error {
	log.Debug("HeartbeatListener")
	robot := RobotStatus{
		Key:     "is_online",
		Value:   "1",
		RobotSN: data.RobotSN,
	}
	if err := robot.Store(nil); err != nil {
		log.Error("HeartbeatListener Processor robot.Store(tx) error : ", err)
		return err
	}
	return nil
}

func AddTimer(rate time.Duration) func() {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		t := time.NewTicker(rate)
		for {
			select {
			case <-t.C:
				log.Debug("timer start")
				UpdataOnline()
			case <-ctx.Done():
				t.Stop()
				break
			}
		}
	}(ctx)
	return cancel
}

func UpdataOnline() {
	robot := RobotStatus{}
	if err := robot.UpdateOnline(nil, 3*time.Minute); err != nil {
		log.Error("robot status UpdataOnline() failed", err)
	}
}

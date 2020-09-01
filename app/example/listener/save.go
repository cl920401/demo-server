package listener

import (
	"demo-server/app/internal/listener"
	"encoding/json"
	"fmt"
	"time"

	"demo-server/lib/log"
	"demo-server/lib/mysql"
)

type RobotStatus struct {
	ID        int       `json:"id"`
	Key       string    `json:"status_key"`
	Value     string    `json:"status_value"`
	RobotSN   string    `json:"robot_sn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" `
	DeletedAt int64     `json:"deleted_at,omitempty"`
}
type SaveListener struct{}

func (h *SaveListener) Init() {
	common := listener.GetStatusHandler()
	common.Init()
	common.AddListener(h)
}

func (h *SaveListener) Close() {
	common := listener.GetStatusHandler()
	common.Close()
}

func (h *SaveListener) Event() string {
	return listener.EVT_ALL
}

func (h *SaveListener) Processor(data *listener.RobotStatus) error {
	input := make(map[string]string)
	input["event"] = data.Event
	input["app_id"] = data.AppID
	input["appv"] = data.Version
	input["brand"] = data.Brand
	input["ch"] = data.Ch
	input["family_id"] = data.FamilyID
	input["ctime"] = data.Ctime
	input["hwid"] = data.Hwid
	input["osv"] = data.Osv
	input["pf"] = data.Pf
	input["robot_sn"] = data.RobotSN
	input["robot_id"] = data.RobotID
	input["token"] = data.Token
	var kv map[string]interface{}

	if err := json.Unmarshal([]byte(data.Data), &kv); err != nil {
		log.Error("json.Unmarshal([]byte(data.Data)) error, data=", data, "error=", err)
		return err
	}
	for k, v := range kv {
		strValue, err := json.Marshal(v)
		if err != nil {
			log.Error("json.Marshal(v) error, v=", v, "error=", err)
			continue
		}
		input[k] = string(strValue)
	}
	db := mysql.Default()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("SaveListener Processor tx Begin err:[%w]", err)
	}
	defer mysql.Rollback(tx)

	// 保存机器人状态
	for key, value := range input {
		if key == "" {
			log.Error("SaveListener Processor error : key invalid")
			continue
		}
		robot := RobotStatus{
			Key:     key,
			Value:   value,
			RobotSN: data.RobotSN,
		}
		if err := robot.Store(tx); err != nil {
			log.Error("SaveListener Processor robot.Store(tx) error : ", err)
			return err
		}
	}
	return tx.Commit()
}

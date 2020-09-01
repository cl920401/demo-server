package status

import (
	"fmt"
	"time"

	"demo-server/lib/log"

	"demo-server/lib/mysql"

	sq "github.com/Masterminds/squirrel"
)

//RoomMap 地图表
type RobotStatus struct {
	ID        int       `json:"id"`
	Key       string    `json:"status_key"`
	Value     string    `json:"status_value"`
	RobotSN   string    `json:"robot_sn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" `
	DeletedAt int64     `json:"deleted_at,omitempty"`
}

//TableName 返回表名
func (*RobotStatus) TableName() string {
	return "robot_status"
}

func (t *RobotStatus) Store(tx sq.BaseRunner) error {
	if tx == nil {
		tx = mysql.Default()
	}
	_, err := tx.Exec(
		"REPLACE INTO robot_status (status_key, status_value, robot_sn, updated_at) VALUES (?,?,?,?)",
		t.Key, t.Value, t.RobotSN, time.Now())
	if err != nil {
		return err
	}
	return nil

}

func (t *RobotStatus) GetStatus(tx sq.BaseRunner, robotSN []string) (map[string][]RobotStatus, error) {
	if tx == nil {
		tx = mysql.Default()
	}
	b := sq.Select("status_key", "status_value", "robot_sn", "updated_at").
		From(t.TableName()).
		Where(sq.Eq{"robot_sn": robotSN, "deleted_at": 0}).
		Limit(1000).
		RunWith(tx)
	log.Debug(b.ToSql())

	rows, err := b.Query()
	if err != nil {
		return nil, err
	}
	defer mysql.CloseRows(rows)

	retMap := make(map[string][]RobotStatus)
	for rows != nil && rows.Next() {
		var item RobotStatus
		if err := rows.Scan(
			&item.Key,
			&item.Value,
			&item.RobotSN,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("GetStatus row err:[%w]", err)
		}
		if _, ok := retMap[item.RobotSN]; !ok {
			retMap[item.RobotSN] = make([]RobotStatus, 0)
		}
		retMap[item.RobotSN] = append(retMap[item.RobotSN], item)
	}
	return retMap, nil
}

func (t *RobotStatus) UpdateOnline(tx sq.BaseRunner, timeout time.Duration) error {
	if tx == nil {
		tx = mysql.Default()
	}
	// 获取本地时间
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	localTime := time.Now().In(cstSh)
	b := sq.Update(t.TableName()).SetMap(map[string]interface{}{
		"status_value": "0",
		"updated_at":   localTime,
	}).Where(sq.And{sq.Lt{"updated_at": localTime.Add(-timeout)}, sq.Eq{"status_key": "is_online", "deleted_at": 0}}).RunWith(tx)
	log.Debug(b.ToSql())

	if _, err := b.Exec(); err != nil {
		return fmt.Errorf("model status UpdateOnline() error:[%w]", err)
	}

	return nil

}

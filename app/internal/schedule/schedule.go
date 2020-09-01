package schedule

const (
	VERSION1 = 1
	VERSION2 = 2
)

//Event 日程中的事件
type Event struct {
	OpkID         string  `json:"opk_id"`
	Title         string  `json:"title"`
	Desc          string  `json:"desc"`
	Status        int     `json:"status"`
	StartTime     int64   `json:"start_time"`
	ExecutionTime int     `json:"execution_time"`
	Icon          string  `json:"icon"`
	Children      []Event `json:"children"`
}

// DaySchedule 一天的日程
type DaySchedule struct {
	Date     int64   `json:"date"`
	Schedule []Event `json:"schedule"`
}

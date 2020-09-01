package mysql

import (
	"database/sql"
	"demo-server/lib/config"
	"demo-server/lib/log"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // init mysql
	"sync"
)

type mysqlConf struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	DBName       string `json:"dbname"`
	Host         string `json:"host"`
	Port         int16  `json:"port"`
	Charset      string `json:"charset"`
	Timeout      int    `json:"timeout"`
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
	Loc          string `json:"loc"`
}

//保存连接对象
var conn = make(map[string]*sql.DB)

var isInit bool
var rw sync.RWMutex

func (v mysqlConf) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&readTimeout=%ds&writeTimeout=%ds&timeout=%ds&parseTime=true&loc=%s",
		v.Username, v.Password, v.Host, v.Port, v.DBName, v.Charset, v.Timeout, v.Timeout, v.Timeout, v.Loc,
	)
}

func (v mysqlConf) maxOpenConns() int {
	if v.MaxOpenConns == 0 {
		log.Fatal("MaxOpenConns Must set")
	}
	return v.MaxOpenConns
}

func (v mysqlConf) maxIdleConns() int {
	if v.MaxIdleConns == 0 {
		log.Fatal("MaxIdleConns Must set")
	}
	return v.MaxIdleConns
}

//InitDB db初始化
func InitDB() {
	if isInit {
		return
	}
	dbConf := config.Get("mysql")
	var data map[string]*mysqlConf
	if err := dbConf.Scan(&data); err != nil {
		log.Fatal("Error parsing database configuration file ", err)
		return
	}
	for k, v := range data {
		dbConn, err := NewMysql(v)
		if err != nil {
			log.Fatal("InitDB ERROR ", k, " ", v.String())
			return
		}
		conn[k] = dbConn
	}
	isInit = true
}

//NewMysql 实例化db
func NewMysql(conf *mysqlConf) (Db *sql.DB, err error) {

	Db, err = sql.Open("mysql", conf.String())
	if err != nil {
		log.Errorf("mysql conn Error %s", err.Error())
		return
	}
	err = Db.Ping()
	if err != nil {
		log.Errorf("Could not establish a connection with the database, detail: %s", err.Error())
		return
	}
	Db.SetMaxOpenConns(conf.maxOpenConns())
	Db.SetMaxIdleConns(conf.maxIdleConns())
	return
}

//DB 对外获取db实例
func DB(key string) *sql.DB {
	if !isInit {
		rw.Lock()
		defer rw.Unlock()
		InitDB()
	}
	if db, ok := conn[key]; ok && db != nil {
		return db
	}
	return nil
}

//Default 获取默认数据库实例
func Default() *sql.DB {
	return DB("default")
}

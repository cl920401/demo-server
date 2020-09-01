package elasticsearch

import (
	"demo-server/lib/config"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/olivere/elastic"
	aws "github.com/olivere/elastic/aws/v4"
	"log"
	"net/http"
	"sync"
)

var once = sync.Once{}
var conn = make(map[string]*elastic.Client)

const EsTypeAws = "aws"

type conf struct {
	Host      string `json:"host"`
	Type      string `json:"type"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
	Token     string `json:"token"`
}

//初始化
func InitDB() {
	once.Do(func() {
		esConfig := config.Get("elasticsearch")
		var data map[string]*conf
		if err := esConfig.Scan(&data); err != nil {
			log.Fatal("Error parsing elasticsearch configuration file ", err)
			return
		}
		for k, v := range data {
			var client *http.Client
			if v.Type == EsTypeAws {
				client = aws.NewV4SigningClient(
					credentials.NewStaticCredentials(v.AccessKey, v.SecretKey, v.Token), v.Region)
			}

			esConn, err := elastic.NewClient(
				elastic.SetURL(v.Host),
				elastic.SetSniff(false),
				elastic.SetHttpClient(client))
			if err != nil {
				log.Fatalf("error: %s", err.Error())
				return
			}
			conn[k] = esConn
		}
	})
}

//DB 对外获取db实例
func DB(key string) *elastic.Client {
	InitDB()
	if db, ok := conn[key]; ok && db != nil {
		return db
	}
	return nil
}

//Default 获取默认数据库实例
func Default() *elastic.Client {
	return DB("default")
}

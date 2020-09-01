package config

import (
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/encoder/yaml"
	"github.com/micro/go-micro/v2/config/reader"
	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/config/source/env"
	"github.com/micro/go-micro/v2/config/source/etcd"
	"github.com/micro/go-micro/v2/config/source/flag"
	"github.com/micro/go-micro/v2/config/source/memory"
	"io/ioutil"
	"strings"
)

//data 保存配置文件信息
var data, _ = config.NewConfig()

func Data() config.Config {
	return data
}

//ToString 输出字符串
func ToString() string {
	return string(data.Bytes())
}

//LoadMemory 从字符串解析配置信息
func LoadMemory(confData []byte) error {
	return data.Load(memory.NewSource(memory.WithJSON(confData)))
}

//LoadEnv 从环境变量解析配置信息
func LoadEnv(prefix string) error {
	return data.Load(NewSourceEnv(prefix))
}

//LoadEtcd 从etcd加载配置文件
func LoadEtcd(prefix string, username, password string, address ...string) error {
	etcdSource := etcd.NewSource(
		source.WithEncoder(yaml.NewEncoder()),
		// optionally specify etcd address; default to localhost:2379
		etcd.WithAddress(address...),
		etcd.WithPrefix(prefix),
		// optionally strip the provided prefix from the keys, defaults to false
		etcd.StripPrefix(true),
		//etcd auth
		etcd.Auth(username, password),
	)
	return data.Load(etcdSource)
}

//LoadFlag 从flag解析配置信息
func LoadFlag() error {
	return data.Load(NewSourceFlag())
}

//LoadMultiple 多数据源加载
func LoadMultiple(source ...source.Source) error {
	return data.Load(source...)
}

//NewSourceEnv 环境变量
func NewSourceEnv(prefix string) source.Source {
	return env.NewSource(env.WithPrefix(prefix))
}

//NewSourceFlag 数据源flag
func NewSourceFlag() source.Source {
	return flag.NewSource(flag.IncludeUnset(true))
}

//NewSourceFile 数据源文件
func NewSourceFile(filePath string) source.Source {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil
	}
	return memory.NewSource(memory.WithJSON(data))
}

//LoadFile 从文件解析配置信息
func LoadFile(filePath string) error {
	confData, err := ioutil.ReadFile(filePath)
	if err == nil {
		return LoadMemory(confData)
	}
	return err
}

//Get 获取配置信息
func Get(key string) reader.Value {
	path := strings.Split(key, ".")
	return data.Get(path...)
}

//Scan 配置信息解析到结构体
func Scan(v interface{}) error {
	return data.Scan(&v)
}

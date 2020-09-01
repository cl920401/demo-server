package apm

import (
	"demo-server/lib/apm/output"
	"demo-server/lib/config"
	"demo-server/lib/log"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"
	"net/http"
	"sync"
	"time"
)

//全局registry
var mRegistry = metrics.NewRegistry()

const tick = 3 * time.Second

var keys = sync.Map{}

func init() {
	//添加http Handle /debug/metrics
	http.HandleFunc("/debug/metrics", func(writer http.ResponseWriter, request *http.Request) {
		if config.Get("apm.enable").Bool(false) {
			exp.ExpHandler(mRegistry).ServeHTTP(writer, request)
		} else {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte("{}"))
		}
	})

	metrics.RegisterRuntimeMemStats(mRegistry)
	metrics.RegisterDebugGCStats(mRegistry)
	//prometheusRegistry := prometheus.NewRegistry()
	//pClient := output.NewPrometheusProvider(mRegistry, "test", "subsys", prometheusRegistry)
	es := output.NewElasticSearch(mRegistry)
	go func() {
		for range time.Tick(tick) {
			if config.Get("apm.enable").Bool(false) {
				metrics.CaptureDebugGCStatsOnce(mRegistry)    //更新gc信息
				metrics.CaptureRuntimeMemStatsOnce(mRegistry) //更新内存状态信息
				err := es.UpdateElasticSearchMetricsOnce()    //上报到es
				if err != nil {
					log.Error(err)
				}
			} else {
				//关闭apm监控
				keys.Range(func(key, value interface{}) bool {
					mRegistry.Unregister(key.(string))
					keys.Delete(key)
					return true
				})
			}
		}
	}()
}

func GetRegistry() metrics.Registry {
	return mRegistry
}

//简单的计数器 可以增加 减少
func Counter(key, t string) metrics.Counter {
	if config.Get("apm.enable").Bool(false) {
		k := "counter." + t + "." + key
		if _, ok := keys.LoadOrStore(k, true); ok {
			return mRegistry.Get(k).(metrics.Counter)
		}
		return mRegistry.GetOrRegister(k, metrics.NewCounter()).(metrics.Counter)
	}
	return metrics.NilCounter{}
}

//自增的计数器,用来度量一系列事件发生的比率 提供了平均速率，以及指数平滑平均速率，以及采样后的1分钟，5分钟，15分钟速率
func Meter(key, t string) metrics.Meter {
	if config.Get("apm.enable").Bool(false) {
		k := "meter." + t + "." + key
		if _, ok := keys.LoadOrStore(k, true); ok {
			return mRegistry.Get(k).(metrics.Meter)
		}
		return mRegistry.GetOrRegister(k, metrics.NewMeter()).(metrics.Meter)
	}
	return metrics.NilMeter{}
}

//用来记录一些对象或者事物的瞬时值
func Gauges(key, t string) metrics.Gauge {
	if config.Get("apm.enable").Bool(false) {
		k := "gauge." + t + "." + key
		if _, ok := keys.LoadOrStore(k, true); ok {
			return mRegistry.Get(k).(metrics.Gauge)
		}
		return mRegistry.GetOrRegister(k, metrics.NewGauge()).(metrics.Gauge)
	}
	return metrics.NilGauge{}
}

//统计数据的分布情况 比如最小值，最大值，中间值，还有中位数，75百分位, 90百分位, 95百分位, 98百分位, 99百分位, 和 99.9百分位的值
func Histograms(key, t string) metrics.Histogram {
	if config.Get("apm.enable").Bool(false) {
		k := "histogram." + t + "." + key
		if _, ok := keys.LoadOrStore(k, true); ok {
			return mRegistry.Get(k).(metrics.Histogram)
		}
		return mRegistry.GetOrRegister(k,
			metrics.NewHistogram(metrics.NewExpDecaySample(1028, 0.015))).(metrics.Histogram)
	}
	return metrics.NilHistogram{}
}

//统计当前请求的速率和处理时间
func Timer(key, t string) metrics.Timer {
	if config.Get("apm.enable").Bool(false) {
		k := "timer." + t + "." + key
		if _, ok := keys.LoadOrStore(k, true); ok {
			return mRegistry.Get(k).(metrics.Timer)
		}
		return mRegistry.GetOrRegister(k, metrics.NewTimer()).(metrics.Timer)
	}
	return metrics.NilTimer{}
}

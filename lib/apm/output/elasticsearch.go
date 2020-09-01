package output

import (
	"context"
	"demo-server/lib/config"
	"demo-server/lib/elasticsearch"
	"demo-server/lib/log"
	"errors"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/rcrowley/go-metrics"
	"time"
)

type ElasticSearch struct {
	Registry     metrics.Registry // Registry to be exported
	DurationUnit time.Duration    // Time conversion unit for durations
}

func NewElasticSearch(r metrics.Registry) ElasticSearch {
	return ElasticSearch{
		Registry:     r,
		DurationUnit: time.Nanosecond,
	}
}

func (c *ElasticSearch) UpdateElasticSearchMetricsOnce() error {
	t := time.Now()
	du := float64(c.DurationUnit)
	client := elasticsearch.DB("apm")
	if client == nil {
		return errors.New("elasticsearch.apm profile not configured or connect error")
	}
	bulkService := client.Bulk()
	var n = 0
	var m = config.Get("elasticsearch.apm.bulk").Int(50)
	c.Registry.Each(func(name string, i interface{}) {

		switch metric := i.(type) {
		case metrics.Counter:
			item := map[string]interface{}{
				"name":    name,
				"service": ServiceName(),
				"value":   metric.Count(),
				"host":    ShortHostname(),
				"time":    t,
			}
			index := fmt.Sprintf("%s-apm-counter-%s", config.Get("service.env").String("dev"), time.Now().Format("2006-01-02"))
			bulkService.Add(elastic.NewBulkIndexRequest().Index(index).Type("metrics").Doc(item))
		case metrics.Gauge:
			item := map[string]interface{}{
				"name":    name,
				"service": ServiceName(),
				"value":   metric.Value(),
				"host":    ShortHostname(),
				"time":    t,
			}
			index := fmt.Sprintf("%s-apm-gauge-%s", config.Get("service.env").String("dev"), time.Now().Format("2006-01-02"))
			bulkService.Add(elastic.NewBulkIndexRequest().Index(index).Type("metrics").Doc(item))

		case metrics.GaugeFloat64:
			item := map[string]interface{}{
				"name":    name,
				"service": ServiceName(),
				"value":   metric.Value(),
				"host":    ShortHostname(),
				"time":    t,
			}
			index := fmt.Sprintf("%s-apm-gaugefloat-%s", config.Get("service.env").String("dev"), time.Now().Format("2006-01-02"))
			bulkService.Add(elastic.NewBulkIndexRequest().Index(index).Type("metrics").Doc(item))
		case metrics.Histogram:
			index := fmt.Sprintf("%s-apm-histogram-%s", config.Get("service.env").String("dev"), time.Now().Format("2006-01-02"))
			h := metric.Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			item := map[string]interface{}{
				"name":           name,
				"service":        ServiceName(),
				"host":           ShortHostname(),
				"count":          h.Count(),
				"min":            h.Min(),
				"max":            h.Max(),
				"mean":           h.Mean(),
				"std-dev":        h.StdDev(),
				"50-percentile":  ps[0],
				"75-percentile":  ps[1],
				"95-percentile":  ps[2],
				"99-percentile":  ps[3],
				"999-percentile": ps[4],
				"time":           t,
			}
			bulkService.Add(elastic.NewBulkIndexRequest().Index(index).Type("metrics").Doc(item))

		case metrics.Meter:
			index := fmt.Sprintf("%s-apm-meter-%s", config.Get("service.env").String("dev"), time.Now().Format("2006-01-02"))
			m := metric.Snapshot()
			item := map[string]interface{}{
				"name":           name,
				"service":        ServiceName(),
				"host":           ShortHostname(),
				"count":          m.Count(),
				"one-minute":     m.Rate1(),
				"five-minute":    m.Rate5(),
				"fifteen-minute": m.Rate15(),
				"mean":           m.RateMean(),
				"time":           t,
			}

			bulkService.Add(elastic.NewBulkIndexRequest().Index(index).Type("metrics").Doc(item))

		case metrics.Timer:
			index := fmt.Sprintf("%s-apm-timer-%s", config.Get("service.env").String("dev"), time.Now().Format("2006-01-02"))
			ms := metric.Snapshot()
			ps := ms.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			item := map[string]interface{}{
				"host":           ShortHostname(),
				"name":           name,
				"service":        ServiceName(),
				"count":          ms.Count(),
				"min":            ms.Min() / int64(du),
				"max":            ms.Max() / int64(du),
				"mean":           ms.Mean() / du,
				"std-dev":        ms.StdDev() / du,
				"50-percentile":  ps[0] / du,
				"75-percentile":  ps[1] / du,
				"95-percentile":  ps[2] / du,
				"99-percentile":  ps[3] / du,
				"999-percentile": ps[4] / du,
				"one-minute":     ms.Rate1(),
				"five-minute":    ms.Rate5(),
				"fifteen-minute": ms.Rate15(),
				"mean-rate":      ms.RateMean(),
				"time":           t,
			}
			bulkService.Add(elastic.NewBulkIndexRequest().Index(index).Type("metrics").Doc(item))
		}
		n++
		if n%m == 0 && bulkService.NumberOfActions() > 0 {
			if _, err := bulkService.Do(context.Background()); err != nil {
				log.Error(err)
			}
			bulkService.Reset()
		}
	})
	if bulkService.NumberOfActions() > 0 {
		_, err := bulkService.Do(context.Background())
		return err
	}
	return nil
}

package log

import (
	"strconv"
	"time"

	"github.com/i5heu/simple-S3-cache/internal/config"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type Logger struct {
	Idb     influxdb2.Client
	Write   api.WriteAPI
	Enabled bool
}

func New(conf config.Config) Logger {
	if conf.InfluxDbToken == "" {
		return Logger{Enabled: false}
	}

	client := influxdb2.NewClient(conf.InfluxDbUrl, conf.InfluxDbToken)
	writeAPI := client.WriteAPI(conf.InfluxDbOrg, conf.InfluxDbBucket)

	l := Logger{Idb: client, Write: writeAPI, Enabled: true}
	go l.flushWorker()

	return l
}

func (l *Logger) flushWorker() {
	for {
		time.Sleep(1 * time.Second)
		l.Write.Flush()
	}
}

func (l *Logger) LogRequest(timeStart time.Time, url string, statusCode int, cached bool, fileSize uint) {
	if !l.Enabled {
		return
	}

	p := influxdb2.NewPointWithMeasurement("stat").
		AddTag("statusCode", strconv.Itoa(statusCode)).
		AddTag("cached", strconv.FormatBool(cached)).
		AddField("duration", time.Since(timeStart).Microseconds()/1000).
		AddField("file_size", fileSize).
		AddField("url", url).
		SetTime(time.Now())

	l.Write.WritePoint(p)
}

func (l *Logger) LogCache(timeStart time.Time, cacheName string, cacheSize uint, cacheSizeMax uint) {
	if !l.Enabled {
		return
	}

	p := influxdb2.NewPointWithMeasurement("stat-cache").
		AddTag("cacheName", cacheName).
		AddField("cacheSize", cacheSize).
		AddField("cacheSizeMax", cacheSizeMax).
		AddField("duration", time.Since(timeStart).Microseconds()/1000).
		SetTime(time.Now())

	l.Write.WritePoint(p)
}

package log

import (
	"simple-S3-cache/config"
	"strconv"
	"time"

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

	l := Logger{Idb: client, Write: writeAPI}
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
		AddField("duration", time.Since(timeStart).Milliseconds()).
		AddField("file_size", fileSize).
		AddField("url", url).
		SetTime(time.Now())

	l.Write.WritePoint(p)
}

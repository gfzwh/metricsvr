package controller

import (
	"context"
	"strings"

	"go.uber.org/zap"

	"github.com/gfzwh/gfz/config"
	"github.com/gfzwh/gfz/proto"
	"github.com/gfzwh/gfz/udp"
	"github.com/gfzwh/gfz/zzlog"
)

func (this *controller) ipport(host string) (ip, port string) {
	hts := strings.Split(host, ":")
	if 0 < len(hts) {
		ip = hts[0]
	}

	if 1 < len(hts) {
		port = hts[1]
	}

	return
}

func (this *controller) OnEvent(ctx context.Context, data []byte, ext *udp.EventInfo) (err error) {
	code := "0"
	defer func() {
		Counter(config.Get("server", "name").String(""), "OnEvent", "", "", code)
		if r := recover(); r != nil {
			zzlog.Errorw("controller.OnEvent error", zap.Int("size", len(data)), zap.Error(r.(error)))
		}
	}()
	var metric proto.Metrics
	err = metric.Unmarshal(data)
	if nil != err {
		code = "500"
		zzlog.Errorw("OnEvent Unmarshal error", zap.Any("size", len(data)), zap.Error(err))

		return
	}

	zzlog.Debugw("udp.OnEvent call ", zap.Any("svrname", metric.Svrname), zap.Any("host", metric.Host), zap.Any("size", len(data)))

	ip, port := this.ipport(metric.Host)
	switch metric.Type {
	case proto.MetricType_CounterType:
		Counter(metric.Svrname, metric.Counter.Method, ip, port, metric.Counter.Code)
		Counter(metric.Svrname, metric.Counter.Method, "all", "", metric.Counter.Code)

		break
	case proto.MetricType_GaugeType:
		if metric.Gauge.Inc {
			GaugeInc(metric.Svrname, metric.Gauge.Type, ip, port, metric.Gauge.Value)
			GaugeInc(metric.Svrname, metric.Gauge.Type, "all", "", metric.Gauge.Value)
		} else {
			Gauge(metric.Svrname, metric.Gauge.Type, metric.Gauge.Value, ip, port, metric.Gauge.Add)
			Gauge(metric.Svrname, metric.Gauge.Type, metric.Gauge.Value, "all", "", metric.Gauge.Add)
		}

		break
	case proto.MetricType_SummaryType:
		Summary(metric.Svrname, metric.Summary.Method, ip, port, metric.Micro)
		Summary(metric.Svrname, metric.Summary.Method, "all", "", metric.Micro)

		break
	}
	return
}

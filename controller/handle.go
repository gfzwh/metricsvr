package controller

import (
	"context"

	"go.uber.org/zap"

	"github.com/gfzwh/gfz/config"
	"github.com/gfzwh/gfz/proto"
	"github.com/gfzwh/gfz/udp"
	"github.com/gfzwh/gfz/zzlog"
)

func (this *controller) OnEvent(ctx context.Context, data []byte, ext *udp.EventInfo) (err error) {
	code := "0"
	defer func() {
		Counter(config.Get("server", "name").String(""), "OnEvent", code)
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

	switch metric.Type {
	case proto.MetricType_CounterType:
		Counter(metric.Svrname, metric.Counter.Method, metric.Counter.Code)

		break
	case proto.MetricType_GaugeType:
		if metric.Gauge.Inc {
			GaugeInc(metric.Svrname, metric.Gauge.Type, metric.Gauge.Value)
		} else {
			Gauge(metric.Svrname, metric.Gauge.Type, metric.Gauge.Value, metric.Gauge.Add)
		}
		break
	case proto.MetricType_SummaryType:
		Summary(metric.Svrname, metric.Summary.Method, metric.Micro)
		break
	}
	return
}

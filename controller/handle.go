package controller

import (
	"strings"

	"go.uber.org/zap"

	"github.com/gfzwh/gfz/proto"
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

type IConsumer interface {
	Error(error)
	Notify(interface{})
	Message([]byte)
}

func (this *controller) Error(err error) {
	zzlog.Errorw("kafka error", zap.Error(err))
}

func (this *controller) Notify(notify interface{}) {

}

func (this *controller) Message(partition int32, offset int64, message []byte) {
	data := message
	defer func() {
		if r := recover(); r != nil {
			zzlog.Errorw("controller.Message error", zap.Int("size", len(data)),
				zap.Error(r.(error)))
		}
	}()
	var metric proto.Metrics
	err := metric.Unmarshal(data)
	if nil != err {
		zzlog.Errorw("controller.Message Unmarshal error", zap.Any("size", len(data)), zap.Error(err))

		return
	}

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

	zzlog.Debugw("controller.Message call ", zap.Any("svrname", metric.Svrname),
		zap.Any("host", metric.Host), zap.Any("size", len(data)))
	return
}

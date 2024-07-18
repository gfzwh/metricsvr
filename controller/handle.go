package controller

import (
	"strings"

	"go.uber.org/zap"

	"github.com/shockerjue/gffg/proto"
	"github.com/shockerjue/gffg/zzlog"
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

	for _, it := range metric.Lists {
		ip, port := this.ipport(it.Host)
		switch it.Type {
		case proto.MetricType_CounterType:
			Counter(it.Svrname, it.Counter.Method, ip, port, it.Counter.Code)
			Counter(it.Svrname, it.Counter.Method, "all", "", it.Counter.Code)

			break
		case proto.MetricType_GaugeType:
			if it.Gauge.Inc {
				GaugeInc(it.Svrname, it.Gauge.Type, ip, port, it.Gauge.Value)
				GaugeInc(it.Svrname, it.Gauge.Type, "all", "", it.Gauge.Value)
			} else {
				Gauge(it.Svrname, it.Gauge.Type, it.Gauge.Value, ip, port, it.Gauge.Add)
				Gauge(it.Svrname, it.Gauge.Type, it.Gauge.Value, "all", "", it.Gauge.Add)
			}

			break
		case proto.MetricType_SummaryType:
			Summary(it.Svrname, it.Summary.Method, ip, port, it.Micro)
			Summary(it.Svrname, it.Summary.Method, "all", "", it.Micro)

			break
		}

		zzlog.Debugw("controller.Message call ", zap.Any("svrname", it.Svrname),
			zap.Any("host", it.Host), zap.Any("size", len(data)))
	}

	return
}

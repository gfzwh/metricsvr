package controller

import (
	"sync"
	"time"

	"github.com/gfzwh/gfz/config"
	"github.com/gfzwh/gfz/zzlog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"go.uber.org/zap"
)

type mitem struct {
	registry   *prometheus.Registry
	counterVec *prometheus.CounterVec
	gfzVec     *prometheus.GaugeVec
	summaryVec *prometheus.SummaryVec
}

type metrics struct {
	rw     sync.RWMutex
	metric map[string]*mitem

	l     sync.Mutex
	initm map[string]*sync.Once
}

var instance *metrics
var once sync.Once

func Metrics() *metrics {
	once.Do(func() {
		instance = &metrics{
			rw:     sync.RWMutex{},
			metric: make(map[string]*mitem),
			initm:  make(map[string]*sync.Once),
			l:      sync.Mutex{},
		}
	})

	return instance
}

func (m *metrics) Mitem(svrname string) *mitem {
	m.rw.RLock()
	if v, ok := m.metric[svrname]; ok {
		m.rw.RUnlock()

		return v
	}
	m.rw.RUnlock()

	m.l.Lock()
	defer m.l.Unlock()
	if _, ok := m.initm[svrname]; !ok {
		m.initm[svrname] = &sync.Once{}
	}

	m.initm[svrname].Do(func() {
		instance := &mitem{
			registry: prometheus.NewRegistry(),
			counterVec: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "gfz_call",
					Help: "How many RPC requests processed, partitioned by status code and RPC method.",
				},
				[]string{"code", "method", "host"},
			),
			// 用于统计用（链接、断开连接,panic,错误,...）
			gfzVec: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: "",
					Subsystem: "",
					Name:      "gfz_sys",
					Help:      "Number of blob storage operations waiting to be processed, partitioned by user and type.",
				},
				[]string{
					// Which user has requested the operation?
					"type",
					// Of what type is the operation?
					"value",

					"host",
				},
			),
			summaryVec: prometheus.NewSummaryVec(
				prometheus.SummaryOpts{
					Name:       "gfz_call_delay",
					Help:       "The temperature of the frog pond.",
					Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
				},
				[]string{"method", "host"},
			),
		}

		instance.registry.MustRegister(instance.counterVec)
		instance.registry.MustRegister(instance.gfzVec)
		instance.registry.MustRegister(instance.summaryVec)

		go func() {
			for {
				if err := push.New(config.Get("metrics", "push").String(""), svrname).
					Collector(instance.counterVec).
					Collector(instance.gfzVec).
					Collector(instance.summaryVec).
					Push(); err != nil {
					zzlog.Errorw("push metrics err", zap.Error(err))
				}

				instance.gfzVec.Reset()
				instance.counterVec.Reset()
				instance.summaryVec.Reset()

				time.Sleep(1 * time.Second)
			}
		}()

		m.rw.Lock()
		m.metric[svrname] = instance
		m.rw.Unlock()
	})

	return m.metric[svrname]
}

func GaugeInc(svrname, _type, host, value string) {
	Metrics().Mitem(svrname).gfzVec.With(prometheus.Labels{"type": _type, "value": value, "host": host}).Inc()
}

func Gauge(svrname, _type, value, host string, add int64) {
	Metrics().Mitem(svrname).gfzVec.With(prometheus.Labels{"type": _type, "value": value, "host": host}).Set(float64(add))
}

func Counter(svrname, method, host, code string) {
	Metrics().Mitem(svrname).counterVec.WithLabelValues(code, method, host).Inc()
}

func Summary(svrname, method, host string, micro int64) {
	Metrics().Mitem(svrname).summaryVec.WithLabelValues(method, host).Observe(float64(micro))
}

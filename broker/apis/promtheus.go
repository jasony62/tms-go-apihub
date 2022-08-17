package apis

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	klog "k8s.io/klog/v2"
)

var httpInPromCounter *prometheus.CounterVec
var httpOutPromCounter *prometheus.CounterVec
var httpInDurationPromHistogram *prometheus.HistogramVec
var httpOutDurationPromHistogram *prometheus.HistogramVec

var httpInPromCounterMap map[string]*prometheus.CounterVec
var httpOutPromCounterMap map[string]*prometheus.CounterVec
var httpInDurationPromHistogramMap map[string]*prometheus.HistogramVec
var httpOutDurationPromHistogramMap map[string]*prometheus.HistogramVec

var httpInPromCounterFlowMap map[string]*prometheus.CounterVec
var httpOutPromCounterFlowMap map[string]*prometheus.CounterVec
var httpInDurationPromHistogramFlowMap map[string]*prometheus.HistogramVec
var httpOutDurationPromHistogramFlowMap map[string]*prometheus.HistogramVec

func promStart(stack *hub.Stack, params map[string]string) (interface{}, int) {
	klog.Infoln("promStart!")
	host := params["host"]
	port := params["port"]

	if len(host) == 0 {
		host = "0.0.0.0"
	}
	if len(port) == 0 {
		port = "8000"
	}
	klog.Infoln("promStart: host: ", host, " port:", port)

	promStartRun(fmt.Sprintf("%s:%s", host, port))
	return nil, http.StatusOK
}

func promStartRun(address string) {
	promInitData()

	// Expose the registered metrics via HTTP
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	go func() {
		klog.Infoln("Listen and Serve: ", address)
		if err := http.ListenAndServe(address, nil); err != nil {
			klog.Fatal("Error in ListenAndServe: %v", err)
		}
	}()
}

func getPromLabels(params map[string]string) map[string]string {
	return prometheus.Labels{
		"code":  params["code"],
		"child": params["child"],
		"root":  params["root"],
		"type":  params["type"]}
}

func promHttpCounterInc(stack *hub.Stack, params map[string]string) (interface{}, int) {
	promLabels := getPromLabels(params)

	klog.Infoln("promLabels: ", promLabels)
	val := params["duration"]
	duration, err := strconv.ParseFloat(val, 64)
	if err != nil {
		str := "解析http out duration失败, err: " + err.Error()
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusBadRequest
	}
	duration = duration * 10
	if params["httpInOut"] == "httpIn" {
		httpInDurationPromHistogram.With(promLabels).Observe(duration)
		httpInPromCounter.With(promLabels).Inc()

		if promLabels["type"] == "httpapi" {
			klog.Infoln("httpInPromCounterMap: ", promLabels["root"])
			httpInPromCounterMap[promLabels["root"]].With(promLabels).Inc()
			httpInDurationPromHistogramMap[promLabels["root"]].With(promLabels).Observe(duration)
		} else if promLabels["type"] == "flow" {
			httpInPromCounterFlowMap[promLabels["root"]].With(promLabels).Inc()
			httpInDurationPromHistogramFlowMap[promLabels["root"]].With(promLabels).Observe(duration)
		}

	} else if params["httpInOut"] == "httpOut" {
		httpOutDurationPromHistogram.With(promLabels).Observe(duration)
		httpOutPromCounter.With(promLabels).Inc()

		if promLabels["type"] == "httpapi" {
			klog.Infoln("httpOutPromCounterMap: ", promLabels["root"])
			httpOutPromCounterMap[promLabels["root"]].With(promLabels).Inc()
			httpOutDurationPromHistogramMap[promLabels["root"]].With(promLabels).Observe(duration)
		} else if promLabels["type"] == "flow" {
			httpOutPromCounterFlowMap[promLabels["root"]].With(promLabels).Inc()
			httpOutDurationPromHistogramFlowMap[promLabels["root"]].With(promLabels).Observe(duration)
		}
	} else {
		str := "httpInOut参数配置错误！"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusBadRequest
	}
	return nil, http.StatusOK
}

func promInitData() {

	//Init total counter and histogram
	httpInPromCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_in",
			Help: "api hub http in counters",
		},
		[]string{"code", "child", "root", "type"},
	)
	httpOutPromCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_out",
			Help: "api hub http out counters",
		},
		[]string{"code", "child", "root", "type"},
	)
	httpInDurationPromHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_in_duration_100ms",
			Help:    "apihub http in latency distributions.",
			Buckets: prometheus.LinearBuckets(0, 1, 101), // bucket从0开始,间隔是100ms,一共101个
		},
		[]string{"code", "child", "root", "type"},
	)
	httpOutDurationPromHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_out_duration_100ms",
			Help:    "apihub http out latency distributions.",
			Buckets: prometheus.LinearBuckets(0, 1, 101),
		},
		[]string{"code", "child", "root", "type"},
	)
	prometheus.MustRegister(httpInPromCounter)
	prometheus.MustRegister(httpOutPromCounter)
	prometheus.MustRegister(httpInDurationPromHistogram)
	prometheus.MustRegister(httpOutDurationPromHistogram)

	httpInPromCounterMap = make(map[string]*prometheus.CounterVec)
	httpOutPromCounterMap = make(map[string]*prometheus.CounterVec)
	httpInDurationPromHistogramMap = make(map[string]*prometheus.HistogramVec)
	httpOutDurationPromHistogramMap = make(map[string]*prometheus.HistogramVec)

	httpInPromCounterFlowMap = make(map[string]*prometheus.CounterVec)
	httpOutPromCounterFlowMap = make(map[string]*prometheus.CounterVec)
	httpInDurationPromHistogramFlowMap = make(map[string]*prometheus.HistogramVec)
	httpOutDurationPromHistogramFlowMap = make(map[string]*prometheus.HistogramVec)

	//Init each httpapi counter and histogram
	for k, v := range util.DefaultConfMap.ApiMap {
		inCounter := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_in_" + v.Id,
				Help: "api hub http in counters",
			},
			[]string{"code", "child", "root", "type"},
		)
		httpInPromCounterMap[k] = inCounter

		outCounter := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_out_" + v.Id,
				Help: "api hub http out counters",
			},
			[]string{"code", "child", "root", "type"},
		)
		httpOutPromCounterMap[k] = outCounter

		inHistogram := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_in_" + v.Id + "_duration_100ms",
				Help:    "apihub http in latency distributions.",
				Buckets: prometheus.LinearBuckets(0, 1, 101), // bucket从0开始,间隔是100ms,一共101个
			},
			[]string{"code", "child", "root", "type"},
		)

		httpInDurationPromHistogramMap[k] = inHistogram

		outHistogram := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_out_" + v.Id + "_duration_100ms",
				Help:    "apihub http out latency distributions.",
				Buckets: prometheus.LinearBuckets(0, 1, 101), // bucket从0开始,间隔是100ms,一共101个
			},
			[]string{"code", "child", "root", "type"},
		)

		httpOutDurationPromHistogramMap[k] = outHistogram

		prometheus.MustRegister(httpInPromCounterMap[k])
		prometheus.MustRegister(httpInDurationPromHistogramMap[k])
		prometheus.MustRegister(httpOutPromCounterMap[k])
		prometheus.MustRegister(httpOutDurationPromHistogramMap[k])
	}

	for k, v := range util.DefaultConfMap.FlowMap {
		if k == "main" {
			klog.Infoln("promInitData flow main, skipped")
			continue
		}
		inCounter := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "flow_in_" + v.Name,
				Help: "api hub flow in counters",
			},
			[]string{"code", "child", "root", "type"},
		)
		httpInPromCounterFlowMap[k] = inCounter

		outCounter := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "flow_out_" + v.Name,
				Help: "api hub flow out counters",
			},
			[]string{"code", "child", "root", "type"},
		)
		httpOutPromCounterFlowMap[k] = outCounter

		inHistogram := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "flow_in_" + v.Name + "_duration_100ms",
				Help:    "apihub flow in latency distributions.",
				Buckets: prometheus.LinearBuckets(0, 1, 101), // bucket从0开始,间隔是100ms,一共101个
			},
			[]string{"code", "child", "root", "type"},
		)
		httpInDurationPromHistogramFlowMap[k] = inHistogram

		outHistogram := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "flow_out_" + v.Name + "_duration_100ms",
				Help:    "apihub flow out latency distributions.",
				Buckets: prometheus.LinearBuckets(0, 1, 101), // bucket从0开始,间隔是100ms,一共101个
			},
			[]string{"code", "child", "root", "type"},
		)
		httpOutDurationPromHistogramFlowMap[k] = outHistogram

		prometheus.MustRegister(httpInPromCounterFlowMap[k])
		prometheus.MustRegister(httpOutPromCounterFlowMap[k])
		prometheus.MustRegister(httpInDurationPromHistogramFlowMap[k])
		prometheus.MustRegister(httpOutDurationPromHistogramFlowMap[k])
	}

	klog.Infoln("promInitData: OK")
}

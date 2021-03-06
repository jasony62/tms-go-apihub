package apis

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	klog "k8s.io/klog/v2"
)

var httpInPromCounter *prometheus.CounterVec
var httpOutPromCounter *prometheus.CounterVec
var httpInDurationPromHistogram *prometheus.HistogramVec
var httpOutDurationPromHistogram *prometheus.HistogramVec

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
	return nil, 200
}

func promStartRun(address string) {
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
			Name:    "http_in_duration_sec",
			Help:    "apihub http in latency distributions.",
			Buckets: prometheus.LinearBuckets(0, 1, 11), // bucket从0开始,间隔是1,一共11个
		},
		[]string{"code", "child", "root", "type"},
	)
	httpOutDurationPromHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_out_duration_sec",
			Help:    "apihub http out latency distributions.",
			Buckets: prometheus.LinearBuckets(0, 1, 11),
		},
		[]string{"code", "child", "root", "type"},
	)
	prometheus.MustRegister(httpInPromCounter)
	prometheus.MustRegister(httpOutPromCounter)
	prometheus.MustRegister(httpInDurationPromHistogram)
	prometheus.MustRegister(httpOutDurationPromHistogram)

	// Expose the registered metrics via HTTP
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	go func() {
		klog.Infoln("Listen and Serve ", address)
		if err := http.ListenAndServe(address, nil); err != nil {
			klog.Fatal("Error in ListenAndServe: %v", err)
		}
	}()
}

func getPromLabels(params map[string]string) map[string]string {
	return prometheus.Labels{
		"code": params["code"],
		"child": params["child"],
		"root": params["root"],
		"type": params["type"]}
}

func promHttpCounterInc(stack *hub.Stack, params map[string]string) (interface{}, int) {
	promLabels := getPromLabels(params)

	klog.Infoln("promLabels: ", promLabels)
	val := params["duration"]
	duration, err := strconv.ParseFloat(val, 64)
	if err != nil {
		klog.Errorln("解析http out duration失败, err: ", err)
		return nil, 400
	}
	if params["httpInOut"] == "httpIn" {
		httpInDurationPromHistogram.With(promLabels).Observe(duration)
		httpInPromCounter.With(promLabels).Inc()
	} else if params["httpInOut"] == "httpOut" {
		httpOutDurationPromHistogram.With(promLabels).Observe(duration)
		httpOutPromCounter.With(promLabels).Inc()
	} else {
		klog.Errorln("httpInOut参数配置错误！")
		return nil, 400
	}
	return nil, 200
}

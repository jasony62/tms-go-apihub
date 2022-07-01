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
			Name: "apihub_http_in",
			Help: "api hub http in counters",
		},
		[]string{"code", "duration", "id", "msg", "name", "root", "start", "type", "uuid"},
	)
	httpOutPromCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "apihub_http_out",
			Help: "api hub http out counters",
		},
		[]string{"code", "duration", "id", "msg", "name", "root", "start", "type", "uuid"},
	)
	httpInDurationPromHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "apihub_http_in_duration_sec",
			Help:    "apihub http in latency distributions.",
			Buckets: prometheus.LinearBuckets(0, 1, 11), // bucket从0开始,间隔是1,一共11个
		},
		[]string{"code", "id", "msg", "name", "root", "start", "type", "uuid"},
	)
	klog.Infoln("111!")
	httpOutDurationPromHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "apihub_http_out_duration_sec",
			Help:    "apihub http out latency distributions.",
			Buckets: prometheus.LinearBuckets(0, 1, 11),
		},
		[]string{"code", "id", "msg", "name", "root", "start", "type", "uuid"},
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
		"code":  params["code"],
		"id":    params["id"],
		"msg":   params["msg"],
		"name":  params["name"],
		"root":  params["root"],
		"start": params["start"],
		"type":  params["type"],
		"uuid":  params["uuid"]}
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
		promLabels["duration"] = params["duration"]
		httpInPromCounter.With(promLabels).Inc()
	} else if params["httpInOut"] == "httpOut" {
		httpOutDurationPromHistogram.With(promLabels).Observe(duration)
		promLabels["duration"] = params["duration"]
		httpOutPromCounter.With(promLabels).Inc()
	}
	return nil, 200
}

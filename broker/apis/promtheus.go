package apis

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/logger"
	"github.com/jasony62/tms-go-apihub/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpInPromCounter *prometheus.CounterVec
var httpOutPromCounter *prometheus.CounterVec
var httpInDurationPromHistogram *prometheus.HistogramVec
var httpOutDurationPromHistogram *prometheus.HistogramVec

func promStart(stack *hub.Stack, params map[string]string) (interface{}, int) {
	logger.LogS().Infoln("promStart!")
	host := params["host"]
	port := params["port"]

	if len(host) == 0 {
		host = "0.0.0.0"
	}
	if len(port) == 0 {
		port = "8000"
	}
	logger.LogS().Infoln("promStart: host: ", host, " port:", port)

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
		logger.LogS().Infoln("Listen and Serve: ", address)
		if err := http.ListenAndServe(address, nil); err != nil {
			logger.LogS().Errorln("Error in ListenAndServe: %v", err)
			panic("ERROR: prometheus ListenAndServe Failed!")
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

	logger.LogS().Infoln("promLabels: ", promLabels)
	val := params["duration"]
	duration, err := strconv.ParseFloat(val, 64)
	if err != nil {
		str := "解析http out duration失败, err: " + err.Error()
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusBadRequest
	}

	if params["httpInOut"] == "httpIn" {
		httpInDurationPromHistogram.With(promLabels).Observe(duration)
		httpInPromCounter.With(promLabels).Inc()
	} else if params["httpInOut"] == "httpOut" {
		httpOutDurationPromHistogram.With(promLabels).Observe(duration)
		httpOutPromCounter.With(promLabels).Inc()
	} else {
		str := "httpInOut参数配置错误！"
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusBadRequest
	}
	return nil, http.StatusOK
}

func promInitData() {
	//Init total counter and histogram
	httpInPromCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_in_total",
			Help: "api hub http in counter",
		},
		[]string{"code", "child", "root", "type"},
	)
	httpOutPromCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_out_total",
			Help: "api hub http out counter",
		},
		[]string{"code", "child", "root", "type"},
	)
	httpInDurationPromHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_in_duration_second",
			Help:    "apihub http in latency distributionsin second.",
			Buckets: prometheus.ExponentialBuckets(1, 2, 8),
		},
		[]string{"code", "child", "root", "type"},
	)
	httpOutDurationPromHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_out_duration_second",
			Help:    "apihub http out latency distributions in second.",
			Buckets: prometheus.ExponentialBuckets(1, 2, 8),
		},
		[]string{"code", "child", "root", "type"},
	)
	prometheus.MustRegister(httpInPromCounter)
	prometheus.MustRegister(httpOutPromCounter)
	prometheus.MustRegister(httpInDurationPromHistogram)
	prometheus.MustRegister(httpOutDurationPromHistogram)
}

package unit

import (
	"time"

	"strconv"

	klog "k8s.io/klog/v2"
	// "github.com/jasony62/tms-go-apihub/hub"
	// "github.com/jasony62/tms-go-apihub/plugin"
	// "github.com/jasony62/tms-go-apihub/util"
)

var funcMap map[string](func() string) = map[string](func() string){"utc": utcTime}

func utcTime() (result string) {
	result = strconv.FormatInt(time.Now().Unix(), 10)
	klog.Infoln("time is ", result, "s")
	return
}

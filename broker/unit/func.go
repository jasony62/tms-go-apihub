package unit

import (
	"fmt"
	klog "k8s.io/klog/v2"
	"time"
	// "github.com/jasony62/tms-go-apihub/hub"
	// "github.com/jasony62/tms-go-apihub/plugin"
	// "github.com/jasony62/tms-go-apihub/util"
)

var funcMap map[string](func() string) = map[string](func() string){"utc": utcTime}

func utcTime() string {
	cur := time.Now()
	timestamp := cur.UnixNano() / 1000000000
	klog.Infoln("time is ", timestamp, "s")
	return fmt.Sprint(timestamp)
}

package core

import (
	"github.com/jasony62/tms-go-apihub/hub"
	klog "k8s.io/klog/v2"
)

func init() {
	klog.Infof("Core register apis\n")
	RegisterApis(map[string]hub.ApiHandler{"flowApi": runFlowApi,
		"scheduleApi": runScheduleApi,
	})
}

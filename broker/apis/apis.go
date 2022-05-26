package apis

import (
	"github.com/jasony62/tms-go-apihub/core"
	"github.com/jasony62/tms-go-apihub/hub"
	klog "k8s.io/klog/v2"
)

func init() {
	klog.Infof("Register apis\n")
	core.RegisterTasks(map[string]hub.TaskHandler{"httpApi": runHttpApi,
		"httpResponse":         httpResponse,
		"checkStringsEqual":    checkStringsEqual,
		"checkStringsNotEqual": checkStringsNotEqual,
		"createJson":           createJson,
		"createHtml":           createHtml,
	})
}

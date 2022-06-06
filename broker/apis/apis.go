package apis

import (
	"github.com/jasony62/tms-go-apihub/core"
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
	klog "k8s.io/klog/v2"
)

func init() {
	klog.Infof("APIs register apis\n")
	core.RegisterApis(map[string]hub.ApiHandler{"httpApi": runHttpApi,
		"httpResponse":         httpResponse,
		"checkStringsEqual":    checkStringsEqual,
		"checkStringsNotEqual": checkStringsNotEqual,
		"createJson":           createJson,
		"createHtml":           createHtml,
		"loadConf":             util.LoadConf,
		"apiGateway":           apiGateway,
		"checkRight":           checkRight,
		"storageStore":         storageStore,
		"storageLoad":          storageLoad,
		"confValidator":        confValidator,
	})
}

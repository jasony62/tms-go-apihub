package apis

import (
	"github.com/Sheng-ZM/tms-go-apihub/broker//core"
	"github.com/Sheng-ZM/tms-go-apihub/broker//hub"
	klog "k8s.io/klog/v2"
)

func init() {
	klog.Infof("APIs register apis\n")
	core.RegisterApis(map[string]hub.ApiHandler{
		"httpApi":               runHttpApi,
		"httpResponse":          httpResponse,
		"checkStringsEqual":     checkStringsEqual,
		"checkStringsNotEqual":  checkStringsNotEqual,
		"createJson":            createJson,
		"createHtml":            createHtml,
		"loadConf":              loadConf,
		"apiGateway":            apiGateway,
		"fillBaseInfo":          fillBaseInfo,
		"setDefaultAccessRight": setDefaultAccessRight,
		"checkRight":            checkRight,
		"storageStore":          storageStore,
		"storageLoad":           storageLoad,
		"confValidator":         confValidator,
		"downloadConf":          downloadConf,
		"decompressZip":         decompressZip,
		"dump":                  dump,
		"promStart":             promStart,
		"promHttpCounterInc":    promHttpCounterInc,
		"logToFile":             logToFile,
	})
}

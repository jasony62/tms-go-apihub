package apis

import (
	"net/http"

	"github.com/Sheng-ZM/tms-go-apihub/broker//hub"
	"github.com/Sheng-ZM/tms-go-apihub/broker//util"
	klog "k8s.io/klog/v2"
)

func checkStringsEqual(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		str := "checkStringsEqual缺少参数"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	for k, v := range params {
		if k != v {
			str := "checkStringsEqual检查错误"
			klog.Errorln(stack.BaseString, str)
			return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
		}
	}
	return nil, http.StatusOK
}

func checkStringsNotEqual(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		str := "checkStringsNotEqual缺少参数"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	for k, v := range params {
		if k == v {
			str := "checkStringsNotEqual检查错误"
			klog.Errorln(stack.BaseString, str)
			return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
		}
	}
	return nil, http.StatusOK
}

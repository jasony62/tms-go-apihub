package apis

import (
	"net/http"

	"github.com/jasony62/tms-go-apihub/broker/hub"
	"github.com/jasony62/tms-go-apihub/broker/util"
	klog "k8s.io/klog/v2"
)

func createJson(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		str := "createJson,缺少参数"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	key, OK := params["key"]
	if !OK {
		str := "createJson,key为空"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}
	tmp := stack.Heap[hub.HeapOriginName].(map[string]interface{})
	result := tmp[key]
	delete(tmp, key)
	return result, http.StatusOK
}

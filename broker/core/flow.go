package core

import (
	"net/http"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/broker/hub"
	"github.com/jasony62/tms-go-apihub/broker/util"
)

func handleOneApi(stack *hub.Stack, apiDef *hub.ApiDef, private string) (result interface{}, ret int) {
	if len(apiDef.Command) > 0 {
		result, ret = ApiRun(stack, apiDef, private, false)
	} else {
		result, ret = nil, http.StatusInternalServerError
	}
	return
}

func runFlow(stack *hub.Stack, name string, private string) (result interface{}, ret int) {
	var code int
	var lastResult string
	flowDef, ok := util.FindFlowDef(name)
	if !ok || flowDef == nil {
		str := "获得Flow定义失败：" + name
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorCoreId, str, nil), http.StatusForbidden
	}

	for i := range flowDef.Steps {
		apiDef := flowDef.Steps[i]

		result, code = handleOneApi(stack, &apiDef, private)
		if code != http.StatusOK {
			str := "运行API：" + apiDef.Name + "失败"
			klog.Errorln(stack.BaseString, str)
			return util.CreateTmsError(hub.TmsErrorCoreId, str, nil), code
		}

		if len(apiDef.ResultKey) > 0 {
			stack.Heap[apiDef.ResultKey] = result
			lastResult = apiDef.ResultKey
		}
	}

	if len(lastResult) > 0 {
		return stack.Heap[lastResult], http.StatusOK
	} else {
		return nil, http.StatusOK
	}
}

func runFlowApi(stack *hub.Stack, params map[string]string) (interface{}, int) {
	name, OK := params["name"]
	if !OK {
		str := "缺少flow名称"
		klog.Errorln(stack.BaseString, str)
		return nil, http.StatusForbidden
	}
	private := params["private"]

	return runFlow(stack, name, private)
}

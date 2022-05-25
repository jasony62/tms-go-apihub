package core

import (
	"net/http"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
)

func handleOneTask(stack *hub.Stack, apiDef *hub.ApiDef) (result interface{}, ret int) {
	if len(apiDef.Command) > 0 {
		return ApiRun(stack, apiDef)
	}
	return nil, 500
}

func RunFlow(stack *hub.Stack) (result interface{}, ret int) {
	var code int
	flowDef, err := util.FindFlowDef(stack, stack.ChildName)
	if flowDef == nil {
		klog.Errorln("获得Flow定义失败：", err)
		panic(err)
	}

	for i := range flowDef.Steps {
		apiDef := flowDef.Steps[i]

		result, code = handleOneTask(stack, &apiDef)
		if code != http.StatusOK {
			return nil, code
		}

		if len(apiDef.ResultKey) > 0 {
			stack.StepResult[apiDef.ResultKey] = result
		}
	}

	return nil, http.StatusOK
}

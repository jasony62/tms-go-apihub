package core

import (
	"net/http"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
)

func handleOneApi(stack *hub.Stack, apiDef *hub.ApiDef, private string) (result interface{}, ret int) {
	if len(apiDef.Command) > 0 {
		klog.Infoln("=========执行API name："+apiDef.Name+", command:"+apiDef.Command, " args:", apiDef.Args)
		result, ret = ApiRun(stack, apiDef, private)
	} else {
		result, ret = nil, http.StatusInternalServerError
	}
	klog.Infoln("=========完成API name："+apiDef.Name+", command:"+apiDef.Command, " result:", result)
	return
}

func runFlow(stack *hub.Stack, name string, private string) (result interface{}, ret int) {
	var code int
	var lastResult string
	flowDef, err := util.FindFlowDef(name)
	if flowDef == nil {
		klog.Errorln("获得Flow定义失败：", err)
		panic(err)
	}

	for i := range flowDef.Steps {
		apiDef := flowDef.Steps[i]

		result, code = handleOneApi(stack, &apiDef, private)
		if code != http.StatusOK {
			klog.Errorln("运行API：" + apiDef.Name + "失败")
			return nil, code
		}

		if len(apiDef.ResultKey) > 0 {
			stack.StepResult[apiDef.ResultKey] = result
			lastResult = apiDef.ResultKey
		}
	}

	if len(lastResult) > 0 {
		return stack.StepResult[lastResult], http.StatusOK
	} else {
		return nil, http.StatusOK
	}
}

func runFlowApi(stack *hub.Stack, params map[string]string) (interface{}, int) {
	name, OK := params["name"]
	if !OK {
		str := "缺少flow名称"
		klog.Errorln(str)
		panic(str)
	}
	private := params["private"]

	return runFlow(stack, name, private)
}

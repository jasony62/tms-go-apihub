package flow

import (
	"net/http"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/api"
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/unit"
	"github.com/jasony62/tms-go-apihub/util"
)

func fillOrigin(stack *hub.Stack, parameters *[]hub.OriginDefParam) {
	var value string
	origin := stack.StepResult[hub.OriginName].(map[string]interface{})

	for _, parameter := range *parameters {
		if len(parameter.Value) > 0 {
			value = parameter.Value
		} else {
			value = unit.GetParameterValue(stack, nil, parameter.From)
		}

		oldValue, isOk := origin[parameter.Name]
		if isOk {
			klog.Infoln("replace ", parameter.Name, " from ", oldValue, " to ", value)
		}
		origin[parameter.Name] = value
	}
}

func Run(stack *hub.Stack) (interface{}, int) {
	var lastResultKey string
	flowDef, err := unit.FindFlowDef(stack, stack.ChildName)

	if flowDef == nil {
		klog.Errorln("获得Flow定义失败：", err)
		panic(err)
	}

	for _, step := range flowDef.Steps {
		if step.Api != nil && len(step.Api.Id) > 0 {

			if step.Api.Parameters != nil && len(*step.Api.Parameters) > 0 {
				// 根据flow的定义改写origin
				fillOrigin(stack, step.Api.Parameters)
			}

			// 调用api
			stack.ChildName = step.Api.Id
			jsonOutRspBody, _ := api.Run(stack)

			// 在上下文中保存结果
			if len(step.ResultKey) > 0 {
				stack.StepResult[step.ResultKey] = jsonOutRspBody
			}
		} else if step.Response != nil && len(step.ResultKey) > 0 {
			// 处理响应结果
			stack.StepResult[step.ResultKey] = util.Json2Json(stack.StepResult, step.Response.Json)
		}

		lastResultKey = step.ResultKey
	}

	return stack.StepResult[lastResultKey], http.StatusOK
}

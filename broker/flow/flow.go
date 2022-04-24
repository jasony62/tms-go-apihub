package flow

import (
	"log"
	"net/http"

	"github.com/jasony62/tms-go-apihub/api"
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/unit"
	"github.com/jasony62/tms-go-apihub/util"
)

func Run(stack *hub.Stack) (interface{}, int) {
	var lastResultKey string
	flowDef := stack.FlowDef
	for _, step := range flowDef.Steps {
		//stack.CurrentStep = &step
		if step.Api != nil && len(step.Api.Id) > 0 {
			// 执行API并记录结果
			//apiDef, err := unit.FindApiDef(stack, "", step.Api.Id)
			apiDef, err := unit.FindApiDef(stack, step.Api.Id)

			if apiDef == nil {
				log.Panic("获得API", step.Api.Id, "定义失败：", err)
			}
			// 根据flow的定义改写api定义
			if step.Api.Parameters != nil && len(*step.Api.Parameters) > 0 {
				unit.RewriteApiDefInFlow(apiDef, step.Api)
			}
			// 调用api
			stack.ApiDef = apiDef

			api.Relay(stack, step.ResultKey)
		} else if step.Response != nil {
			// 处理响应结果
			outBodyRules := step.Response.Json
			jsonRspBody := util.Json2Json(stack.StepResult, outBodyRules)
			stack.StepResult[step.ResultKey] = jsonRspBody
		}

		lastResultKey = step.ResultKey
	}

	return stack.StepResult[lastResultKey], http.StatusOK
}

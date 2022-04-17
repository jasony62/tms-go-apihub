package flow

import (
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
		stack.CurrentStep = &step
		if step.Api != nil && len(step.Api.Id) > 0 {
			// 执行API并记录结果
			apiDef, _ := unit.FindApiDef(stack, "", step.Api.Id)
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

	/*获得最后一个任务的结果*/
	jsonOutBody := stack.StepResult[lastResultKey]

	return jsonOutBody, http.StatusOK
}

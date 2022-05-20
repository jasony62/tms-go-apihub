package flow

import (
	"net/http"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/api"
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/unit"
	"github.com/jasony62/tms-go-apihub/util"
)

type concurrentFlowIn struct {
	step *hub.FlowStepDef
}

type concurrentFlowOut struct {
	step   *hub.FlowStepDef
	result interface{}
}

func copyFlowStack(src *hub.Stack) *hub.Stack {
	stepResult := make(map[string]interface{})
	for k, v := range src.StepResult {
		stepResult[k] = v
	}

	//avoid vars race conditions
	return &hub.Stack{
		GinContext: src.GinContext,
		RootName:   src.RootName,
		ChildName:  "",
		StepResult: stepResult,
	}
}

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

func handleOneApi(stack *hub.Stack, step *hub.FlowStepDef) (result interface{}) {
	if step.Api != nil && len(step.Api.Id) > 0 {
		if step.Api.Parameters != nil && len(*step.Api.Parameters) > 0 {
			// 根据flow的定义改写origin
			fillOrigin(stack, step.Api.Parameters)
		}

		// 调用api
		stack.ChildName = step.Api.Id
		jsonOutRspBody, _ := api.Run(stack, step.Api.PrivateName)

		// 在上下文中保存结果
		if len(step.ResultKey) > 0 {
			result = jsonOutRspBody
		}
	} else if step.Response != nil && step.Response.From != nil {
		// 处理响应结果
		klog.Infoln("handleOneApi响应文本格式", step.Response.Type)
		var rules interface{}
		if step.Response.Type == "json" {
			rules = step.Response.From.Template
		} else if step.Response.Type == "html" {
			rules = step.Response.From.Name
		}
		result = util.Json2Json(stack.StepResult, rules)
		if result == nil {
			klog.Infoln("get final result failed：", rules, "\r\n", stack.StepResult, "\r\n", result)
		} else {
			klog.Infoln("get final result：", result)
		}
	}
	return result
}

func concurrentFlowWorker(stack *hub.Stack, steps chan concurrentFlowIn, out chan concurrentFlowOut) {
	for step := range steps {
		result := handleOneApi(copyFlowStack(stack), step.step)
		out <- concurrentFlowOut{step: step.step, result: result}
	}
}

func waitConcurrentFlowResult(stack *hub.Stack, out chan concurrentFlowOut, counter int) (lastKey string) {
	results := make(map[string]interface{}, counter)
	for counter > 0 {
		//等待结果
		result := <-out
		key := result.step.ResultKey
		if len(key) > 0 {
			results[key] = result.result
			lastKey = key
		}
		counter--
	}
	//防止并发读写crash
	for k, v := range results {
		stack.StepResult[k] = v
	}

	//由于并行，最后的结果并不确定，所以并行的返回结果不是固定的，因此当需要返回值时，最后一个应该是非并行的
	return lastKey
}

func Run(stack *hub.Stack) (interface{}, string, int) {
	var lastResultKey string
	var lastTypeKey string
	var counter int
	var in chan concurrentFlowIn
	var out chan concurrentFlowOut
	var result interface{}

	flowDef, err := unit.FindFlowDef(stack, stack.ChildName)
	if flowDef == nil {
		klog.Errorln("获得Flow定义失败：", err)
		panic(err)
	}

	if flowDef.ConcurrentNum > 1 {
		in = make(chan concurrentFlowIn, len(flowDef.Steps))
		defer close(in)
		out = make(chan concurrentFlowOut, len(flowDef.Steps))
		defer close(out)
		for i := 0; i < flowDef.ConcurrentNum; i++ {
			go concurrentFlowWorker(stack, in, out)
		}
	}

	lastTypeKey = "json" //默认类型为json

	for i := range flowDef.Steps {
		step := flowDef.Steps[i]
		if flowDef.ConcurrentNum > 1 {
			if step.Concurrent {
				in <- concurrentFlowIn{step: &step}
				counter++
				continue
			} else {
				//避免并发读写ResultKey
				if counter > 0 {
					lastResultKey = waitConcurrentFlowResult(stack, out, counter)
					counter = 0
				}
			}
		}

		result = handleOneApi(stack, &step)
		if len(step.ResultKey) > 0 {
			stack.StepResult[step.ResultKey] = result
			lastResultKey = step.ResultKey
			if step.Response != nil {
				lastTypeKey = step.Response.Type
			}
		}
	}

	//当最后一个step也是并行，等待全部执行完
	if counter > 0 {
		lastResultKey = waitConcurrentFlowResult(stack, out, counter)
	}

	//由于并行，最后的结果并不确定，所以并行的返回结果不是固定的，因此当需要返回值时，最后一个应该是非并行的
	return stack.StepResult[lastResultKey], lastTypeKey, http.StatusOK
}

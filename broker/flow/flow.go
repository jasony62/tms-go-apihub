package flow

import (
	"net/http"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/task"
	"github.com/jasony62/tms-go-apihub/unit"
	"github.com/jasony62/tms-go-apihub/util"
)

type concurrentFlowIn struct {
	taskDef *hub.TaskDef
}
type concurrentFlowOut struct {
	taskDef *hub.TaskDef
	result  interface{}
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
		StepResult: stepResult,
	}
}

func fillOrigin(stack *hub.Stack, parameters *[]hub.OriginDefParam) {
	var value string
	origin := stack.StepResult[hub.OriginName].(map[string]interface{})

	for _, parameter := range *parameters {
		value = unit.GetParameterValue(stack, nil, parameter.From)
		oldValue, isOk := origin[parameter.Name]
		if isOk {
			klog.Infoln("replace ", parameter.Name, " from ", oldValue, " to ", value)
		}
		origin[parameter.Name] = value
	}
}

func handleOneTask(stack *hub.Stack, taskDef *hub.TaskDef) (result interface{}, ret int) {
	if len(taskDef.Command) > 0 {
		return task.Run(stack, taskDef)
	} else if taskDef.Response != nil && taskDef.Response.From != nil {
		// 处理响应结果
		klog.Infoln("handleOneApi响应文本格式", taskDef.Response.Type)
		var rules interface{}
		if taskDef.Response.Type == "json" {
			rules = taskDef.Response.From.Json
		} else if taskDef.Response.Type == "html" {
			rules = taskDef.Response.From.Content
		}
		result = util.Json2Json(stack.StepResult, rules)
		if result == nil {
			klog.Infoln("get final result failed：", rules, "\r\n", stack.StepResult, "\r\n", result)
		} else {
			klog.Infoln("get final result：", result)
		}
	}
	return result, 200
}

func concurrentFlowWorker(stack *hub.Stack, tasks chan concurrentFlowIn, out chan concurrentFlowOut) {
	for taskDef := range tasks {
		result, _ := handleOneTask(copyFlowStack(stack), taskDef.taskDef)
		out <- concurrentFlowOut{taskDef: taskDef.taskDef, result: result}
	}
}

func waitConcurrentTaskResult(stack *hub.Stack, out chan concurrentFlowOut, counter int) (lastKey string) {
	results := make(map[string]interface{}, counter)
	for counter > 0 {
		//等待结果
		result := <-out
		key := result.taskDef.ResultKey
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
	var code int

	flowDef, err := unit.FindFlowDef(stack, stack.ChildName)
	if flowDef == nil {
		klog.Errorln("获得Flow定义失败：", err)
		panic(err)
	}

	if flowDef.ConcurrentNum > 1 {
		in = make(chan concurrentFlowIn, len(flowDef.Tasks))
		defer close(in)
		out = make(chan concurrentFlowOut, len(flowDef.Tasks))
		defer close(out)
		for i := 0; i < flowDef.ConcurrentNum; i++ {
			go concurrentFlowWorker(stack, in, out)
		}
	}

	lastTypeKey = "json" //默认类型为json

	for i := range flowDef.Tasks {
		taskDef := flowDef.Tasks[i]
		if flowDef.ConcurrentNum > 1 {
			if taskDef.Concurrent {
				in <- concurrentFlowIn{taskDef: &taskDef}
				counter++
				continue
			} else {
				//避免并发读写ResultKey
				if counter > 0 {
					lastResultKey = waitConcurrentTaskResult(stack, out, counter)
					counter = 0
				}
			}
		}

		result, code = handleOneTask(stack, &taskDef)
		if code != 200 {
			return nil, "", code
		}

		if len(taskDef.ResultKey) > 0 {
			stack.StepResult[taskDef.ResultKey] = result
			lastResultKey = taskDef.ResultKey
			if taskDef.Response != nil {
				lastTypeKey = taskDef.Response.Type
			}
		}
	}

	//当最后一个step也是并行，等待全部执行完
	if counter > 0 {
		lastResultKey = waitConcurrentTaskResult(stack, out, counter)
	}

	//由于并行，最后的结果并不确定，所以并行的返回结果不是固定的，因此当需要返回值时，最后一个应该是非并行的
	return stack.StepResult[lastResultKey], lastTypeKey, http.StatusOK
}

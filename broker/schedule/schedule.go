package schedule

import (
	"strconv"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/api"
	"github.com/jasony62/tms-go-apihub/flow"
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/unit"
)

func generateStepResult(stack *hub.Stack, parameters *[]hub.OriginDefParam) interface{} {
	var value string
	result := make(map[string]interface{}, len(*parameters))
	for _, parameter := range *parameters {
		if len(parameter.Value) > 0 {
			value = parameter.Value
		} else {
			value = unit.GetParameterValue(stack, nil, parameter.From)
		}
		result[parameter.Name] = value
	}
	return result
}

func copyStack(src *hub.Stack, task *hub.ScheduleTaskDef) *hub.Stack {
	result := hub.Stack{
		GinContext: src.GinContext,
		Name:       task.Commond,
	}

	if task.Parameters != nil {
		result.StepResult = map[string]interface{}{hub.OriginName: generateStepResult(src, task.Parameters)}
	} else {
		result.StepResult = map[string]interface{}{hub.OriginName: src.StepResult[hub.OriginName]}
	}

	oriLoop := src.StepResult["loop"].(map[string]int)
	loop := make(map[string]int, len(oriLoop))
	for index, element := range oriLoop {
		loop[index] = element
	}
	result.StepResult["loop"] = loop
	return &result
}

func handleSwitchTask(stack *hub.Stack, task *hub.ScheduleTaskDef) (interface{}, int) {
	key := unit.GetParameterValue(stack, nil, &task.Key)

	if len(key) == 0 {
		err := "invalid switch key"
		klog.Errorln(err)
		panic(err)
	}

	for _, item := range *task.Cases {
		if item.Value == key {
			return handleTasks(stack, item.Tasks)
		}
	}
	return nil, 500
}

func handleLoopTask(stack *hub.Stack, task *hub.ScheduleTaskDef) (interface{}, int) {
	var result interface{}
	keyStr := unit.GetParameterValue(stack, nil, &task.Key)

	if len(keyStr) == 0 {
		err := "invalid loop key"
		klog.Errorln(err)
		panic(err)
	}
	max, _ := strconv.Atoi(keyStr)
	loopResult := make([]interface{}, max)
	if len(task.ResultKey) > 0 {
		stack.StepResult[task.ResultKey] = loopResult
	}
	for i := 0; i < max; i++ {
		loop := stack.StepResult["loop"].(map[string]int)
		loop[task.Name] = i
		result, _ = handleTasks(stack, task.Tasks)
		loopResult[i] = result
	}
	return loopResult, 200
}

func handleControlTask(stack *hub.Stack, task *hub.ScheduleTaskDef) (interface{}, int) {
	switch task.Commond {
	case "switch":
		if task.Cases != nil {
			return handleSwitchTask(stack, task)
		} else {
			err := "No switch cases"
			klog.Errorln(err)
			panic(err)
		}
	case "loop":
		return handleLoopTask(stack, task)
	default:
		err := "don't support command " + task.Type
		klog.Errorln(err)
		panic(err)
	}
}

func handleFlowTask(stack *hub.Stack, task *hub.ScheduleTaskDef) (result interface{}, status int) {
	tmpStack := copyStack(stack, task)

	// 执行编排
	result, status = flow.Run(tmpStack)

	if len(task.ResultKey) > 0 {
		stack.StepResult[task.ResultKey] = result
	}
	return
}

func handleApiTask(stack *hub.Stack, task *hub.ScheduleTaskDef) (result interface{}, status int) {
	tmpStack := copyStack(stack, task)

	// 执行API
	result, status = api.Run(tmpStack)

	if len(task.ResultKey) > 0 {
		stack.StepResult[task.ResultKey] = result
	}
	return
}

func handleTasks(stack *hub.Stack, tasks *[]hub.ScheduleTaskDef) (result interface{}, status int) {
	for _, task := range *tasks {
		if len(task.Type) > 0 && len(task.Commond) > 0 {
			switch task.Type {
			case "control":
				result, status = handleControlTask(stack, &task)
			case "flow":
				result, status = handleFlowTask(stack, &task)
			case "api":
				result, status = handleApiTask(stack, &task)
			default:
				err := "don't support type " + task.Type
				klog.Errorln(err)
				panic(err)
			}
		}
	}
	return result, status
}

func Run(stack *hub.Stack) (interface{}, int) {
	scheduleDef, err := unit.FindScheduleDef(stack, stack.Name)
	if scheduleDef == nil || scheduleDef.Tasks == nil {
		klog.Errorln("获得Schedule定义失败：", err)
		panic(err)
	}
	stack.StepResult["loop"] = make(map[string]int)
	return handleTasks(stack, scheduleDef.Tasks)
}

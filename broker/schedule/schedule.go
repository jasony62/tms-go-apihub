package schedule

import (
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
	var value interface{}
	if task.Parameters != nil {
		value = generateStepResult(src, task.Parameters)
	} else {
		value = src.StepResult[hub.OriginName]
	}

	return &hub.Stack{GinContext: src.GinContext,
		StepResult: map[string]interface{}{hub.OriginName: value},
		Name:       task.Commond,
	}
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
			return handleTasks(stack, *item.Tasks)
		}
	}
	return nil, 500
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

	// 执行编排
	result, status = api.Run(tmpStack)

	if len(task.ResultKey) > 0 {
		stack.StepResult[task.ResultKey] = result
	}
	return
}

func handleTasks(stack *hub.Stack, tasks []hub.ScheduleTaskDef) (result interface{}, status int) {
	for _, task := range tasks {
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
	if scheduleDef == nil {
		klog.Errorln("获得Schedule定义失败：", err)
		panic(err)
	}
	return handleTasks(stack, scheduleDef.Tasks)
}

package schedule

import (
	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/flow"
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/unit"
)

func copyStack(src *hub.Stack) *hub.Stack {
	return &hub.Stack{GinContext: src.GinContext,
		StepResult: make(map[string]interface{}),
	}
}

func handleSwitchTask(stack *hub.Stack, task *hub.ScheduleTaskDef) (interface{}, int) {
	key := unit.GetParameterValue(stack, &task.Key)

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

func generateStepResult(stack *hub.Stack, parameters *[]hub.ScheduleDefParam) interface{} {
	var value string
	result := make(map[string]string, len(*parameters))
	for _, parameter := range *parameters {
		if len(parameter.Value) > 0 {
			value = parameter.Value
		} else {
			value = unit.GetParameterValue(stack, parameter.From)
		}
		result[parameter.Name] = value
	}
	return result
}

func handleFlowTask(stack *hub.Stack, task *hub.ScheduleTaskDef) (result interface{}, status int) {
	var err error
	tmpStack := copyStack(stack)
	tmpStack.FlowDef, err = unit.FindFlowDef(stack, task.Commond)

	if tmpStack.FlowDef == nil {
		str := "获得Flow定义失败：" + err.Error()
		klog.Errorln(str)
		panic(str)
	}

	if task.Parameters != nil {
		tmpStack.StepResult["origin"] = generateStepResult(stack, task.Parameters)
	} else {
		tmpStack.StepResult["origin"] = stack.StepResult["origin"]
	}

	// 执行编排
	result, status = flow.Run(tmpStack)

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
				err := "don't support api"
				klog.Errorln(err)
				panic(err)
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
	return handleTasks(stack, stack.ScheduleDef.Tasks)
}

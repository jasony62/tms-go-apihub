package schedule

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/jasony62/tms-go-apihub/flow"
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/unit"
)

func copyStack(src *hub.Stack) *hub.Stack {
	// 收到的数据
	stack := new(hub.Stack)
	stack.GinContext = src.GinContext
	stack.StepResult = make(map[string]interface{})
	return stack
}

func handleSwitchTask(stack *hub.Stack, task *hub.ScheduleTaskDef) (interface{}, int) {
	key := unit.GetParameterValue(stack, task.Key.From, task.Key.Name, task.Key.Template)

	if len(key) == 0 {
		log.Panic("invalid switch key")
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
			log.Panic("No switch cases")
		}
	default:
		log.Panic("don't support command ", task.Type)
	}
	return nil, 500
}

func generateOrigin(stack *hub.Stack, parameters *[]hub.ScheduleDefParam) interface{} {
	var b bytes.Buffer
	var value string
	count := 0
	b.Grow(len(*parameters) * 64)
	b.WriteString("{")
	for _, parameter := range *parameters {
		if count > 0 {
			b.WriteString(",")
		}
		count++
		b.WriteString("\"")
		b.WriteString(parameter.Name)
		b.WriteString("\":\"")
		if len(parameter.Value) > 0 {
			value = parameter.Value
		} else {
			value = unit.GetParameterValue(stack, parameter.From.From, parameter.From.Name, parameter.From.Template)
		}
		b.WriteString(value)
		b.WriteString("\"")
	}
	b.WriteString("}")

	var v interface{}
	err := json.Unmarshal(b.Bytes(), &v)
	if err != nil {
		log.Panic("fail to generate new origin ", err)
	}
	return v
}

func handleFlowTask(stack *hub.Stack, task *hub.ScheduleTaskDef) (result interface{}, status int) {
	var err error
	tmpStack := copyStack(stack)
	tmpStack.FlowDef, err = unit.FindFlowDef(stack, task.Commond)

	if tmpStack.FlowDef == nil {
		log.Panic("获得Flow定义失败：", err)
		return nil, 500
	}

	if task.Parameters != nil {
		tmpStack.StepResult["origin"] = generateOrigin(stack, task.Parameters)
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
				log.Panic("don't support api")
			default:
				log.Panic("don't support type ", task.Type)
			}
		}
	}
	return result, status
}

func Run(stack *hub.Stack) (interface{}, int) {
	return handleTasks(stack, stack.ScheduleDef.Tasks)
}

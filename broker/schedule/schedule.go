package schedule

import (
	"strconv"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/api"
	"github.com/jasony62/tms-go-apihub/flow"
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/unit"
)

type concurrentLoopIn struct {
	index int
	stack *hub.Stack
	task  *[]hub.ScheduleTaskDef
}

type concurrentLoopOut struct {
	index  int
	result interface{}
}

type concurrentScheIn struct {
	task *hub.ScheduleTaskDef
}

type concurrentScheOut struct {
	task   *hub.ScheduleTaskDef
	result interface{}
}

func generateStepResult(stack *hub.Stack, parameters *[]hub.OriginDefParam) interface{} {
	var value string
	result := make(map[string]interface{}, len(*parameters))
	for _, parameter := range *parameters {
		value = unit.GetParameterValue(stack, nil, parameter.From)
		result[parameter.Name] = value
	}
	return result
}

func copyScheduleStack(src *hub.Stack, task *hub.ScheduleTaskDef) *hub.Stack {
	result := hub.Stack{
		GinContext: src.GinContext,
		RootName:   src.RootName,
		ChildName:  task.Name,
		StepResult: make(map[string]interface{}),
	}

	for k, v := range src.StepResult {
		switch k {
		case hub.OriginName:
			if task.Parameters != nil {
				result.StepResult[k] = generateStepResult(src, task.Parameters)
			} else {
				result.StepResult[k] = v
			}
		case hub.LoopName:
			oriLoop := src.StepResult[k].(map[string]int)
			loop := make(map[string]int, len(oriLoop))
			for index, element := range oriLoop {
				loop[index] = element
			}
			result.StepResult[k] = loop
		default:
			result.StepResult[k] = v
		}

	}

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
			return handleTasks(stack, item.Tasks, task.ConcurrentNum)
		}
	}
	return nil, 500
}

func concurrentLoopWorker(tasks chan concurrentLoopIn, out chan concurrentLoopOut) {
	for task := range tasks {
		result, _ := handleTasks(task.stack, task.task, 0)
		out <- concurrentLoopOut{index: task.index, result: result}
	}
}

func triggerConcurrentLoop(stack *hub.Stack, task *hub.ScheduleTaskDef, loopLength int, loop map[string]int, loopResult []interface{}) {
	var taskCount, msgCount int
	counter := loopLength

	if task.ConcurrentLoopNum > loopLength {
		taskCount = loopLength
		msgCount = loopLength
	} else {
		taskCount = task.ConcurrentLoopNum
		msgCount = taskCount * 2
		if msgCount > loopLength {
			msgCount = loopLength
		}
	}

	in := make(chan concurrentLoopIn, msgCount)
	defer close(in)
	out := make(chan concurrentLoopOut, msgCount)
	defer close(out)

	i := 0

	for ; i < msgCount; i++ {
		loop[task.ResultKey] = i
		in <- concurrentLoopIn{index: i, stack: copyScheduleStack(stack, task), task: task.Tasks}
	}

	for i = 0; i < taskCount; i++ {
		go concurrentLoopWorker(in, out)
	}

	i = msgCount
	for result := range out {
		loopResult[result.index] = result.result
		counter--
		if i < loopLength {
			loop[task.ResultKey] = i
			tmpStack := copyScheduleStack(stack, task)
			in <- concurrentLoopIn{index: i, stack: tmpStack, task: task.Tasks}
			i++
		} else {
			if counter == 0 {
				break
			}
		}
	}
}

func handleLoopTask(stack *hub.Stack, task *hub.ScheduleTaskDef) (interface{}, int) {
	var result interface{}
	keyStr := unit.GetParameterValue(stack, nil, &task.Key)

	if len(keyStr) == 0 {
		err := "invalid loop key"
		klog.Errorln(err)
		panic(err)
	}
	loopLength, _ := strconv.Atoi(keyStr)
	loopResult := make([]interface{}, loopLength)
	if !task.Concurrent && len(task.ResultKey) > 0 {
		stack.StepResult[task.ResultKey] = loopResult
	}

	loop := stack.StepResult[hub.LoopName].(map[string]int)
	if task.ConcurrentLoopNum > 1 && loopLength > 1 {
		triggerConcurrentLoop(stack, task, loopLength, loop, loopResult)
	} else {
		for i := 0; i < loopLength; i++ {
			loop[task.Name] = i
			result, _ = handleTasks(stack, task.Tasks, task.ConcurrentNum)
			loopResult[i] = result
		}
	}
	return loopResult, 200
}

func handleControlTask(stack *hub.Stack, task *hub.ScheduleTaskDef) (interface{}, int) {
	switch task.Name {
	case "switch":
		if task.Cases != nil {
			return handleSwitchTask(stack, task)
		} else {
			err := "No switch cases"
			klog.Errorln(err)
			panic(err)
		}
	case hub.LoopName:
		return handleLoopTask(stack, task)
	default:
		err := "don't support command " + task.Type
		klog.Errorln(err)
		panic(err)
	}
}

func handleFlowTask(stack *hub.Stack, task *hub.ScheduleTaskDef) (result interface{}, status int) {
	// 执行编排
	result, _, status = flow.Run(copyScheduleStack(stack, task))

	if !task.Concurrent && len(task.ResultKey) > 0 {
		stack.StepResult[task.ResultKey] = result
	}
	return
}

func handleApiTask(stack *hub.Stack, task *hub.ScheduleTaskDef) (result interface{}, status int) {
	// 执行API
	result, status = api.Run(copyScheduleStack(stack, task), task.Name, task.PrivateName)

	if !task.Concurrent && len(task.ResultKey) > 0 {
		stack.StepResult[task.ResultKey] = result
	}
	return
}
func handleOneTask(stack *hub.Stack, task *hub.ScheduleTaskDef) (result interface{}, status int) {
	if len(task.Type) > 0 && len(task.Name) > 0 {
		switch task.Type {
		case "control":
			result, status = handleControlTask(stack, task)
		case "flow":
			result, status = handleFlowTask(stack, task)
		case "api":
			result, status = handleApiTask(stack, task)
		default:
			err := "don't support type " + task.Type
			klog.Errorln(err)
			panic(err)
		}
	}
	return result, status
}

func concurrentScheWorker(stack *hub.Stack, tasks chan concurrentScheIn, out chan concurrentScheOut) {
	for task := range tasks {
		result, _ := handleOneTask(stack, task.task)
		out <- concurrentScheOut{task: task.task, result: result}
	}
}

func waitConcurrentScheResult(stack *hub.Stack, out chan concurrentScheOut, counter int) interface{} {
	results := make(map[string]interface{}, counter)
	var lastKey string
	for counter > 0 {
		//等待结果
		result := <-out
		key := result.task.ResultKey
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
	return results[lastKey]
}

func handleTasks(stack *hub.Stack, tasks *[]hub.ScheduleTaskDef, concurrentNum int) (result interface{}, status int) {
	var counter int
	var in chan concurrentScheIn
	var out chan concurrentScheOut

	if concurrentNum > 1 {
		/*假设所有的task都是并行的，多留buffer，提升性能*/
		in = make(chan concurrentScheIn, len(*tasks))
		defer close(in)
		out = make(chan concurrentScheOut, len(*tasks))
		defer close(out)
		for i := 0; i < concurrentNum; i++ {
			go concurrentScheWorker(stack, in, out)
		}
	}

	for index := range *tasks {
		task := &(*tasks)[index]
		if concurrentNum > 1 {
			if task.Concurrent {
				in <- concurrentScheIn{task: task}
				counter++
				continue
			} else {
				//避免并发读写ResultKey
				if counter > 0 {
					result = waitConcurrentScheResult(stack, out, counter)
					counter = 0
				}
			}
		}

		result, status = handleOneTask(stack, task)
	}

	//防止都是并行任务
	if counter > 0 {
		result = waitConcurrentScheResult(stack, out, counter)
	}
	return result, status
}

func Run(stack *hub.Stack) (interface{}, int) {
	scheduleDef, err := unit.FindScheduleDef(stack, stack.ChildName)
	if scheduleDef == nil || scheduleDef.Tasks == nil {
		klog.Errorln("获得Schedule定义失败：", err)
		panic(err)
	}
	stack.StepResult[hub.LoopName] = make(map[string]int)

	return handleTasks(stack, scheduleDef.Tasks, scheduleDef.ConcurrentNum)
}

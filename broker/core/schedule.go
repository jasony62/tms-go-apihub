package core

import (
	"net/http"
	"strconv"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/logger"
	"github.com/jasony62/tms-go-apihub/util"
)

type concurrentLoopIn struct {
	index int
	stack *hub.Stack
	task  *[]hub.ScheduleApiDef
}

type concurrentLoopOut struct {
	index  int
	result interface{}
}

type concurrentScheIn struct {
	task *hub.ScheduleApiDef
}

type concurrentScheOut struct {
	task   *hub.ScheduleApiDef
	result interface{}
}

func isNormalMode(task *hub.ScheduleApiDef) bool {
	mode := task.Mode
	return (mode != "concurrent") && (mode != "background")
}

func generateStepResult(stack *hub.Stack, args *[]hub.BaseParamDef) interface{} {
	result := make(map[string]interface{}, len(*args))
	for _, parameter := range *args {
		result[parameter.Name], _ = util.GetParameterRawValue(stack, nil, &parameter.Value)
	}
	return result
}

func copyScheduleStack(src *hub.Stack, task *hub.ScheduleApiDef) *hub.Stack {
	result := hub.Stack{
		GinContext: src.GinContext,
		Heap:       make(map[string]interface{}),
		BaseString: src.BaseString,
	}

	for k, v := range src.Heap {
		switch k {
		case hub.HeapOriginName:
			if task.Type == "api" && task.Api.Args != nil {
				result.Heap[k] = generateStepResult(src, task.Api.Args)
			} else {
				result.Heap[k] = v
			}
		case hub.HeapLoopName:
			oriLoop := src.Heap[k].(map[string]int)
			loop := make(map[string]int, len(oriLoop))
			for index, element := range oriLoop {
				loop[index] = element
			}
			result.Heap[k] = loop
		default:
			result.Heap[k] = v
		}

	}

	return &result
}

func handleSwitchTask(stack *hub.Stack, task *hub.ScheduleApiDef) (interface{}, int) {
	key, _ := util.GetParameterStringValue(stack, nil, &task.Control.Key)

	if len(key) == 0 {
		err := "invalid switch key"
		logger.LogS().Errorln(stack.BaseString, err)
		return util.CreateTmsError(hub.TmsErrorCoreId, err, nil), http.StatusInternalServerError
	}

	for _, item := range *task.Control.Cases {
		if item.Value == key {
			return handleTasks(stack, item.Steps, task.Control.ConcurrentNum)
		}
	}
	return util.CreateTmsError(hub.TmsErrorCoreId, "No task control", nil), http.StatusInternalServerError
}

func concurrentLoopWorker(apis chan concurrentLoopIn, out chan concurrentLoopOut) {
	for task := range apis {
		handleTasks(task.stack, task.task, 0)
		out <- concurrentLoopOut{index: task.index, result: ""}
	}
}

func triggerConcurrentLoop(stack *hub.Stack, task *hub.ScheduleApiDef, loopLength int, loop map[string]int, loopResult []interface{}) {
	var taskCount, msgCount int
	counter := loopLength

	if task.Control.ConcurrentLoopNum > loopLength {
		taskCount = loopLength
		msgCount = loopLength
	} else {
		taskCount = task.Control.ConcurrentLoopNum
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
		loop[task.Control.ResultKey] = i
		in <- concurrentLoopIn{index: i, stack: copyScheduleStack(stack, task), task: task.Control.Steps}
	}

	for i = 0; i < taskCount; i++ {
		go concurrentLoopWorker(in, out)
	}

	i = msgCount
	for result := range out {
		loopResult[result.index] = result.result
		counter--
		logger.LogS().Infoln(stack.BaseString, "loop并行处理结束：", counter, " result:", result)
		if i < loopLength {
			loop[task.Control.ResultKey] = i
			tmpStack := copyScheduleStack(stack, task)
			in <- concurrentLoopIn{index: i, stack: tmpStack, task: task.Control.Steps}
			i++
		} else if counter == 0 {
			break
		}
	}
}

func handleLoopTask(stack *hub.Stack, task *hub.ScheduleApiDef) (interface{}, int) {
	var result interface{}
	keyStr, _ := util.GetParameterStringValue(stack, nil, &task.Control.Key)

	if len(keyStr) == 0 {
		err := "invalid loop key"
		logger.LogS().Errorln(stack.BaseString, err)
		return util.CreateTmsError(hub.TmsErrorCoreId, err, nil), http.StatusInternalServerError
	}
	loopLength, _ := strconv.Atoi(keyStr)
	loopResult := make([]interface{}, loopLength)
	if isNormalMode(task) && len(task.Control.ResultKey) > 0 {
		stack.Heap[task.Control.ResultKey] = loopResult
	}

	loop := stack.Heap[hub.HeapLoopName].(map[string]int)
	if task.Control.ConcurrentLoopNum > 1 && loopLength > 1 {
		triggerConcurrentLoop(stack, task, loopLength, loop, loopResult)
	} else {
		for i := 0; i < loopLength; i++ {
			loop[task.Control.Name] = i
			result, _ = handleTasks(stack, task.Control.Steps, task.Control.ConcurrentNum)
			loopResult[i] = result
		}
	}
	return loopResult, 200
}

func handleApiTask(stack *hub.Stack, task *hub.ScheduleApiDef) (result interface{}, status int) {
	// 执行API
	result, status = ApiRun(stack, task.Api, task.Private, false)

	if isNormalMode(task) && len(task.Api.ResultKey) > 0 {
		stack.Heap[task.Api.ResultKey] = result
	}
	return
}

func handleOneScheduleApi(stack *hub.Stack, task *hub.ScheduleApiDef) (result interface{}, status int) {
	if len(task.Type) > 0 {
		switch task.Type {
		case "switch":
			logger.LogS().Infoln(stack.BaseString, "运行 switch name：", task.Control.Name)
			if task.Control.Cases != nil {
				return handleSwitchTask(stack, task)
			} else {
				err := "No switch cases"
				logger.LogS().Errorln(stack.BaseString, err)
				return util.CreateTmsError(hub.TmsErrorCoreId, err, nil), http.StatusInternalServerError
			}
		case hub.HeapLoopName:
			logger.LogS().Infoln(stack.BaseString, "运行 loop name", task.Control.Name)
			return handleLoopTask(stack, task)
		case "api":
			logger.LogS().Infoln(stack.BaseString, "运行 api name", task.Api.Name)
			result, status = handleApiTask(stack, task)
		default:
			err := "don't support type " + task.Type
			logger.LogS().Errorln(stack.BaseString, err)
			return util.CreateTmsError(hub.TmsErrorCoreId, err, nil), http.StatusInternalServerError
		}
	}
	return result, status
}

func concurrentScheWorker(stack *hub.Stack, apis chan concurrentScheIn, out chan concurrentScheOut) {
	for task := range apis {
		logger.LogS().Infoln(stack.BaseString, "并行运行 type：", task.task.Type)
		result, _ := handleOneScheduleApi(stack, task.task)
		out <- concurrentScheOut{task: task.task, result: result}
	}
}

func waitConcurrentScheResult(stack *hub.Stack, out chan concurrentScheOut, counter int) interface{} {
	results := make(map[string]interface{}, counter)
	var lastKey string
	for counter > 0 {
		//等待结果
		result := <-out
		var key string
		switch result.task.Type {
		case "api":
			key = result.task.Api.ResultKey
		default:
			key = result.task.Control.ResultKey
		}
		logger.LogS().Infoln(stack.BaseString, "并行处理结束：", counter, " result:", result)
		if len(key) > 0 {
			results[key] = result.result
			lastKey = key
		}
		counter--
	}
	//防止并发读写crash
	for k, v := range results {
		stack.Heap[k] = v
	}

	//由于并行，最后的结果并不确定，所以并行的返回结果不是固定的，因此当需要返回值时，最后一个应该是非并行的
	return results[lastKey]
}

func handleTasks(stack *hub.Stack, apis *[]hub.ScheduleApiDef, concurrentNum int) (result interface{}, status int) {
	var counter int
	var in chan concurrentScheIn
	var out chan concurrentScheOut
	if concurrentNum > 1 {
		/*假设所有的task都是并行的，多留buffer，提升性能*/
		in = make(chan concurrentScheIn, len(*apis))
		defer close(in)
		out = make(chan concurrentScheOut, len(*apis))
		defer close(out)
		for i := 0; i < concurrentNum; i++ {
			go concurrentScheWorker(stack, in, out)
		}
	}
	if apis == nil {
		str := "apis nil"
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorCoreId, str, nil), http.StatusInternalServerError
	}
	logger.LogS().Infoln(stack.BaseString, "apis lens：", len(*apis))
	for index := range *apis {
		task := &(*apis)[index]
		if concurrentNum > 1 {
			if task.Mode == "concurrent" {
				logger.LogS().Infoln(stack.BaseString, "准备并行运行 type：", task.Type, ",concurrentNum:", concurrentNum)
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
		if task.Mode == "background" {
			logger.LogS().Infoln(stack.BaseString, "后台 type：", task.Type)
			go handleOneScheduleApi(stack, task)
		} else {
			logger.LogS().Infoln(stack.BaseString, "串行 type：", task.Type, ", concurrentNum:", concurrentNum)
			result, status = handleOneScheduleApi(stack, task)
		}
	}

	//防止都是并行任务
	if counter > 0 {
		result = waitConcurrentScheResult(stack, out, counter)
	}
	return result, status
}

func runSchedule(stack *hub.Stack, name string, private string) (interface{}, int) {
	scheduleDef, ok := util.FindScheduleDef(name)
	if !ok || scheduleDef == nil || scheduleDef.Steps == nil {
		str := "获得Schedule定义失败：" + name
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorCoreId, str, nil), http.StatusInternalServerError
	}
	stack.Heap[hub.HeapLoopName] = make(map[string]int)

	return handleTasks(stack, scheduleDef.Steps, scheduleDef.ConcurrentNum)
}

func runScheduleApi(stack *hub.Stack, params map[string]string) (interface{}, int) {
	name, OK := params["name"]
	if !OK {
		str := "缺少flow名称"
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorCoreId, str, nil), http.StatusInternalServerError
	}
	private := params["private"]

	return runSchedule(stack, name, private)
}

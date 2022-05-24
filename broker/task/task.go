package task

import (
	"sync"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/unit"
	klog "k8s.io/klog/v2"
)

var mapLock sync.Mutex
var taskMap = make(map[string]hub.TaskHandler)

func RegisterTasks(list map[string]hub.TaskHandler) {
	mapLock.Lock()
	defer mapLock.Unlock()

	for k, v := range list {
		_, OK := taskMap[k]
		if OK {
			str := "task重名" + k
			klog.Errorln(str)
			panic(str)
		} else {
			taskMap[k] = v
		}
	}
}

// task调用
func Run(stack *hub.Stack, task *hub.TaskDef) (jsonOutRspBody interface{}, ret int) {
	function := taskMap[task.Command]
	if function == nil {
		str := "不能执行" + stack.ChildName
		klog.Errorln(str)
		panic(str)
	}

	var parameters map[string]string
	var origin map[string]interface{}

	if task.Parameters != nil {
		parameters = make(map[string]string)
		for index := range *task.Parameters {
			item := (*task.Parameters)[index]
			parameters[item.Name] = unit.GetParameterValue(stack, nil, item.From)
		}
	}

	if task.OriginParameters != nil {
		origin = stack.StepResult[hub.OriginName].(map[string]interface{})
		for index := range *task.OriginParameters {
			item := (*task.OriginParameters)[index]
			origin[item.Name] = unit.GetParameterValue(stack, nil, item.From)
		}
	}
	return function(stack, parameters)
}

package core

import (
	"net/http"
	"sync"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
	klog "k8s.io/klog/v2"
)

var mapLock sync.Mutex
var apiMap = make(map[string]hub.ApiHandler)

func RegisterApis(list map[string]hub.ApiHandler) {
	mapLock.Lock()
	defer mapLock.Unlock()

	for k, v := range list {
		_, OK := apiMap[k]
		if OK {
			str := "task重名:" + k
			klog.Errorln(str)
			panic(str)
		} else {
			apiMap[k] = v
		}
	}
}

// task调用
func ApiRun(stack *hub.Stack, api *hub.ApiDef, private string) (result interface{}, ret int) {
	function := apiMap[api.Command]
	var err error
	if function == nil {
		str := "不能执行" + stack.ChildName
		klog.Errorln(str)
		panic(str)
	}

	var origin map[string]interface{}
	parameters := make(map[string]string)
	var privateDef *hub.PrivateArray
	if len(api.Private) > 0 {
		parameters["private"] = api.Private
		privateDef, err = util.FindPrivateDef(stack, api.Private)
		if err != nil {
			klog.Errorln("获得private定义失败：", err)
			panic(err)
		}
	}

	if api.Parameters != nil {
		for index := range *api.Parameters {
			item := (*api.Parameters)[index]
			parameters[item.Name], err = util.GetParameterStringValue(stack, privateDef, &item.Value)
			if err != nil {
				return nil, http.StatusInternalServerError
			}
		}
	}

	if api.OriginParameters != nil {
		origin = stack.StepResult[hub.OriginName].(map[string]interface{})
		for index := range *api.OriginParameters {
			item := (*api.OriginParameters)[index]
			origin[item.Name], err = util.GetParameterRawValue(stack, privateDef, &item.Value)
			if err != nil {
				return nil, http.StatusInternalServerError
			}
		}
	}
	return function(stack, parameters)
}

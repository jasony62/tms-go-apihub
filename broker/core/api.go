package core

import (
	"net/http"
	"sync"
	"time"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
	klog "k8s.io/klog/v2"
)

var mapLock sync.Mutex
var apiMap = make(map[string]hub.ApiHandler)

func RegisterApis(list map[string]hub.ApiHandler) {
	mapLock.Lock()
	defer mapLock.Unlock()

	for k := range list {
		_, OK := apiMap[k]
		if OK {
			str := "main : task重名:" + k
			klog.Errorln(str)
			panic(str)
		} else {
			apiMap[k] = list[k]
		}
	}
}

func preApis(stack *hub.Stack, apiDef *hub.ApiDef) {
	klog.Infoln("___pre API,", stack.BaseString, "command:", apiDef.Command, "name:"+apiDef.Name)
}

func postApis(stack *hub.Stack, apiDef *hub.ApiDef, result interface{}, code int, duration float64) {
	if stack == nil {
		return
	}

	if code == http.StatusOK {
		klog.Infoln("___post API OK:", stack.BaseString, "command:"+apiDef.Command, " name："+apiDef.Name, " result:", result, " duration:", duration)
	} else {
		klog.Errorln("!!!post API NOK:", stack.BaseString, "command :"+apiDef.Command, " name："+apiDef.Name, " result:", result, " duration:", duration)
	}
}

// task调用
func ApiRun(stack *hub.Stack, api *hub.ApiDef, private string, internal bool) (result interface{}, ret int) {
	var t time.Time
	if !internal {
		t = time.Now()
	}
	function := apiMap[api.Command]
	var err error
	if function == nil {
		str := "不能执行" + api.Command
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorCoreId, str, nil), http.StatusForbidden
	}

	var origin map[string]interface{}
	args := make(map[string]string)
	var privateDef *hub.PrivateArray
	if len(api.Private) > 0 {
		var ok bool
		args["private"] = api.Private
		privateDef, ok = util.FindPrivateDef(api.Private)
		if !ok || err != nil {
			klog.Errorln(stack.BaseString, "获得private定义失败：", api.Private)
			return nil, http.StatusForbidden
		}
	}

	if api.Args != nil {
		for index := range *api.Args {
			item := (*api.Args)[index]
			args[item.Name], err = util.GetParameterStringValue(stack, privateDef, &item.Value)
			if err != nil {
				klog.Errorln(stack.BaseString, "获得value失败：", err)
				return nil, http.StatusInternalServerError
			}
		}
	}

	if api.OriginParameters != nil {
		origin = stack.Heap[hub.HeapOriginName].(map[string]interface{})
		for index := range *api.OriginParameters {
			item := (*api.OriginParameters)[index]
			origin[item.Name], err = util.GetParameterRawValue(stack, privateDef, &item.Value)
			if err != nil {
				klog.Errorln(stack.BaseString, "获得origin失败：", err)
				return nil, http.StatusInternalServerError
			}
		}
	}
	if !internal {
		preApis(stack, api)
	}
	result, ret = function(stack, args)
	if !internal {
		duration := time.Since(t).Seconds()
		postApis(stack, api, result, ret, duration)
	}
	return
}

package core

import (
	"net/http"
	"sync"
	"time"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
	"go.uber.org/zap"
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
			zap.S().Errorln(str)
			panic(str)
		} else {
			apiMap[k] = list[k]
		}
	}
}

func preApis(stack *hub.Stack, apiDef *hub.ApiDef) {
	//	zap.S().Infoln("___pre API,", stack.BaseString, "command:", apiDef.Command, "name:"+apiDef.Name)
}

func postApis(stack *hub.Stack, apiDef *hub.ApiDef, result interface{}, code int, duration float64) {
	if stack == nil {
		return
	}

	if code == http.StatusOK {
		zap.S().Infoln("___post API OK: ", stack.BaseString, "command:"+apiDef.Command, " name："+apiDef.Name, " result:", result, " duration(s):", duration)
	} else {
		zap.S().Errorln("!!!post API NOK:", stack.BaseString, "command :"+apiDef.Command, " name："+apiDef.Name, " result:", result, " duration(s):", duration)
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
		zap.S().Errorln(stack.BaseString, str)
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
			str := "获得private定义失败：" + api.Private
			zap.S().Errorln(stack.BaseString, str)
			return util.CreateTmsError(hub.TmsErrorCoreId, str, nil), http.StatusForbidden
		}
	}

	if api.Args != nil {
		for index := range *api.Args {
			item := (*api.Args)[index]
			args[item.Name], err = util.GetParameterStringValue(stack, privateDef, &item.Value)
			if err != nil {
				str := "获得value失败：" + err.Error()
				zap.S().Errorln(stack.BaseString, str)
				return util.CreateTmsError(hub.TmsErrorCoreId, str, nil), http.StatusInternalServerError
			}
		}
	}

	if api.OriginParameters != nil {
		origin = stack.Heap[hub.HeapOriginName].(map[string]interface{})
		for index := range *api.OriginParameters {
			item := (*api.OriginParameters)[index]
			origin[item.Name], err = util.GetParameterRawValue(stack, privateDef, &item.Value)
			if err != nil {
				str := "获得origin失败：" + err.Error()
				zap.S().Errorln(stack.BaseString, str)
				return util.CreateTmsError(hub.TmsErrorCoreId, str, nil), http.StatusInternalServerError
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

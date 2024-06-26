package apis

import (
	"net/http"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/logger"
	"github.com/jasony62/tms-go-apihub/util"
	"github.com/valyala/fasthttp"
)

type storageMap struct {
	StorageMap map[string]string
}

var storeMap = storageMap{
	StorageMap: make(map[string]string),
}

func storageStore(stack *hub.Stack, params map[string]string) (interface{}, int) {
	/*
		引入storage层，引入两个API storageStore，storageLoad
		   	*parameters：
		   **source： local（内存kv） 后面再支持mongodb（支持数据库）
		   **index：索引
		   **content
		   store时是json，则从origin中创建string，否则按string直接存储
		   load时是json需要从获取的string解析为json，其他忽略

		   *origin：store时候为json时要存入的信息
	*/

	var user string
	var key string
	var index string
	var source string
	var content string
	var OK bool

	// user, OK = params["user"]
	// if !OK {
	// 	str := "storageStore缺少user定义"
	// 	logger.LogS().Errorln(stack.BaseString, str)
	// 	return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	// }

	// logger.LogS().Infoln("storageStore user: ", user)
	// if len(user) == 0 {
	// 	str := "storageStore缺少user"
	// 	logger.LogS().Errorln(stack.BaseString, str)
	// 	return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	// }

	key, OK = params["key"]
	if !OK {
		str := "storageStore缺少key索引"
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}

	index, OK = params["index"]
	if !OK {
		str := "storageStore缺少index索引"
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}

	source, OK = params["source"]
	if !OK {
		source = "local"
	}

	content, OK = params["content"]
	if !OK {
		str := "storageStore缺少存储内容content"
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}

	logger.LogS().Infoln("storageStore: user:", user, "params:", params)
	if source == "local" {
		return storeLocal(stack, user, key, index, content)
	}
	return nil, fasthttp.StatusInternalServerError
}

func storageLoad(stack *hub.Stack, params map[string]string) (interface{}, int) {
	var index string
	var source string
	var content string
	var OK bool

	index, OK = params["index"]
	if !OK {
		str := "storageLoad缺少index索引"
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}

	source, OK = params["source"]
	if !OK {
		source = "local"
	}

	content, OK = params["content"]
	if !OK {
		str := "storageLoad缺少存储内容content"
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}

	if source == "local" {
		return loadLocal(stack, index, source, content)
	}

	return nil, fasthttp.StatusInternalServerError
}

func storeLocal(stack *hub.Stack, user string, key string, index string, content string) (interface{}, int) {
	if _, ok := storeMap.StorageMap[index]; ok {
		logger.LogS().Infoln("storageStore索引为:", key, "的值已经存在，覆盖之")
	}

	logger.LogS().Infoln("storeLocal: index:", index, " user:", user, " content:", content)
	tmp := stack.Heap[hub.HeapOriginName].(map[string]interface{})
	result := tmp[key]
	byteJson, err := jsonEx.Marshal(result)
	if err != nil {
		str := "storeLocal解析失败" + err.Error()
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	if content == "json" {
		storeMap.StorageMap[index] = string(byteJson)
	} else {
		storeMap.StorageMap[index] = content
	}

	logger.LogS().Infoln("storeLocal:", storeMap.StorageMap[index])
	return nil, fasthttp.StatusOK
}

func loadLocal(stack *hub.Stack, index string, source string, content string) (interface{}, int) {
	var val string
	var ok bool
	logger.LogS().Infoln("loadLocal:", index, " source:", source, " content:", content)

	if val, ok = storeMap.StorageMap[index]; !ok {
		str := "loadLocal加载" + index + "失败"
		logger.LogS().Warnln(stack.BaseString, str)
	}

	var ret interface{}
	if content == "json" {
		if val == "" {
			str := "加载" + index + "失败"
			logger.LogS().Errorln(stack.BaseString, str)
			return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
		}
		jsonEx.Unmarshal([]byte(val), &ret)
	} else {
		if val == "" {
			ret = content
		} else {
			ret = val
		}
	}
	logger.LogS().Infoln("loadLocal value:", val)

	return ret, fasthttp.StatusOK
}

func storageClear(stack *hub.Stack, params map[string]string) (interface{}, int) {
	var index string
	var source string
	var OK bool

	index, OK = params["index"]
	if !OK {
		str := "storageClear缺少index索引"
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}
	source, OK = params["source"]
	if !OK {
		source = "local"
	}

	if source == "local" {
		delete(storeMap.StorageMap, index)
	}
	return nil, fasthttp.StatusOK
}

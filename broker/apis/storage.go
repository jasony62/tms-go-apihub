package apis

import (
	"net/http"

	"github.com/Sheng-ZM/tms-go-apihub/broker//hub"
	"github.com/Sheng-ZM/tms-go-apihub/broker//util"
	"github.com/valyala/fasthttp"
	klog "k8s.io/klog/v2"
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

	user, OK = params["user"]
	if !OK {
		str := "storageStore缺少user定义"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}

	klog.Infoln("storageStore user: ", user)
	if len(user) == 0 {
		str := "storageStore缺少user"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}

	key, OK = params["key"]
	if !OK {
		str := "storageStore缺少key索引"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}

	index, OK = params["index"]
	if !OK {
		str := "storageStore缺少index索引"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}

	source, OK = params["source"]
	if !OK {
		source = "local"
	}

	content, OK = params["content"]
	if !OK {
		str := "storageStore缺少存储内容content"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}

	klog.Infoln("storageStore: user:", user, "params:", params)
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
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}

	source, OK = params["source"]
	if !OK {
		source = "local"
	}

	content, OK = params["content"]
	if !OK {
		str := "storageLoad缺少存储内容content"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}

	if source == "local" {
		return loadLocal(stack, index, source, content)
	}

	return nil, fasthttp.StatusInternalServerError
}

func storeLocal(stack *hub.Stack, user string, key string, index string, content string) (interface{}, int) {
	if _, ok := storeMap.StorageMap[index]; ok {
		klog.Infoln("storageStore索引为:", key, "的值已经存在，覆盖之")
	}

	klog.Infoln("storeLocal: index:", index, " user:", user, " content:", content)
	tmp := stack.Heap[hub.HeapOriginName].(map[string]interface{})
	result := tmp[key]
	klog.Infoln("storeLocal: result:", result)
	byteJson, err := jsonEx.Marshal(result)
	if err != nil {
		str := "storeLocal解析失败" + err.Error()
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	if content == "json" {
		storeMap.StorageMap[index] = string(byteJson)
	} else {
		storeMap.StorageMap[index] = content
	}

	klog.Infoln("storeLocal:", storeMap.StorageMap[index])
	return nil, fasthttp.StatusOK
}

func loadLocal(stack *hub.Stack, index string, source string, content string) (interface{}, int) {
	var val string
	var ok bool
	klog.Infoln("loadLocal:", index, " source:", source, " content:", content)

	if val, ok = storeMap.StorageMap[index]; !ok {
		str := "loadLocal加载失败" + index
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	klog.Infoln("loadLocal value:", val)
	var ret interface{}
	if content == "json" {
		jsonEx.Unmarshal([]byte(val), &ret)
	} else {
		ret = val
	}

	//	klog.Infoln("loadLocal ret:", ret)
	delete(storeMap.StorageMap, index)
	return ret, fasthttp.StatusOK
}

package apis

import (
	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
	"github.com/valyala/fasthttp"
)

func checkRight(stack *hub.Stack, params map[string]string) (interface{}, int) {
	var userKey string
	var name string
	var apiType string
	var OK bool

	userKey, OK = params["userKey"]
	if !OK {
		str := "缺少userKey定义"
		klog.Errorln(str)
		panic(str)
	}

	name, OK = params["name"]
	if !OK {
		str := "缺少api名称"
		klog.Errorln(str)
		panic(str)
	}

	apiType, OK = params["type"]
	if !OK {
		str := "缺少type类型"
		klog.Errorln(str)
		panic(str)
	}

	klog.Infoln("checkRight name: ", name, " type: ", apiType)

	//判断执行权限
	if !util.CheckRight(stack, userKey, name, apiType) {
		return nil, fasthttp.StatusInternalServerError
	}

	return nil, fasthttp.StatusOK
}

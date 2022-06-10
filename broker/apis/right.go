package apis

import (
	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
	"github.com/valyala/fasthttp"
)

func checkRight(stack *hub.Stack, params map[string]string) (interface{}, int) {
	var user string
	var name string
	var apiType string
	var OK bool

	user, OK = params["user"]
	if !OK {
		klog.Infoln("缺少user定义，不检查权限")
		return nil, fasthttp.StatusOK
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

	//判断执行权限
	if !util.CheckRight(stack, user, name, apiType) {
		return nil, fasthttp.StatusInternalServerError
	}

	return nil, fasthttp.StatusOK
}

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
		str := "缺少user定义"
		klog.Errorln(str)
		panic(str)
	}

	klog.Infoln("checkRight user: ", user)
	if len(user) == 0 {
		klog.Errorln("checkRight user is null")
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

	klog.Infoln("checkRight name: ", name, " type: ", apiType)

	//判断执行权限
	if !util.CheckRight(stack, user, name, apiType) {
		return nil, fasthttp.StatusInternalServerError
	}

	return nil, fasthttp.StatusOK
}

package apis

import (
	"net/http"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/broker/hub"
	"github.com/jasony62/tms-go-apihub/broker/util"
	"github.com/valyala/fasthttp"
)

var allowDefaultAccess = true

func userInList(arr *hub.RightArray, user string) bool {
	if arr.List != nil {
		for _, u := range *arr.List {
			if user == u.User {
				return true
			}
		}
	}
	return false
}

func hasRight(stack *hub.Stack, user string, name string, callType string) (interface{}, int) {
	// check是否有权限
	//	klog.Infoln("CheckRight user:", user, " callType:", callType, " name:", name)
	rightInfo := util.FindRightDef(user, name, callType)

	haveRight := false
	if rightInfo != nil {
		switch rightInfo.Right {
		case hub.Right_Pulbic:
			haveRight = true
		case hub.Right_Internal:
			haveRight = false
		case hub.Right_Whitelist:
			haveRight = userInList(rightInfo, user)
		case hub.Right_Blacklist:
			haveRight = !userInList(rightInfo, user)
		default:
			klog.Infoln("CheckRight invalid right: ", rightInfo.Right)
			haveRight = false
		}
	} else if allowDefaultAccess {
		haveRight = true
	}

	if !haveRight {
		str := "Deny access right for: " + user + ",api " + name
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	} else {
		return nil, http.StatusOK
	}
}

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
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}

	apiType, OK = params["type"]
	if !OK {
		str := "缺少type类型"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusForbidden
	}

	//判断执行权限
	return hasRight(stack, user, name, apiType)
}

func setDefaultAccessRight(stack *hub.Stack, params map[string]string) (interface{}, int) {
	policy, OK := params["default"]
	if !OK {
		str := "缺少default权限值"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), 400
	}
	switch policy {
	case "deny":
		allowDefaultAccess = false
		klog.Infoln("default access policy: deny")
	default:
		klog.Infoln("default access policy: access")
	}
	return nil, http.StatusOK
}

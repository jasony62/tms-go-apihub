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
		//	user = stack.GinContext.Query("appID")
		// 作为api执行时，params["user"]不会定义，需要去apidef 中查询是否定义了query或者header中有user
		var err error
		var value string
		HttpApi, err := util.FindHttpApiDef(stack.RootName)
		if err != nil {
			klog.Errorln("获得private定义失败：", err)
			panic(err)
		}

		if HttpApi.Args != nil {
			paramLen := len(*HttpApi.Args)
			if paramLen > 0 {
				for _, param := range *HttpApi.Args {
					if len(param.Name) > 0 {
						value, err = util.GetParameterStringValue(stack, nil, &param.Value)
						if err != nil {
							return nil, fasthttp.StatusOK
						}

						if param.Name == "user" {
							if len(value) > 0 {
								user = value
								klog.Infoln("查找到user定义：", user)
							}
							break
						}
					}
				}
			}
		}

		if len(user) == 0 {
			klog.Infoln("缺少user定义，不检查权限")
			return nil, fasthttp.StatusOK
		}
	}

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

	//判断执行权限
	if !util.CheckRight(stack, user, name, apiType) {
		return nil, fasthttp.StatusInternalServerError
	}

	return nil, fasthttp.StatusOK
}

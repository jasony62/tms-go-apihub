package apis

import (
	"net/http"

	"github.com/jasony62/tms-go-apihub/broker/hub"
	"github.com/jasony62/tms-go-apihub/broker/util"
	klog "k8s.io/klog/v2"
)

func createHtml(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		str := "createHtml,缺少参数"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	name, OK := params["type"]
	if !OK {
		str := "createHtml,type为空"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	content, OK := params["content"]
	if !OK {
		str := "createHtml,content为空"
		klog.Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	if name == "resource" {
		content, OK = util.FindResourceDef(content)
		if !OK {
			str := "createHtml FindResourceDef failed"
			klog.Errorln(stack.BaseString, str)
			return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
		}
	}

	result, err := util.Json2Html(stack.Heap, content)
	if err != nil {
		return util.CreateTmsError(hub.TmsErrorApisId, err.Error(), nil), http.StatusInternalServerError
	}
	return result, 200
}

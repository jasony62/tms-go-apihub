package apis

import (
	"net/http"

	"github.com/jasony62/tms-go-apihub/broker/hub"
	"github.com/jasony62/tms-go-apihub/broker/util"
	klog "k8s.io/klog/v2"
)

func dump(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		str := "dump参数为空"
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}
	klog.Infoln("\r\n****************DUMP:\r\n", stack.BaseString, " params:", params, "\r\n")

	return nil, http.StatusOK
}

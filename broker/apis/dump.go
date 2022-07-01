package apis

import (
	"net/http"

	"github.com/jasony62/tms-go-apihub/hub"
	klog "k8s.io/klog/v2"
)

func dump(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		return nil, http.StatusInternalServerError
	}
	klog.Infoln("\r\n****************DUMP:\r\n", params, "\r\n")

	return nil, 200
}

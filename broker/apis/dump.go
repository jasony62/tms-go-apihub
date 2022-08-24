package apis

import (
	"net/http"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/logger"
	"github.com/jasony62/tms-go-apihub/util"
)

func dump(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		str := "dump参数为空"
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}
	logger.LogS().Infoln("\r\n****************DUMP:\r\n", stack.BaseString, " params:", params, "\r\n")

	return nil, http.StatusOK
}

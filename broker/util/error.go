package util

import (
	"encoding/base64"
	"runtime"

	"github.com/jasony62/tms-go-apihub/broker/hub"
)

func CreateTmsError(id uint, msg string, err error) (ret hub.TmsError) {
	//获取的是 Caller这个函数的调用栈
	_, ret.Module, ret.Line, _ = runtime.Caller(1)
	ret.Module = base64.StdEncoding.EncodeToString([]byte(ret.Module))
	ret.Id = id
	ret.ErrorMsg = msg
	ret.Error = err
	/*TODO add to heap?*/
	return
}

package apis

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/logger"
	"github.com/jasony62/tms-go-apihub/util"
)

var LogConf logger.LogConfigs
var LogWithLevel bool

func logOutput(stack *hub.Stack, params map[string]string) (interface{}, int) {
	// 1. 配置log参数
	LogConf.LogPath = getLogConf(params["filepath"], "../log/")
	LogConf.LogFileName = getLogConf(params["filename"], "")
	LogConf.LogFormat = getLogConf(params["logformat"], "logfmt")
	LogConf.LogLevel = getLogConf(params["loglevel"], "info")
	LogConf.LogFileMaxSize, _ = strconv.Atoi(getLogConf(params["fileMaxSize"], "50"))
	LogConf.LogFileMaxBackups, _ = strconv.Atoi(getLogConf(params["fileMaxBackups"], "100"))
	LogConf.LogMaxAge, _ = strconv.Atoi(getLogConf(params["maxDays"], "10"))
	LogConf.LogCompress, _ = strconv.ParseBool(getLogConf(params["compress"], "false"))
	LogConf.LogStdout, _ = strconv.ParseBool(getLogConf(params["stdout"], "true"))
	LogWithLevel, _ = strconv.ParseBool(getLogConf(params["logwithlevel"], "false"))

	// 2. 初始化log
	if err := logger.InitLogger(LogConf, LogWithLevel); err != nil {
		str := "初始化日志系统失败" + err.Error()
		fmt.Println(str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func getLogConf(inputStr string, defaultStr string) string {
	if len(inputStr) == 0 {
		return defaultStr
	} else {
		return inputStr
	}
}

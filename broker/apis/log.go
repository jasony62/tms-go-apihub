package apis

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/logger"
	"github.com/jasony62/tms-go-apihub/util"
)

var logConfig logger.LogConfigs
var logWithLevel bool

func logOutput(stack *hub.Stack, params map[string]string) (interface{}, int) {
	// 1. 配置log参数
	logConfig.LogPath = getLogConf(params["filepath"], "../log/")
	logConfig.LogFileName = getLogConf(params["filename"], "")
	logConfig.LogFormat = getLogConf(params["logformat"], "logfmt")
	logConfig.LogLevel = getLogConf(params["loglevel"], "info")
	logConfig.LogFileMaxSize, _ = strconv.Atoi(getLogConf(params["fileMaxSize"], "50"))
	logConfig.LogFileMaxBackups, _ = strconv.Atoi(getLogConf(params["fileMaxBackups"], "100"))
	logConfig.LogMaxAge, _ = strconv.Atoi(getLogConf(params["maxDays"], "10"))
	logConfig.LogCompress, _ = strconv.ParseBool(getLogConf(params["compress"], "false"))
	logConfig.LogStdout, _ = strconv.ParseBool(getLogConf(params["stdout"], "true"))
	logWithLevel, _ = strconv.ParseBool(getLogConf(params["logwithlevel"], "false"))

	// 2. 初始化log
	if err := logger.InitLogger(logConfig, logWithLevel); err != nil {
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

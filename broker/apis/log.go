package apis

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/jasony62/tms-go-apihub/hub"

	klog "k8s.io/klog/v2"
)

func logToFile(stack *hub.Stack, params map[string]string) (interface{}, int) {
	filename := params["filename"]

	klog.Infoln("logToFile ", filename)
	if len(filename) == 0 {
		return nil, http.StatusOK
	}

	folder := "../log/"
	exist, _ := pathExists(folder)
	if !exist {
		os.Mkdir(folder, os.ModePerm)
	}

	//获取当前时间，命名log文件
	timeStr := time.Now().Format("20060102150405")
	flag.Set("logtostderr", "false") // By default klog logs to stderr, switch that off

	level := params["loglevel"]
	if len(level) > 0 {
		//klog写日志到文件不支持分级别输出，打印到屏幕支持，如果按级别打印到屏幕需要关闭alsologtostderr
		flag.Set("stderrthreshold", level) //InfoLog:"INFO",WarningLog:"WARNING",ErrorLog:"ERROR",FatalLog:"FATAL",
	} else {
		flag.Set("alsologtostderr", "true") // false is default, but this is informative
	}

	logName := folder + filename + "_" + timeStr + ".log"
	flag.Set("log_file", logName)
	flag.Parse()
	return nil, http.StatusOK
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

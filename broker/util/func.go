package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"plugin"
	"strconv"
	"strings"
	"time"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/logger"
)

func utcFunc(params []string) string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func utcmsFunc(params []string) string {
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
}

func utcTemplate(args ...interface{}) string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
func utcmsTemplate(args ...interface{}) string {
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
}

func md5Template(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	str := fmt.Sprint(args...)
	w := md5.New()
	io.WriteString(w, str)
	checksum := fmt.Sprintf("%x", w.Sum(nil))
	return checksum
}

func md5Func(params []string) string {
	if len(params) == 0 {
		return ""
	}
	var str string
	for _, v := range params {
		str = str + v
	}
	w := md5.New()
	io.WriteString(w, str)
	checksum := fmt.Sprintf("%x", w.Sum(nil))
	return checksum
}

func timestampFunc(params []string) string {
	//	return time.Now().Format("2006-01-02 15:04:05")
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	return time.Now().In(cstZone).Format("2006-01-02 15:04:05")
}

func timestampTemplate(args ...interface{}) string {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	return time.Now().In(cstZone).Format("2006-01-02 15:04:05")
}

var funcMap map[string]hub.FuncHandler = map[string]hub.FuncHandler{
	"utc":       utcFunc,
	"utc_ms":    utcmsFunc,
	"md5":       md5Func,
	"timestamp": timestampFunc,
}

var funcMapForTemplate map[string](interface{}) = map[string](interface{}){
	"utc":       utcTemplate,
	"utc_ms":    utcmsTemplate,
	"md5":       md5Template,
	"timestamp": timestampTemplate,
}

func loadConfigPluginData(path string) {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		logger.LogS().Errorln(err.Error())
		return
	}
	var prefix string
	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())

		if fileInfoList[i].IsDir() {
			logger.LogS().Infoln("Plugin子目录: ", fileName)
			prefix = fileInfoList[i].Name()
			loadConfigPluginData(path + "/" + prefix)
		} else {
			if !strings.HasSuffix(fileName, ".so") {
				continue
			}
			logger.LogS().Infoln("加载Plugin(*.so)文件: ", fileName)
			p, err := plugin.Open(fileName)
			if err != nil {
				logger.LogS().Errorln(err.Error())
				panic(err)
			}

			// 导入动态库，注册函数到funcMap
			registerFunc, err := p.Lookup("Register")
			if err != nil {
				logger.LogS().Errorln(err.Error())
				panic(err)
			}
			mapFunc, mapFuncForTemplate := registerFunc.(func() (map[string]interface{}, map[string]interface{}))()
			loadPluginFuncs(mapFunc, mapFuncForTemplate)
		}
	}
	logger.LogS().Infoln("加载func so文件完成！\r\n")
}

func loadPluginFuncs(mapFunc map[string]interface{}, mapFuncForTemplate map[string]interface{}) {
	for k, v := range mapFunc {
		if _, ok := funcMap[k]; ok {
			logger.LogS().Errorln("加载(%s)失败,FuncMap存在重名函数！\r\n", k)
		} else {
			funcMap[k] = v.(hub.FuncHandler)
		}
	}

	for k, v := range mapFuncForTemplate {
		if _, ok := funcMapForTemplate[k]; ok {
			logger.LogS().Errorln("加载(%s)失败,FuncMapForTemplate存在重名函数！\r\n", k)
		} else {
			funcMapForTemplate[k] = v
		}
	}
}

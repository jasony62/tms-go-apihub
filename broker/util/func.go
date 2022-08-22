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
	"go.uber.org/zap"
)

func utcFunc(params []string) string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func utcTemplate(args ...interface{}) string {
	return strconv.FormatInt(time.Now().Unix(), 10)
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

var funcMap map[string]hub.FuncHandler = map[string]hub.FuncHandler{
	"utc": utcFunc,
	"md5": md5Func,
}

var funcMapForTemplate map[string](interface{}) = map[string](interface{}){
	"utc": utcTemplate,
	"md5": md5Template,
}

func loadConfigPluginData(path string) {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		zap.S().Errorln(err.Error())
		return
	}
	var prefix string
	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())

		if fileInfoList[i].IsDir() {
			zap.S().Infoln("Plugin子目录: ", fileName)
			prefix = fileInfoList[i].Name()
			loadConfigPluginData(path + "/" + prefix)
		} else {
			if !strings.HasSuffix(fileName, ".so") {
				continue
			}
			zap.S().Infoln("加载Plugin(*.so)文件: ", fileName)
			p, err := plugin.Open(fileName)
			if err != nil {
				zap.S().Errorln(err.Error())
				panic(err)
			}

			// 导入动态库，注册函数到funcMap
			registerFunc, err := p.Lookup("Register")
			if err != nil {
				zap.S().Errorln(err.Error())
				panic(err)
			}
			mapFunc, mapFuncForTemplate := registerFunc.(func() (map[string]interface{}, map[string]interface{}))()
			loadPluginFuncs(mapFunc, mapFuncForTemplate)
		}
	}
	zap.S().Infoln("加载func so文件完成！\r\n")
}

func loadPluginFuncs(mapFunc map[string]interface{}, mapFuncForTemplate map[string]interface{}) {
	for k, v := range mapFunc {
		if _, ok := funcMap[k]; ok {
			zap.S().Errorln("加载(%s)失败,FuncMap存在重名函数！\r\n", k)
		} else {
			funcMap[k] = v.(hub.FuncHandler)
		}
	}

	for k, v := range mapFuncForTemplate {
		if _, ok := funcMapForTemplate[k]; ok {
			zap.S().Errorln("加载(%s)失败,FuncMapForTemplate存在重名函数！\r\n", k)
		} else {
			funcMapForTemplate[k] = v
		}
	}
}

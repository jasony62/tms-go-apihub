package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/logger"
)

// 应用的基本信息
type confMap struct {
	BasePath         string
	ApiMap           map[string]*hub.HttpApiDef
	PrivateMap       map[string]*hub.PrivateArray
	FlowMap          map[string]*hub.FlowDef
	ScheduleMap      map[string]*hub.ScheduleDef
	SourceMap        map[string]string
	ApiRightMap      map[string]*hub.RightArray
	FlowRightMap     map[string]*hub.RightArray
	ScheduleRightMap map[string]*hub.RightArray
}

var DefaultConfMap = confMap{
	BasePath:         "./conf/",
	ApiMap:           make(map[string]*hub.HttpApiDef),
	FlowMap:          make(map[string]*hub.FlowDef),
	ScheduleMap:      make(map[string]*hub.ScheduleDef),
	PrivateMap:       make(map[string]*hub.PrivateArray),
	SourceMap:        make(map[string]string),
	ApiRightMap:      make(map[string]*hub.RightArray),
	FlowRightMap:     make(map[string]*hub.RightArray),
	ScheduleRightMap: make(map[string]*hub.RightArray),
}

func loadConfigJsonData(paths []string) {
	logger.LogS().Infoln("加载API def文件...")
	for i := hub.JSON_TYPE_PRIVATE; i <= hub.JSON_TYPE_SCHEDULE; i++ {
		/*TODO add error return and panic if failure*/
		loadJsonDefData(i, paths[i], "", true)
	}

	for i := hub.JSON_TYPE_API_RIGHT; i <= hub.JSON_TYPE_SCHEDULE_RIGHT; i++ {
		/*TODO add error return and panic if failure*/
		loadJsonDefData(i, paths[i], "", true)
	}
}

func loadJsonDefData(jsonType int, path string, prefix string, includeDir bool) {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		logger.LogS().Errorln(err.Error())
		return
	}

	oldPrefix := prefix
	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())

		if fileInfoList[i].IsDir() && includeDir {
			prefix = fileInfoList[i].Name()
			loadJsonDefData(jsonType, path+"/"+prefix, prefix, true)
		} else {
			if !strings.HasSuffix(fileName, ".json") {
				continue
			}

			prefix = oldPrefix

			byteFile, err := ioutil.ReadFile(fileName)
			if err != nil {
				str := "获得Json定义失败：" + err.Error()
				logger.LogS().Errorln(str)
				panic(str)
			}

			if !json.Valid(byteFile) {
				str := "Json文件无效：" + fileName
				logger.LogS().Errorln(str)
				panic(str)
			}

			var key string
			fname := fileInfoList[i].Name()
			index := strings.Index(fname, ".json")
			if index >= 0 {
				fname = fname[:index]
			}

			key = fname

			decoder := json.NewDecoder(bytes.NewReader(byteFile))
			switch jsonType {
			case hub.JSON_TYPE_API:
				def := new(hub.HttpApiDef)
				decoder.Decode(&def)
				DefaultConfMap.ApiMap[key] = def
			case hub.JSON_TYPE_FLOW:
				def := new(hub.FlowDef)
				decoder.Decode(&def)
				DefaultConfMap.FlowMap[key] = def
			case hub.JSON_TYPE_SCHEDULE:
				def := new(hub.ScheduleDef)
				decoder.Decode(&def)
				DefaultConfMap.ScheduleMap[key] = def
			case hub.JSON_TYPE_PRIVATE:
				def := new(hub.PrivateArray)
				decoder.Decode(&def)
				DefaultConfMap.PrivateMap[key] = def
			case hub.JSON_TYPE_API_RIGHT:
				def := new(hub.RightArray)
				decoder.Decode(&def)
				DefaultConfMap.ApiRightMap[key] = def
			case hub.JSON_TYPE_FLOW_RIGHT:
				def := new(hub.RightArray)
				decoder.Decode(&def)
				DefaultConfMap.FlowRightMap[key] = def
			case hub.JSON_TYPE_SCHEDULE_RIGHT:
				def := new(hub.RightArray)
				decoder.Decode(&def)
				DefaultConfMap.ScheduleRightMap[key] = def
			default:
			}
		}
	}
}

func loadTemplateData(path string, prefix string) {
	logger.LogS().Infoln("加载Template文件...")
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		str := "invalid path " + path
		logger.LogS().Errorln(str)
		return
	}

	oldPrefix := prefix
	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())

		if fileInfoList[i].IsDir() {
			prefix = fileInfoList[i].Name()
			loadTemplateData(path+"/"+prefix, prefix)
		} else {
			prefix = oldPrefix

			byteFile, err := ioutil.ReadFile(fileName)
			if err != nil {
				str := "获得tmpl定义失败：" + err.Error()
				logger.LogS().Errorln(str)
				panic(str)
			}

			fname := fileInfoList[i].Name()
			DefaultConfMap.SourceMap[fname] = string(byteFile)
		}
	}
}

func FindHttpApiDef(name string) (value *hub.HttpApiDef, ok bool) {
	if len(name) == 0 {
		return nil, false
	}

	value, ok = DefaultConfMap.ApiMap[name]
	return
}

func FindPrivateDef(name string) (value *hub.PrivateArray, ok bool) {
	if len(name) == 0 {
		return nil, false
	}

	value, ok = DefaultConfMap.PrivateMap[name]
	return
}

func FindFlowDef(id string) (value *hub.FlowDef, ok bool) {
	value, ok = DefaultConfMap.FlowMap[id]
	return
}

func FindScheduleDef(id string) (value *hub.ScheduleDef, ok bool) {
	value, ok = DefaultConfMap.ScheduleMap[id]
	return
}

func FindResourceDef(id string) (value string, ok bool) {
	value, ok = DefaultConfMap.SourceMap[id]
	return
}

func FindRightDef(user string, name string, callType string) *hub.RightArray {
	// check是否有权限
	logger.LogS().Infoln("CheckRight user:", user, " callType:", callType, " name:", name)
	//map
	switch callType {
	case "httpapi":
		return DefaultConfMap.ApiRightMap[name]
	case "flow":
		return DefaultConfMap.FlowRightMap[name]
	case "schedule":
		return DefaultConfMap.ScheduleRightMap[name]
	default:
		return DefaultConfMap.ApiRightMap[name]

	}
}

func GetBasePath() string {
	return DefaultConfMap.BasePath
}

func LoadConf(basePath string) {
	loadConfigJsonData([]string{basePath + "privates",
		basePath + "httpapis", basePath + "flows",
		basePath + "schedules", basePath + "rights/httpapi",
		basePath + "rights/flow", basePath + "rights/schedule"})

	loadTemplateData(basePath+"templates", "")
	loadConfigPluginData(basePath + "plugins")
}

func LoadMainFlow(path string) (interface{}, int) {
	if len(path) > 0 {
		if path[len(path)-1] != '/' && path[len(path)-1] != '\\' {
			path = path + "/"
		}
		DefaultConfMap.BasePath = path
	}
	logger.LogS().Infoln("Load main flow from %s\n", DefaultConfMap.BasePath)
	loadJsonDefData(hub.JSON_TYPE_FLOW, DefaultConfMap.BasePath, "", false)
	return nil, http.StatusOK
}

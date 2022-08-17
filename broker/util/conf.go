package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"plugin"
	"strings"

	"github.com/jasony62/tms-go-apihub/hub"
	klog "k8s.io/klog/v2"
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
	klog.Infoln("加载API def文件...")
	loadJsonDefData(JSON_TYPE_API, paths[JSON_TYPE_API], "", true)
	klog.Infoln("加载Flow def文件...")
	loadJsonDefData(JSON_TYPE_FLOW, paths[JSON_TYPE_FLOW], "", true)
	klog.Infoln("加载Schedule def文件...")
	loadJsonDefData(JSON_TYPE_SCHEDULE, paths[JSON_TYPE_SCHEDULE], "", true)
	klog.Infoln("加载Private def文件...")
	loadJsonDefData(JSON_TYPE_PRIVATE, paths[JSON_TYPE_PRIVATE], "", true)
	klog.Infoln("加载Rights文件...")
	loadJsonDefData(JSON_TYPE_API_RIGHT, paths[JSON_TYPE_API_RIGHT], "", false)
	loadJsonDefData(JSON_TYPE_FLOW_RIGHT, paths[JSON_TYPE_FLOW_RIGHT], "", false)
	loadJsonDefData(JSON_TYPE_SCHEDULE_RIGHT, paths[JSON_TYPE_SCHEDULE_RIGHT], "", false)
}

func loadJsonDefData(jsonType int, path string, prefix string, includeDir bool) {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		klog.Errorln(err)
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
				klog.Errorln(str)
				panic(str)
			}

			if !json.Valid(byteFile) {
				str := "Json文件无效：" + fileName
				klog.Errorln(str)
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
			case JSON_TYPE_API:
				def := new(hub.HttpApiDef)
				decoder.Decode(&def)
				DefaultConfMap.ApiMap[key] = def
			case JSON_TYPE_FLOW:
				def := new(hub.FlowDef)
				decoder.Decode(&def)
				DefaultConfMap.FlowMap[key] = def
			case JSON_TYPE_SCHEDULE:
				def := new(hub.ScheduleDef)
				decoder.Decode(&def)
				DefaultConfMap.ScheduleMap[key] = def
			case JSON_TYPE_PRIVATE:
				def := new(hub.PrivateArray)
				decoder.Decode(&def)
				DefaultConfMap.PrivateMap[key] = def
			case JSON_TYPE_API_RIGHT:
				def := new(hub.RightArray)
				decoder.Decode(&def)
				DefaultConfMap.ApiRightMap[key] = def
			case JSON_TYPE_FLOW_RIGHT:
				def := new(hub.RightArray)
				decoder.Decode(&def)
				DefaultConfMap.FlowRightMap[key] = def
			case JSON_TYPE_SCHEDULE_RIGHT:
				def := new(hub.RightArray)
				decoder.Decode(&def)
				DefaultConfMap.ScheduleRightMap[key] = def
			default:
			}
		}
	}
}

func loadConfigPluginData(path string) {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		klog.Errorln(err)
		return
	}
	var prefix string
	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())

		if fileInfoList[i].IsDir() {
			klog.Infoln("Plugin子目录: ", fileName)
			prefix = fileInfoList[i].Name()
			loadConfigPluginData(path + "/" + prefix)
		} else {
			if !strings.HasSuffix(fileName, ".so") {
				continue
			}
			klog.Infoln("加载Plugin(*.so)文件: ", fileName)
			p, err := plugin.Open(fileName)
			if err != nil {
				klog.Errorln(err)
				panic(err)
			}

			// 导入动态库，注册函数到funcMap
			registerFunc, err := p.Lookup("Register")
			if err != nil {
				klog.Errorln(err)
				panic(err)
			}
			mapFunc, mapFuncForTemplate := registerFunc.(func() (map[string]interface{}, map[string]interface{}))()
			loadPluginFuncs(mapFunc, mapFuncForTemplate)
			klog.Infof("加载Json文件完成！\r\n")
		}
	}
}

func loadPluginFuncs(mapFunc map[string]interface{}, mapFuncForTemplate map[string]interface{}) {
	for k, v := range mapFunc {
		if _, ok := funcMap[k]; ok {
			klog.Errorf("加载(%s)失败,FuncMap存在重名函数！\r\n", k)
		} else {
			funcMap[k] = v.(hub.FuncHandler)
		}
	}

	for k, v := range mapFuncForTemplate {
		if _, ok := funcMapForTemplate[k]; ok {
			klog.Errorf("加载(%s)失败,FuncMapForTemplate存在重名函数！\r\n", k)
		} else {
			funcMapForTemplate[k] = v
		}
	}
}

func loadTemplateData(path string, prefix string) {
	klog.Infoln("加载Template文件...")
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		str := "invalid path " + path
		klog.Errorln(str)
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
				klog.Errorln(str)
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
	klog.Infoln("CheckRight user:", user, " callType:", callType, " name:", name)
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
		DefaultConfMap.BasePath = path
	}
	klog.Infof("Load main flow from %s\n", DefaultConfMap.BasePath)
	loadJsonDefData(JSON_TYPE_FLOW, DefaultConfMap.BasePath, "", false)
	return nil, http.StatusOK
}

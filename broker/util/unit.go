package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"plugin"
	"strings"

	"github.com/jasony62/tms-go-apihub/hub"
	klog "k8s.io/klog/v2"
)

const (
	JSON_TYPE_PRIVATE = iota
	JSON_TYPE_API
	JSON_TYPE_FLOW
	JSON_TYPE_SCHEDULE
	TMPL_TYPE
)

func FindApiDef(stack *hub.Stack, id string) (*hub.HttpApiDef, error) {
	key := GetBucketKey(stack, id)
	apiDef, ok := hub.DefaultApp.ApiMap[key]
	if !ok {
		return nil, errors.New("No api found")
	}

	return apiDef, nil
}

func FindPrivateDef(stack *hub.Stack, name string) (*hub.PrivateArray, error) {
	if len(name) == 0 {
		return nil, nil
	}

	key := GetBucketKey(stack, name)
	result, ok := hub.DefaultApp.PrivateMap[key]
	if !ok {
		return nil, errors.New("No private found")
	}
	return result, nil
}

func FindFlowDef(stack *hub.Stack, id string) (*hub.FlowDef, error) {
	key := GetBucketKey(stack, id)
	value, ok := hub.DefaultApp.FlowMap[key]
	if !ok {
		return nil, errors.New("No flow found")
	}
	return value, nil
}

func FindScheduleDef(stack *hub.Stack, id string) (*hub.ScheduleDef, error) {
	key := GetBucketKey(stack, id)
	value, ok := hub.DefaultApp.ScheduleMap[key]
	if !ok {
		return nil, errors.New("No schedule found")
	}
	return value, nil
}

func findPrivateValue(private *hub.PrivateArray, name string) string {
	for _, pair := range *private.Pairs {
		if pair.Name == name {
			return pair.Value
		}
	}
	return ""
}

func getArgsVal(stepResult map[string]interface{}, args []string) []string {
	vars := (stepResult["vars"]).(map[string]string)
	argsV := []string{}
	for _, v := range args {
		argsV = append(argsV, vars[v])
	}
	return argsV
}

func GetParameterRawValue(stack *hub.Stack, private *hub.PrivateArray, from *hub.BaseValueDef) (value interface{}, err error) {
	switch from.From {
	case "literal":
		value = from.Content
	case "query":
		value = stack.Query(from.Content)
	case hub.OriginName:
		value, err = stack.QueryFromStepResult("{{.origin." + from.Content + "}}")
	case "private":
		value = findPrivateValue(private, from.Content)
	case "template":
		value, err = stack.QueryFromStepResult(from.Content)
	case "StepResult":
		value, err = stack.QueryFromStepResult("{{." + from.Content + "}}")
	case "json":
		jsonOutBody, err := Json2Json(stack.StepResult, from.Json)
		if err != nil {
			return "", err
		}
		byteJson, err := json.Marshal(jsonOutBody)
		if err != nil {
			return "", err
		}
		value = RemoveOutideQuote(byteJson)
	case "jsonRaw":
		value, err = Json2Json(stack.StepResult, from.Json)
	case "env":
		value = os.Getenv(from.Content)
	case "func":
		function := hub.FuncMap[from.Content]
		if function == nil {
			str := "获取function定义失败："
			klog.Errorln(str)
			panic(str)
		}
		var params []string
		if len(from.Args) > 0 {
			strs := strings.Fields(from.Args)
			params = getArgsVal(stack.StepResult, strs)
		}
		value = function(params)
	case hub.ResultName:
		value, err = stack.QueryFromStepResult("{{.result." + from.Content + "}}")
	default:
		str := "不支持的type" + from.From
		klog.Errorln(str)
		panic(str)
	}
	return
}

func GetParameterStringValue(stack *hub.Stack, private *hub.PrivateArray, from *hub.BaseValueDef) (value string, err error) {
	result, err := GetParameterRawValue(stack, private, from)
	if err == nil {
		return result.(string), err
	}
	return "", err
}

func LoadConfigJsonData(paths []string) {
	hub.DefaultApp.ApiMap = make(map[string]*hub.HttpApiDef)
	hub.DefaultApp.FlowMap = make(map[string]*hub.FlowDef)
	hub.DefaultApp.ScheduleMap = make(map[string]*hub.ScheduleDef)
	hub.DefaultApp.PrivateMap = make(map[string]*hub.PrivateArray)
	hub.DefaultApp.SourceMap = make(map[string]string)

	klog.Infoln("加载API def文件...")
	LoadJsonDefData(JSON_TYPE_API, paths[JSON_TYPE_API], "")
	klog.Infoln("加载Flow def文件...")
	LoadJsonDefData(JSON_TYPE_FLOW, paths[JSON_TYPE_FLOW], "")
	klog.Infoln("加载Schedule def文件...")
	LoadJsonDefData(JSON_TYPE_SCHEDULE, paths[JSON_TYPE_SCHEDULE], "")
	klog.Infoln("加载Private def文件...")
	LoadJsonDefData(JSON_TYPE_PRIVATE, paths[JSON_TYPE_PRIVATE], "")
	klog.Infoln("加载Template文件...")
	LoadTemplateData(TMPL_TYPE, paths[TMPL_TYPE], "")
}

func LoadJsonDefData(jsonType int, path string, prefix string) {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		klog.Errorln(err)
		return
	}

	oldPrefix := prefix
	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())

		if fileInfoList[i].IsDir() {
			prefix = fileInfoList[i].Name()
			LoadJsonDefData(jsonType, path+"/"+prefix, prefix)
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

			if prefix == "" {
				key = fname
			} else {
				key = prefix + "/" + fname
			}

			decoder := json.NewDecoder(bytes.NewReader(byteFile))
			switch jsonType {
			case JSON_TYPE_API:
				def := new(hub.HttpApiDef)
				decoder.Decode(&def)
				hub.DefaultApp.ApiMap[key] = def
			case JSON_TYPE_FLOW:
				def := new(hub.FlowDef)
				decoder.Decode(&def)
				hub.DefaultApp.FlowMap[key] = def
			case JSON_TYPE_SCHEDULE:
				def := new(hub.ScheduleDef)
				decoder.Decode(&def)
				hub.DefaultApp.ScheduleMap[key] = def
			case JSON_TYPE_PRIVATE:
				def := new(hub.PrivateArray)
				decoder.Decode(&def)
				hub.DefaultApp.PrivateMap[key] = def
			default:
			}

			klog.Infof("加载Json文件成功: key: %s\r\n", key)
		}
	}
}

func LoadConfigPluginData(path string) {
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
			LoadConfigPluginData(path + "/" + prefix)
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

			// 导入动态库，注册函数到hub.FuncMap
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

func GetBucketKey(stack *hub.Stack, fileName string) string {
	var bucket string
	var key string
	if hub.DefaultApp.BucketEnable {
		bucket = stack.GinContext.Param(`bucket`)
	}

	if bucket == "" {
		key = fileName
	} else {
		key = bucket + "/" + fileName
	}
	klog.Infof("GetBucketKey key: %s, bucket: %s", key, bucket)
	return key
}

func loadPluginFuncs(mapFunc map[string]interface{}, mapFuncForTemplate map[string]interface{}) {
	for k, v := range mapFunc {
		if _, ok := hub.FuncMap[k]; ok {
			klog.Errorf("加载(%s)失败,FuncMap存在重名函数！\r\n", k)
		} else {
			hub.FuncMap[k] = v.(hub.FuncHandler)
		}
	}

	for k, v := range mapFuncForTemplate {
		if _, ok := hub.FuncMapForTemplate[k]; ok {
			klog.Errorf("加载(%s)失败,FuncMapForTemplate存在重名函数！\r\n", k)
		} else {
			hub.FuncMapForTemplate[k] = v
		}
	}
}

func LoadTemplateData(jsonType int, path string, prefix string) {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		klog.Errorln(err)
		return
	}

	oldPrefix := prefix
	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())

		if fileInfoList[i].IsDir() {
			prefix = fileInfoList[i].Name()
			LoadJsonDefData(jsonType, path+"/"+prefix, prefix)
		} else {
			prefix = oldPrefix

			byteFile, err := ioutil.ReadFile(fileName)
			if err != nil {
				str := "获得tmpl定义失败：" + err.Error()
				klog.Errorln(str)
				panic(str)
			}

			var key string
			fname := fileInfoList[i].Name()

			if prefix == "" {
				key = fname
			} else {
				key = prefix + "/" + fname
			}

			hub.DefaultApp.SourceMap[key] = string(byteFile)
			klog.Infof("加载Template文件成功: key: %s\r\n", key)
		}
	}
}

package unit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/plugin"
	"github.com/jasony62/tms-go-apihub/util"
)

const (
	JSON_TYPE_API      = 0
	JSON_TYPE_FLOW     = 1
	JSON_TYPE_SCHEDULE = 2
	JSON_TYPE_PRIVATE  = 3
)

func loadPrivateData(bucket string, name string) (*hub.PrivateArray, error) {
	key := initBucketKey(bucket, name+".json")
	klog.Infoln("loadPrivateData key: ", key)
	value := hub.DefaultApp.PrivateMap[key]

	return &value, nil
}

func FindApiDef(stack *hub.Stack, id string) (*hub.ApiDef, error) {
	var err error
	key, bucket := GetBucketKey(stack, id+".json")
	value := hub.DefaultApp.ApiMap[key]

	apiDef := &value
	if len(apiDef.PrivateName) > 0 {
		//需要load秘钥
		apiDef.Privates, err = loadPrivateData(bucket, apiDef.PrivateName)
		if err != nil {
			str := "获得Private数据失败：" + err.Error()
			klog.Errorln(str)
			panic(str)
		}
	}

	// 通过插件改写API定义
	if apiDef.Plugins != nil && len(*apiDef.Plugins) > 0 {
		for _, pluginDef := range *apiDef.Plugins {
			if len(pluginDef.Path) > 0 {
				pluginFn, _ := plugin.RewriteApiDef(pluginDef.Path)
				if pluginFn != nil {
					pluginFn(apiDef)
				}
			}
		}
	}

	return apiDef, nil
}

func FindFlowDef(stack *hub.Stack, id string) (*hub.FlowDef, error) {
	key, _ := GetBucketKey(stack, id+".json")
	value := hub.DefaultApp.FlowMap[key]
	return &value, nil
}

func FindScheduleDef(stack *hub.Stack, id string) (*hub.ScheduleDef, error) {
	key, _ := GetBucketKey(stack, id+".json")
	value := hub.DefaultApp.ScheduleMap[key]
	return &value, nil
}

func findPrivateValue(private *hub.PrivateArray, name string) string {
	for _, pair := range *private.Pairs {
		if pair.Name == name {
			return pair.Value
		}
	}
	return ""
}

func GetParameterValue(stack *hub.Stack, private *hub.PrivateArray, from *hub.ApiDefParamFrom) string {
	var value string
	switch from.From {
	case "query":
		value = stack.Query(from.Name)
	case hub.OriginName:
		value = stack.QueryFromStepResult("{{.origin." + from.Name + "}}")
	case "private":
		value = findPrivateValue(private, from.Name)
	case "template":
		value = stack.QueryFromStepResult(from.Name)
	case "StepResult":
		value = stack.QueryFromStepResult("{{." + from.Name + "}}")
	case "JsonTemplate":
		jsonOutBody := util.Json2Json(stack.StepResult, from.Template)
		byteJson, _ := json.Marshal(jsonOutBody)
		value = string(byteJson)
	case "func":
		function := funcMap[from.Name]
		if function != nil {
			value = function()
		} else {
			str := "获取function定义失败："
			klog.Errorln(str)
			panic(str)
		}
	}
	return value
}

func LoadConfigJsonData(ApiDefPath string, FlowDefPath string, ScheduleDefPath string, PrivateDefPath string) {
	hub.DefaultApp.ApiMap = make(map[string]hub.ApiDef)
	hub.DefaultApp.FlowMap = make(map[string]hub.FlowDef)
	hub.DefaultApp.ScheduleMap = make(map[string]hub.ScheduleDef)
	hub.DefaultApp.PrivateMap = make(map[string]hub.PrivateArray)

	klog.Infoln("加载API def文件...")
	LoadJsonDefData(JSON_TYPE_API, ApiDefPath, "")
	klog.Infoln("\r\n")
	klog.Infoln("加载Flow def文件...")
	LoadJsonDefData(JSON_TYPE_FLOW, FlowDefPath, "")
	klog.Infoln("\r\n")
	klog.Infoln("加载Schedule def文件...")
	LoadJsonDefData(JSON_TYPE_SCHEDULE, ScheduleDefPath, "")
	klog.Infoln("\r\n")
	klog.Infoln("加载Private def文件...")
	LoadJsonDefData(JSON_TYPE_PRIVATE, PrivateDefPath, "")
}

func LoadJsonDefData(jsonType int, path string, prefix string) {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		klog.Errorln(err)
		return
	}

	num := len(fileInfoList)

	klog.Infoln("\r\n")
	klog.Infoln("加载Json def文件，本目录文件数: ", num)

	oldPrefix := prefix
	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())
		klog.Infoln("Json file: ", fileName)

		if fileInfoList[i].IsDir() {
			klog.Infoln("Json子目录: ", fileName)
			prefix = fileInfoList[i].Name()
			LoadJsonDefData(jsonType, path+"/"+prefix, prefix)
			klog.Infoln("\r\n")
		} else {
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
			if prefix == "" {
				key = fileInfoList[i].Name()
			} else {
				key = prefix + "/" + fileInfoList[i].Name()
			}

			decoder := json.NewDecoder(bytes.NewReader(byteFile))
			switch jsonType {
			case JSON_TYPE_API:
				def := new(hub.ApiDef)
				decoder.Decode(&def)
				hub.DefaultApp.ApiMap[key] = *def
			case JSON_TYPE_FLOW:
				def := new(hub.FlowDef)
				decoder.Decode(&def)
				hub.DefaultApp.FlowMap[key] = *def
			case JSON_TYPE_SCHEDULE:
				def := new(hub.ScheduleDef)
				decoder.Decode(&def)
				hub.DefaultApp.ScheduleMap[key] = *def
			case JSON_TYPE_PRIVATE:
				def := new(hub.PrivateArray)
				decoder.Decode(&def)
				hub.DefaultApp.PrivateMap[key] = *def
			default:
			}

			klog.Infof("加载Json文件成功: 文件名和map的key: %s\r\n", key)
		}
	}
}

func initBucketKey(bucket string, fileName string) string {
	var key string
	if bucket == "" {
		key = fileName
	} else {
		key = bucket + "/" + fileName
	}
	return key
}

func GetBucketKey(stack *hub.Stack, fileName string) (string, string) {
	var bucket string
	if hub.DefaultApp.BucketEnable {
		bucket = stack.GinContext.Param(`bucket`)
	}

	var key string
	if bucket == "" {
		key = fileName
	} else {
		key = bucket + "/" + fileName
	}
	klog.Infof("GetBucketKey key: %s, bucket: %s", key, bucket)
	return key, bucket
}

package unit

import (
	"encoding/json"
	"fmt"
	"os"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/plugin"
	"github.com/jasony62/tms-go-apihub/util"
)

func loadDefFromFile(stack *hub.Stack, path string, id string, def interface{}) {
	var filePath string
	var bucket string
	if hub.DefaultApp.BucketEnable {
		bucket = stack.GinContext.Param(`bucket`)
	}
	if bucket != "" {
		filePath = fmt.Sprintf("%s/%s/%s.json", path, bucket, id)
	} else {
		filePath = fmt.Sprintf("%s/%s.json", path, id)
	}
	filePtr, err := os.Open(filePath)
	if err != nil {
		str := "获得API定义失败：" + err.Error()
		klog.Errorln(str)
		panic(str)
	}
	defer filePtr.Close()

	decoder := json.NewDecoder(filePtr)
	decoder.Decode(def)
}

func FindApiDef(stack *hub.Stack, id string) (*hub.ApiDef, error) {
	apiDef := new(hub.ApiDef)
	loadDefFromFile(stack, hub.DefaultApp.ApiDefPath, id, apiDef)

	if len(apiDef.PrivateName) > 0 {
		//需要load秘钥
		apiDef.Privates = new(hub.PrivateArray)
		loadDefFromFile(stack, hub.DefaultApp.PrivateDefPath, apiDef.PrivateName, apiDef.Privates)
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
	scheduleDef := new(hub.FlowDef)
	loadDefFromFile(stack, hub.DefaultApp.FlowDefPath, id, scheduleDef)
	return scheduleDef, nil
}

func FindScheduleDef(stack *hub.Stack, id string) (*hub.ScheduleDef, error) {
	scheduleDef := new(hub.ScheduleDef)
	loadDefFromFile(stack, hub.DefaultApp.ScheduleDefPath, id, scheduleDef)
	return scheduleDef, nil
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

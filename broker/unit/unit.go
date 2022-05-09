package unit

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	klog "k8s.io/klog/v2"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/plugin"
	"github.com/jasony62/tms-go-apihub/util"
)

func loadPrivateData(path string, bucket string, name string) (*hub.PrivateArray, error) {
	var filePath string
	if bucket != "" {
		filePath = fmt.Sprintf("%s/%s/%s.json", path, bucket, name)
	} else {
		filePath = fmt.Sprintf("%s/%s.json", path, name)
	}
	filePtr, err := os.Open(filePath)
	if err != nil {
		str := "获得API定义失败：" + err.Error()
		klog.Errorln(str)
		panic(str)
	}
	defer filePtr.Close()

	result := new(hub.PrivateArray)
	decoder := json.NewDecoder(filePtr)
	decoder.Decode(result)
	return result, nil
}

func FindApiDef(stack *hub.Stack, id string) (*hub.ApiDef, error) {
	var filePath string
	var bucket string
	if hub.DefaultApp.BucketEnable {
		bucket = stack.GinContext.Param(`bucket`)
	}
	if bucket != "" {
		filePath = fmt.Sprintf("%s/%s/%s.json", hub.DefaultApp.ApiDefPath, bucket, id)
	} else {
		filePath = fmt.Sprintf("%s/%s.json", hub.DefaultApp.ApiDefPath, id)
	}
	filePtr, err := os.Open(filePath)
	if err != nil {
		str := "获得API定义失败：" + err.Error()
		klog.Errorln(str)
		panic(str)
	}
	defer filePtr.Close()

	apiDef := new(hub.ApiDef)
	decoder := json.NewDecoder(filePtr)
	decoder.Decode(apiDef)

	if len(apiDef.PrivateName) > 0 {
		//需要load秘钥
		apiDef.Privates, err = loadPrivateData(hub.DefaultApp.PrivateDefPath, bucket, apiDef.PrivateName)
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
	var filePath string

	var bucket string
	if hub.DefaultApp.BucketEnable {
		bucket = stack.GinContext.Param(`bucket`)
	}

	if bucket != "" {
		filePath = fmt.Sprintf("%s/%s/%s.json", hub.DefaultApp.FlowDefPath, bucket, id)
	} else {
		filePath = fmt.Sprintf("%s/%s.json", hub.DefaultApp.FlowDefPath, id)
	}
	filePtr, err := os.Open(filePath)
	if err != nil {
		str := "获得Flow定义失败：" + filePath + "due to" + err.Error()
		klog.Errorln(str)
		panic(str)
	}
	defer filePtr.Close()

	flowDef := new(hub.FlowDef)
	decoder := json.NewDecoder(filePtr)
	decoder.Decode(&flowDef)

	return flowDef, nil
}

func FindScheduleDef(stack *hub.Stack, id string) (*hub.ScheduleDef, error) {
	var filePath string
	var bucket string

	if hub.DefaultApp.BucketEnable {
		bucket = stack.GinContext.Param(`bucket`)
	}

	if bucket != "" {
		filePath = fmt.Sprintf("%s/%s/%s.json", hub.DefaultApp.ScheduleDefPath, bucket, id)
	} else {
		filePath = fmt.Sprintf("%s/%s.json", hub.DefaultApp.ScheduleDefPath, id)
	}
	filePtr, err := os.Open(filePath)
	if err != nil {
		str := "获得Schedule定义失败：" + err.Error()
		klog.Errorln(str)
		panic(str)
	}
	defer filePtr.Close()

	scheduleDef := new(hub.ScheduleDef)
	decoder := json.NewDecoder(filePtr)
	decoder.Decode(&scheduleDef)

	return scheduleDef, nil
}

func FindPrivateValue(api *hub.ApiDef, name string) string {
	for _, pair := range *api.Privates.Pairs {
		if pair.Name == name {
			return pair.Value
		}
	}
	return ""
}

// 在调用流中根据规则改写api定义
func RewriteApiDefInFlow(api *hub.ApiDef, flowApi *hub.FlowStepApiDef) error {
	if len(*flowApi.Parameters) > 0 {
		for _, newParam := range *flowApi.Parameters {
			for i, oldParam := range *api.Parameters {
				if newParam.Name == oldParam.Name {
					(*api.Parameters)[i] = newParam
				}
			}
		}
	}

	return nil
}
func getArgsVal(stepResult map[string]interface{}, args []string) []string {
	vars := (stepResult["vars"]).(map[string]string)
	argsV := []string{}
	for _, v := range args {
		argsV = append(argsV, vars[v])
	}
	return argsV
}

func GetParameterValue(stack *hub.Stack, from *hub.ApiDefParamFrom) string {
	var value string
	switch from.From {
	case "query":
		value = stack.Query(from.Name)
	case "origin":
		value = stack.QueryFromStepResult("{{.origin." + from.Name + "}}")
	case "private":
		value = FindPrivateValue(stack.ApiDef, from.Name)
	case "template":
		value = stack.QueryFromStepResult(from.Name)
	case "StepResult":
		value = stack.QueryFromStepResult("{{." + from.Name + "}}")
	case "JsonTemplate":
		jsonOutBody := util.Json2Json(stack.StepResult, from.Template)
		byteJson, _ := json.Marshal(jsonOutBody)
		value = util.RemoveOutideQuote(byteJson)
	case "func":
		function := hub.FuncMap[from.Name]
		if function == nil {
			str := "获取function定义失败："
			klog.Errorln(str)
			panic(str)
		}
		switch funcV := function.(type) {
		case func() string:
			value = funcV()
			break
		case func([]string) string:
			strs := strings.Fields(from.Args)
			argsV := getArgsVal(stack.StepResult, strs)
			value = funcV(argsV)
		default:
			str := "function不能执行"
			klog.Errorln(str)
			panic(str)
		}
	}
	return value
}

package util

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"text/template"

	"github.com/jasony62/tms-go-apihub/hub"
	klog "k8s.io/klog/v2"
)

const (
	JSON_TYPE_PRIVATE = iota
	JSON_TYPE_API
	JSON_TYPE_FLOW
	JSON_TYPE_SCHEDULE
	JSON_TYPE_API_RIGHT
	JSON_TYPE_FLOW_RIGHT
	JSON_TYPE_SCHEDULE_RIGHT
)

// 从执行结果中获取查询参数
func queryFromStepResult(stack *hub.Stack, name string) (string, error) {
	tmpl, err := template.New("key").Funcs(funcMapForTemplate).Parse(name)
	if err != nil {
		klog.Infoln("queryFromStepResult 创建并解析template失败:", err)
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, stack.StepResult)
	if err != nil {
		klog.Infoln("渲染template失败:", err)
		return "", err
	}
	return buf.String(), err
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
	case "header":
		value = stack.GinContext.GetHeader(from.Content)
	case "query":
		value = stack.Query(from.Content)
	case "header":
		stack.GinContext.Request.Header.Get(from.Content)
	case hub.OriginName:
		value, err = queryFromStepResult(stack, "{{.origin."+from.Content+"}}")
	case "private":
		value = findPrivateValue(private, from.Content)
	case "template":
		value, err = queryFromStepResult(stack, from.Content)
	case "StepResult":
		value, err = queryFromStepResult(stack, "{{."+from.Content+"}}")
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
		function := funcMap[from.Content]
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
		value, err = queryFromStepResult(stack, "{{.result."+from.Content+"}}")
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

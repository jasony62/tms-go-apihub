package util

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"text/template"

	"github.com/jasony62/tms-go-apihub/hub"
	"go.uber.org/zap"
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
func queryFromHeap(stack *hub.Stack, name string) (string, error) {
	tmpl, err := template.New("key").Funcs(funcMapForTemplate).Parse(name)
	if err != nil {
		zap.S().Errorln("NOK 创建并解析template失败:", err)
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, stack.Heap)
	if err != nil {
		zap.S().Errorln("NOK execute template:", err, " name:", name, "heap:", stack.Heap)
		return "", err
	}
	return buf.String(), err
}

func findPrivateValue(private *hub.PrivateArray, name string) string {
	if private == nil {
		return ""
	}

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
		// 从请求参数中获取查询参数
		value = stack.GinContext.Query(from.Content)
	case hub.HeapOriginName:
		value, err = queryFromHeap(stack, "{{.origin."+from.Content+"}}")
	case "private":
		value = findPrivateValue(private, from.Content)
	case "template":
		value, err = queryFromHeap(stack, from.Content)
	case "heap":
		value, err = queryFromHeap(stack, "{{."+from.Content+"}}")
	case "json":
		jsonOutBody, err := Json2Json(stack.Heap, from.Json)
		if err != nil {
			return "", err
		}
		byteJson, err := json.Marshal(jsonOutBody)
		if err != nil {
			return "", err
		}
		value = RemoveOutideQuote(byteJson)
	case "jsonRaw":
		value, err = Json2Json(stack.Heap, from.Json)
	case "env":
		value = os.Getenv(from.Content)
	case "func":
		function := funcMap[from.Content]
		if function == nil {
			str := "获取function定义失败："
			zap.S().Errorln(str)
			panic(str)
		}
		var params []string
		if len(from.Args) > 0 {
			strs := strings.Fields(from.Args)
			params = getArgsVal(stack.Heap, strs)
		}
		value = function(params)
	default:
		str := "不支持的type " + from.From
		zap.S().Errorln(str)
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

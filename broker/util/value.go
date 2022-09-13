package util

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"text/template"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/logger"
)

func executeTemplate(source interface{}, rules interface{}) (*bytes.Buffer, error) {
	byteTempl, err := json.Marshal(rules)
	if err != nil {
		return nil, err
	}

	strTempl := string(byteTempl)

	// 处理数组
	strTempl = strings.ReplaceAll(strTempl, "\"{{range", "{{range")
	strTempl = strings.ReplaceAll(strTempl, "end}}\"", "end}}")
	strTempl = strings.ReplaceAll(strTempl, "\\\"", "\"")

	tmpl, err := template.New("json").Funcs(funcMapForTemplate).Parse(strTempl)
	if err != nil {
		logger.LogS().Infoln("get template result：", strTempl, byteTempl, " error: ", err)
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, source)
	if err != nil {
		logger.LogS().Infoln("get template result：", err)
		return nil, err
	}
	return buf, err
}

func json2Json(source interface{}, rules interface{}) (interface{}, error) {
	var target interface{}
	buf, err := executeTemplate(source, rules)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf.Bytes(), &target)
	return target, err
}

func rvemoveOutideQuote(s []byte) string {
	if len(s) > 0 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1:(len(s) - 1)]
	}
	return string(s)
}

// 从执行结果中获取查询参数
func queryFromHeap(stack *hub.Stack, name string) (string, error) {
	tmpl, err := template.New("key").Funcs(funcMapForTemplate).Parse(name)
	if err != nil {
		logger.LogS().Errorln("NOK 创建并解析template失败:", err)
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, stack.Heap)
	if err != nil {
		logger.LogS().Errorln("NOK execute template:", err, " name:", name, "heap:", stack.Heap)
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
		jsonOutBody, err := json2Json(stack.Heap, from.Json)
		if err != nil {
			return "", err
		}
		byteJson, err := json.Marshal(jsonOutBody)
		if err != nil {
			return "", err
		}
		value = rvemoveOutideQuote(byteJson)
	case "jsonRaw":
		value, err = json2Json(stack.Heap, from.Json)
	case "env":
		value = os.Getenv(from.Content)
	case "func":
		function := funcMap[from.Content]
		if function == nil {
			str := "获取function定义失败："
			logger.LogS().Errorln(str)
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
		logger.LogS().Errorln(str)
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

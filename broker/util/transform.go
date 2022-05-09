package util

import (
	"bytes"
	"encoding/json"
	"github.com/jasony62/tms-go-apihub/hub"
	"strings"
	"text/template"
)

func executeTemplate(source interface{}, rules interface{}, translate bool) *bytes.Buffer {
	byteTempl, _ := json.Marshal(rules)
	strTempl := string(byteTempl)
	if translate {
		// 处理数组
		strTempl = strings.ReplaceAll(strTempl, "\"{{range", "{{range")
		strTempl = strings.ReplaceAll(strTempl, "end}}\"", "end}}")
		strTempl = strings.ReplaceAll(strTempl, "\\\"", "\"")
	}

	tmpl, _ := template.New("json").Funcs(hub.FuncMap).Parse(strTempl)
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, source)
	return buf
}

func Json2Json(source interface{}, rules interface{}) interface{} {
	buf := executeTemplate(source, rules, true)
	var target interface{}
	json.Unmarshal(buf.Bytes(), &target)
	return target
}

func RemoveOutideQuote(s []byte) string {
	if len(s) > 0 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1:(len(s) - 1)]
	}
	return string(s)
}

//func HandleTemplate(source interface{}, rules interface{}) string {
///	buf := executeTemplate(source, rules, false)
//	return buf.String()
//}

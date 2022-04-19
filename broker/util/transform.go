package util

import (
	"bytes"
	"encoding/json"
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

	tmpl, _ := template.New("json").Parse(strTempl)
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

func HandleTemplate(source interface{}, rules interface{}) string {
	buf := executeTemplate(source, rules, false)
	return buf.String()
}

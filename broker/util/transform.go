package util

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"
)

func Json2Json(source interface{}, rules interface{}) interface{} {
	byteTempl, _ := json.Marshal(rules)
	strTempl := string(byteTempl)
	// 处理数组
	strTempl = strings.ReplaceAll(strTempl, "\"{{range", "{{range")
	strTempl = strings.ReplaceAll(strTempl, "end}}\"", "end}}")
	strTempl = strings.ReplaceAll(strTempl, "\\\"", "\"")

	tmpl, _ := template.New("json").Parse(strTempl)

	buf := new(bytes.Buffer)

	tmpl.Execute(buf, source)

	var target interface{}
	json.Unmarshal(buf.Bytes(), &target)

	return target
}

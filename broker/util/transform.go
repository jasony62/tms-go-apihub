package util

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/jasony62/tms-go-apihub/hub"
	klog "k8s.io/klog/v2"
)

func executeTemplate(source interface{}, rules interface{}) *bytes.Buffer {
	byteTempl, _ := json.Marshal(rules)

	strTempl := string(byteTempl)

	// 处理数组
	strTempl = strings.ReplaceAll(strTempl, "\"{{range", "{{range")
	strTempl = strings.ReplaceAll(strTempl, "end}}\"", "end}}")
	strTempl = strings.ReplaceAll(strTempl, "\\\"", "\"")

	tmpl, err := template.New("json").Funcs(hub.FuncMapForTemplate).Parse(strTempl)
	if err != nil {
		klog.Infoln("get template result：", strTempl, byteTempl, " error: ", err)
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, source)
	if err != nil {
		klog.Infoln("get template result：", err)
	}
	return buf
}

func Json2Json(source interface{}, rules interface{}) interface{} {
	buf := executeTemplate(source, rules)
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

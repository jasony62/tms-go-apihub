package util

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"go.uber.org/zap"
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
		zap.S().Infoln("get template result：", strTempl, byteTempl, " error: ", err)
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, source)
	if err != nil {
		zap.S().Infoln("get template result：", err)
		return nil, err
	}
	return buf, err
}

func Json2Json(source interface{}, rules interface{}) (interface{}, error) {
	var target interface{}
	buf, err := executeTemplate(source, rules)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buf.Bytes(), &target)
	return target, err
}

func RemoveOutideQuote(s []byte) string {
	if len(s) > 0 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1:(len(s) - 1)]
	}
	return string(s)
}

func Json2Html(source interface{}, rules string) (string, error) {
	tmpl, err := template.New("tmpl").Parse(rules)
	if err != nil {
		zap.S().Infoln("get template result：", rules, " error: ", err)
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, source)
	if err != nil {
		zap.S().Infoln("get template result：", err)
		return "", err
	}

	return buf.String(), err
}

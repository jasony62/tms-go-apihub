package apis

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/logger"
	"github.com/jasony62/tms-go-apihub/util"
)

func json2Html(source interface{}, rules string) (string, error) {
	tmpl, err := template.New("tmpl").Parse(rules)
	if err != nil {
		logger.LogS().Errorln("get template result：", rules, " error: ", err)
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, source)
	if err != nil {
		logger.LogS().Errorln("get template result：", err)
		return "", err
	}

	return buf.String(), err
}

func createHtml(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		str := "createHtml,缺少参数"
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	name, OK := params["type"]
	if !OK {
		str := "createHtml,type为空"
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	content, OK := params["content"]
	if !OK {
		str := "createHtml,content为空"
		logger.LogS().Errorln(stack.BaseString, str)
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	if name == "resource" {
		content, OK = util.FindResourceDef(content)
		if !OK {
			str := "createHtml FindResourceDef failed"
			logger.LogS().Errorln(stack.BaseString, str)
			return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
		}
	}

	result, err := json2Html(stack.Heap, content)
	if err != nil {
		return util.CreateTmsError(hub.TmsErrorApisId, err.Error(), nil), http.StatusInternalServerError
	}
	return result, 200
}

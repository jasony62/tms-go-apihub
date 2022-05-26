package apis

import (
	"net/http"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
)

func createHtml(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		return nil, http.StatusInternalServerError
	}

	name, OK := params["type"]
	if !OK {
		return nil, http.StatusInternalServerError
	}

	content, OK := params["content"]
	if !OK {
		return nil, http.StatusInternalServerError
	}

	if name == "resource" {
		content, OK = hub.DefaultApp.SourceMap[content]
		if !OK {
			return nil, http.StatusInternalServerError
		}
	}

	result, err := util.Json2Html(stack.StepResult, content)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return result, 200
}

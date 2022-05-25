package apis

import (
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
)

func createHtml(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		return nil, 500
	}

	name, OK := params["type"]
	if !OK {
		return nil, 500
	}

	content, OK := params["content"]
	if !OK {
		return nil, 500
	}

	if name == "resource" {
		content, OK = hub.DefaultApp.SourceMap[content]
		if !OK {
			return nil, 500
		}
	}

	result := util.Json2Html(stack.StepResult, content)
	return result, 200
}

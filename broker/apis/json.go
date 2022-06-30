package apis

import (
	"net/http"

	"github.com/jasony62/tms-go-apihub/hub"
)

func createJson(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		return nil, http.StatusInternalServerError
	}

	key, OK := params["key"]
	if !OK {
		return nil, http.StatusInternalServerError
	}
	tmp := stack.Heap[hub.OriginName].(map[string]interface{})
	result := tmp[key]
	delete(tmp, key)
	return result, http.StatusOK
}

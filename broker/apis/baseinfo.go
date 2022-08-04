package apis

import (
	"github.com/google/uuid"
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
)

func fillBaseInfo(stack *hub.Stack, params map[string]string) (interface{}, int) {
	base := stack.Heap[hub.HeapBaseName].(map[string]interface{})

	if base == nil {
		return nil, 500

	}

	for k, v := range params {
		base[k] = v
	}

	if len(params["uuid"]) == 0 {
		base["uuid"] = uuid.New().String()
	}

	stack.BaseString = util.MapToString(base)
	return nil, 200
}

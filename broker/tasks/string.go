package tasks

import (
	"github.com/jasony62/tms-go-apihub/hub"
)

func checkStringsEqual(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		return nil, 500
	}

	for k, v := range params {
		if k != v {
			return nil, 500
		}
	}
	return nil, 200
}

func checkStringsNotEqual(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		return nil, 500
	}

	for k, v := range params {
		if k == v {
			return nil, 500
		}
	}
	return nil, 200
}

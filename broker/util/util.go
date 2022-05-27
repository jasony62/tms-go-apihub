package util

import "github.com/jasony62/tms-go-apihub/hub"

func init() {
	hub.FuncMap = map[string]hub.FuncHandler{
		"utc": utc,
		"md5": md5Func,
	}
	hub.FuncMapForTemplate = map[string](interface{}){
		"utc": utcTemplate,
		"md5": md5Template,
	}
}

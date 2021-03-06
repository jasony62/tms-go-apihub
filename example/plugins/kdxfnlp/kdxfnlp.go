package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"tms-go-apihub/broker/hub"
)

func init() {
	log.Println("插件【科大讯飞NLP】执行初始化")

}

func Register() (map[string]hub.FuncHandler, map[string](interface{})) {
	funcMap := map[string](interface{}){
		"md5Func": md5Func,
	}
	funcMapForTemplate := map[string](interface{}){
		"md5Func": md5Template,
	}
	return funcMap, funcMapForTemplate
}

func md5Func(params []string) string {
	if len(params) == 0 {
		return ""
	}
	var str string
	for _, v := range params {
		str = str + v
	}
	w := md5.New()
	io.WriteString(w, str)
	checksum := fmt.Sprintf("%x", w.Sum(nil))
	return checksum
}

func md5Template(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	str := fmt.Sprint(args...)
	w := md5.New()
	io.WriteString(w, str)
	checksum := fmt.Sprintf("%x", w.Sum(nil))
	return checksum
}

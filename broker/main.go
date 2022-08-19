package main

import (
	"flag"
	"net/http"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/jasony62/tms-go-apihub/apis"
	"github.com/jasony62/tms-go-apihub/core"
	"github.com/jasony62/tms-go-apihub/hub"
)

// 命令行指定的环境变量文件
var envfile string
var basePath string

func init() {

	flag.StringVar(&envfile, "env", "", "指定环境变量文件")
	flag.StringVar(&basePath, "base", "../example/", "指定启动路径")
}

func welcome(stack *hub.Stack, params map[string]string) (interface{}, int) {
	content, OK := params["content"]
	if OK {
		zap.S().Infoln(content)
	}
	return nil, http.StatusOK
}

func main() {
	core.RegisterApis(map[string]hub.ApiHandler{"welcome": welcome})
	apis.ApisInit()
	flag.Parse()
	if envfile != "" {
		err := godotenv.Load(envfile)
		if err != nil {
			zap.S().Errorln(err.Error())
		}
	}

	core.ApiHubStartMainFlow(basePath)
}

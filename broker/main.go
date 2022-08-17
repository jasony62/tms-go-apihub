package main

import (
	"flag"
	"net/http"

	"github.com/joho/godotenv"
	klog "k8s.io/klog/v2"

	_ "github.com/jasony62/tms-go-apihub/apis"
	"github.com/jasony62/tms-go-apihub/core"
	"github.com/jasony62/tms-go-apihub/hub"
)

// 命令行指定的环境变量文件
var envfile string
var basePath string

func init() {
	klog.InitFlags(nil)

	flag.StringVar(&envfile, "env", "", "指定环境变量文件")
	flag.StringVar(&basePath, "base", "./conf/", "指定启动路径")
}

func welcome(stack *hub.Stack, params map[string]string) (interface{}, int) {
	content, OK := params["content"]
	if OK {
		klog.Errorln(content)
	}
	return nil, http.StatusOK
}

func main() {
	core.RegisterApis(map[string]hub.ApiHandler{"welcome": welcome})
	flag.Parse()
	if envfile != "" {
		err := godotenv.Load(envfile)
		if err != nil {
			klog.Fatal(err)
		}
	}

	defer klog.Flush()
	core.ApiHubStartMainFlow(basePath)
}

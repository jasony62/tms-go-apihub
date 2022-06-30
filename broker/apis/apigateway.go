package apis

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/jasony62/tms-go-apihub/core"
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"

	"github.com/gin-gonic/gin"
	klog "k8s.io/klog/v2"
)

// 应用的基本信息
type app struct {
	Host         string
	Port         int
	BucketEnable bool
}

var defaultApp = app{
	Host:         "0.0.0.0",
	Port:         8080,
	BucketEnable: false,
}

// 1次请求的上下文
func newStack(c *gin.Context, level string) *hub.Stack {
	// 收到的数据
	var value interface{}
	inReqData := new(interface{})
	c.ShouldBindJSON(&inReqData)

	if *inReqData == nil {
		value = make(map[string]interface{})
	} else {
		value = *inReqData
	}

	base := make(map[string]string)
	name := c.Param(`Id`)
	version := c.Param(`version`)
	if len(version) > 0 {
		name = name + "_" + version
	}
	if defaultApp.BucketEnable {
		name = c.Param(`bucket`) + "/" + name
	}

	base["root"] = name
	base["type"] = level
	base["start"] = strconv.FormatInt(time.Now().Unix(), 10)

	return &hub.Stack{
		GinContext: c,
		Heap:       map[string]interface{}{hub.OriginName: value, hub.BaseName: base},
		Base:       base,
	}
}

// 执行1个API调用
func callHttpApi(c *gin.Context) {
	// 调用api
	tmpStack := newStack(c, "httpapi")

	params := []hub.BaseParamDef{{Name: "name", Value: hub.BaseValueDef{From: "literal", Content: "_HTTPAPI"}}}

	result, status := core.ApiRun(tmpStack, &hub.ApiDef{Name: "main", Command: "flowApi", Args: &params, ResultKey: "main"}, "", false)
	if status != http.StatusOK {
		//成功时的回复应该定义在flow的step中
		c.IndentedJSON(status, result)
	}
}

// 执行一个调用流程
func callFlow(c *gin.Context) {
	// 执行编排
	tmpStack := newStack(c, "flow")
	params := []hub.BaseParamDef{{Name: "name", Value: hub.BaseValueDef{From: "literal", Content: tmpStack.Base[hub.RootParamName]}}}

	result, status := core.ApiRun(tmpStack, &hub.ApiDef{Name: "main", Command: "flowApi", Args: &params, ResultKey: "main"}, "", false)
	if status != http.StatusOK {
		//成功时的回复应该定义在flow的step中
		c.IndentedJSON(status, result)
	}
}

// 执行一个计划流程
func callSchedule(c *gin.Context) {
	// 执行编排
	tmpStack := newStack(c, "schedule")
	params := []hub.BaseParamDef{{Name: "name", Value: hub.BaseValueDef{From: "literal", Content: tmpStack.Base[hub.RootParamName]}}}

	result, status := core.ApiRun(tmpStack, &hub.ApiDef{Name: "main", Command: "scheduleApi", Args: &params, ResultKey: "main"}, "", false)
	if status != http.StatusOK {
		//成功时的回复应该定义在flow的step中
		c.IndentedJSON(status, result)
	}
}

func apiGatewayRun(host string, port string, BucketEnable string) {
	if len(host) > 0 {
		defaultApp.Host = host
	}
	klog.Infoln("host: ", defaultApp.Host)

	if len(port) > 0 {
		defaultApp.Port, _ = strconv.Atoi(port)
	}
	klog.Infoln("port ", defaultApp.Port)

	if len(BucketEnable) > 0 {
		re := regexp.MustCompile(`(?i)yes|true`)
		defaultApp.BucketEnable = re.MatchString(BucketEnable)
	}
	klog.Infoln("bucket enable ", defaultApp.BucketEnable)

	router := gin.Default()
	if defaultApp.BucketEnable {
		router.Any("/httpapi/:bucket/:Id", callHttpApi)
		router.Any("/httpapi/:bucket/:Id/:version", callHttpApi)
		router.Any("/flow:bucket/:Id", callFlow)
		router.Any("/flow:bucket/:Id/:version", callFlow)
		router.Any("/schedule:bucket/:Id", callSchedule)
		router.Any("/schedule:bucket/:Id/:version", callSchedule)
	} else {
		router.Any("/httpapi/:Id", callHttpApi)
		router.Any("/httpapi/:Id/:version", callHttpApi)
		router.Any("/flow/:Id", callFlow)
		router.Any("/flow/:Id/:version", callFlow)
		router.Any("/schedule/:Id", callSchedule)
		router.Any("/schedule/:Id/:version", callSchedule)
	}
	basePath := util.GetBasePath() + "templates"
	if needLoad, _ := util.PathExists(basePath); needLoad {
		router.LoadHTMLGlob(basePath + "/*.tmpl")
	}

	if defaultApp.Port > 0 {
		router.Run(fmt.Sprintf("%s:%d", defaultApp.Host, defaultApp.Port))
	} else {
		router.Run(defaultApp.Host)
	}
}

func apiGateway(stack *hub.Stack, params map[string]string) (interface{}, int) {
	host := params["host"]
	port := params["port"]
	bucket := params["bucket"]

	apiGatewayRun(host, port, bucket)
	return nil, 200
}

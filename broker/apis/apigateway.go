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
	bucketEnable bool
	pre          string
	httpApi      string
	postOK       string
	postNOK      string
}

var defaultApp = app{
	bucketEnable: false,
	pre:          "_APIGATEWAY_PRE",
	httpApi:      "_APIGATEWAY_HTTPAPI",
	postOK:       "_APIGATEWAY_POST_OK",
	postNOK:      "_APIGATEWAY_POST_NOK",
}

func fillStats(stack *hub.Stack, result interface{}, code int) {
	stats := make(map[string]string)
	stack.Heap["stats"] = stats

	stats["child"] = stack.Base["root"]
	stats["duration"] = strconv.FormatFloat(time.Since(stack.Now).Seconds(), 'f', 5, 64)
	stats["code"] = strconv.FormatInt(int64(code), 10)

	if code == http.StatusOK {
		stats["id"] = "0"
		stats["msg"] = "ok"
		klog.Infoln("!!!!post apigateway OK:", stack.Base, "result:", result, " code:", code, " stats:", stats)
		params := []hub.BaseParamDef{{Name: "name", Value: hub.BaseValueDef{From: "literal", Content: "_HTTPOK"}}}
		core.ApiRun(stack, &hub.ApiDef{Name: "HTTPAPI_POST_OK", Command: "flowApi", Args: &params}, "", true)
	} else {
		/*TODO real value*/
		stats["id"] = strconv.FormatInt(int64(code), 10)
		stats["msg"] = "err"
		klog.Errorln("!!!!post apigateway NOK:", stack.Base, ", result:", result, " code:", code, " stats:", stats)
		params := []hub.BaseParamDef{{Name: "name", Value: hub.BaseValueDef{From: "literal", Content: "_HTTPNOK"}}}
		core.ApiRun(stack, &hub.ApiDef{Name: "HTTPAPI_POST_NOK", Command: "flowApi", Args: &params}, "", true)
	}
}

// 1次请求的上下文
func newStack(c *gin.Context, level string) *hub.Stack {
	now := time.Now()
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
	if defaultApp.bucketEnable {
		name = c.Param(`bucket`) + "/" + name
	}

	base["root"] = name
	base["type"] = level
	base["start"] = strconv.FormatInt(now.Unix(), 10)

	return &hub.Stack{
		GinContext: c,
		Heap:       map[string]interface{}{hub.OriginName: value, hub.BaseName: base},
		Base:       base,
		Now:        now,
	}
}

func callCommon(stack *hub.Stack, command string, content string) {
	params := []hub.BaseParamDef{{Name: "name", Value: hub.BaseValueDef{From: "literal", Content: ""}}}
	var result interface{}
	var status int
	if len(defaultApp.pre) != 0 {
		// 调用api
		params[0].Value.Content = defaultApp.pre
		result, status = core.ApiRun(stack, &hub.ApiDef{Name: "main_pre", Command: "flowApi", Args: &params}, "", false)
		if status != http.StatusOK {
			//成功时的回复应该定义在flow的step中
			stack.GinContext.IndentedJSON(status, result)
			klog.Errorln("PRE status:", status, " result:", result)

			if len(defaultApp.postNOK) != 0 {
				fillStats(stack, result, status)
				params[0].Value.Content = defaultApp.postNOK
				result, status = core.ApiRun(stack, &hub.ApiDef{Name: "main_pre_post_nok", Command: "flowApi", Args: &params}, "", true)
				if status != http.StatusOK {
					klog.Errorln("PRE - post NOK status:", status, " result:", result)
				}
			}
			return
		}
	}
	// 调用api
	params[0].Value.Content = content
	result, status = core.ApiRun(stack, &hub.ApiDef{Name: "main", Command: command, Args: &params}, "", false)
	fillStats(stack, result, status)
	if status != http.StatusOK {
		//成功时的回复应该定义在flow的step中
		stack.GinContext.IndentedJSON(status, result)
		klog.Errorln("common status:", status, " result:", result)

		if len(defaultApp.postNOK) != 0 {
			params[0].Value.Content = defaultApp.postNOK
			result, status = core.ApiRun(stack, &hub.ApiDef{Name: "main_post_nok", Command: "flowApi", Args: &params}, "", true)
			if status != http.StatusOK {
				klog.Errorln("common - post NOK - NOK status:", status, " result:", result)
			}
		}
	} else if len(defaultApp.postOK) != 0 {
		params[0].Value.Content = defaultApp.postOK
		result, status = core.ApiRun(stack, &hub.ApiDef{Name: "main_post_ok", Command: "flowApi", Args: &params}, "", true)
		if status != http.StatusOK {
			klog.Errorln("common - post OKstatus:", status, " result:", result)
		}
	}
}

// 执行1个API调用
func callHttpApi(c *gin.Context) {
	stack := newStack(c, "httpapi")
	callCommon(stack, "flowApi", defaultApp.httpApi)
}

// 执行一个调用流程
func callFlow(c *gin.Context) {
	stack := newStack(c, "flow")
	// 执行编排
	callCommon(stack, "flowApi", stack.Base[hub.RootParamName])
}

// 执行一个计划流程
func callSchedule(c *gin.Context) {
	stack := newStack(c, "schedule")
	// 执行编排
	callCommon(stack, "scheduleApi", stack.Base[hub.RootParamName])
}

func apiGatewayRun(host string, portString string, bucketEnable string,
	pre string, postOK string, postNOK string, httpApi string) {
	var port int
	if len(host) == 0 {
		host = "0.0.0.0"
	}

	klog.Infoln("host: ", host)

	if len(portString) > 0 {
		port, _ = strconv.Atoi(portString)
	}
	klog.Infoln("port ", port)

	if len(bucketEnable) > 0 {
		re := regexp.MustCompile(`(?i)yes|true`)
		defaultApp.bucketEnable = re.MatchString(bucketEnable)
	}
	klog.Infoln("bucket enable ", defaultApp.bucketEnable)

	if len(pre) != 0 {
		if pre == "none" {
			defaultApp.pre = ""
		} else {
			defaultApp.pre = pre
		}
	}

	if len(postOK) != 0 {
		if pre == "none" {
			defaultApp.postOK = ""
		} else {
			defaultApp.postOK = postOK
		}
	}

	if len(postNOK) != 0 {
		if postNOK == "none" {
			defaultApp.postNOK = ""
		} else {
			defaultApp.postNOK = postOK
		}
	}

	if len(httpApi) != 0 {
		if postNOK == "none" {
			errStr := "无效httpapi脚本名称"
			klog.Errorln(errStr)
			panic(errStr)
		} else {
			defaultApp.httpApi = httpApi
		}
	}

	router := gin.Default()
	if defaultApp.bucketEnable {
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

	if port > 0 {
		router.Run(fmt.Sprintf("%s:%d", host, port))
	} else {
		router.Run(host)
	}
}

func apiGateway(stack *hub.Stack, params map[string]string) (interface{}, int) {
	apiGatewayRun(params["host"], params["port"], params["bucket"],
		params["pre"], params["postOK"], params["postNOK"], params["httpApi"])
	return nil, 200
}

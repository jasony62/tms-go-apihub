package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jasony62/tms-go-apihub/api"
	"github.com/jasony62/tms-go-apihub/flow"
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/unit"
	"github.com/joho/godotenv"
)

// 应用的基本信息
type App struct {
	host           string
	port           int
	bucketEnable   bool
	apiDefPath     string
	privateDefPath string
	flowDefPath    string
}

// 1次请求的上下文
func (app App) newStack(c *gin.Context) *hub.Stack {
	// 收到的数据
	inReqData := new(interface{})
	c.ShouldBindJSON(&inReqData)

	stack := new(hub.Stack)
	stack.BucketEnable = app.bucketEnable
	stack.ApiDefPath = app.apiDefPath
	stack.PrivateDefPath = app.privateDefPath
	stack.FlowDefPath = app.flowDefPath
	stack.GinContext = c
	stack.RequestBody = inReqData

	return stack
}

var app = new(App)

// 执行1个API调用
func doRelay(c *gin.Context) {
	// 构造运行上下文
	stack := app.newStack(c)

	apiId := c.Param(`apiId`)
	var bucket string
	if app.bucketEnable {
		bucket = c.Param(`bucket`)
	}
	apiDef, _ := unit.FindApiDef(stack, bucket, apiId)

	// 收到的数据
	inReqData := new(interface{})
	c.BindJSON(&inReqData)

	stack.ApiDef = apiDef

	// 调用api
	result, status := api.Relay(stack, "")

	c.IndentedJSON(status, result)
}

// 执行一个调用流程
func runFlow(c *gin.Context) {
	// 构造运行上下文
	stack := app.newStack(c)

	flowId := c.Param(`flowId`)
	var bucket string
	if app.bucketEnable {
		bucket = c.Param(`bucket`)
	}
	flowDef, _ := unit.FindFlowDef(stack, bucket, flowId)

	stack.FlowDef = flowDef
	stack.StepResult = make(map[string]interface{})

	// 执行编排
	result, status := flow.Run(stack)

	c.IndentedJSON(status, result)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 命令行指定的环境变量文件
var envfile string

func init() {
	flag.StringVar(&envfile, "env", "", "指定环境变量文件")
}

func main() {
	flag.Parse()

	if envfile != "" {
		err := godotenv.Load(envfile)
		if err != nil {
			log.Fatal(err)
		}
	}

	host := os.Getenv("TGAH_HOST")
	if host == "" {
		app.host = "0.0.0.0"
	} else {
		app.host = host
	}
	log.Println("host: ", app.host)

	port := os.Getenv("TGAH_PORT")
	if port == "" {
		app.port = 8080
	} else {
		app.port, _ = strconv.Atoi(port)
	}
	log.Println("port ", app.port)

	bucketEnable := os.Getenv("TGAH_BUCKET_ENABLE")
	re := regexp.MustCompile(`(?i)yes|true`)
	app.bucketEnable = re.MatchString(bucketEnable)
	log.Println("bucket enable ", app.bucketEnable)

	app.apiDefPath = os.Getenv("TGAH_API_DEF_PATH")
	if app.apiDefPath == "" {
		log.Println("没有通过环境变量[TGAH_API_DEF_PATH]指定API定义文件存放位置")
	} else {
		if ok, _ := pathExists(app.apiDefPath); ok {
			log.Println("API定义文件存放位置 ", app.apiDefPath)
		} else {
			log.Printf("通过环境变量[TGAH_API_DEF_PATH]指定的API定义文件存放位置[%s]不存在\n", app.apiDefPath)
			app.apiDefPath = ""
		}
	}
	if app.apiDefPath == "" {
		app.apiDefPath = "./conf/apis"
		log.Println("使用默认API定义文件存放位置 ", app.apiDefPath)
	}

	app.privateDefPath = os.Getenv("TGAH_PRIVATE_DEF_PATH")
	if app.privateDefPath == "" {
		log.Println("没有通过环境变量[TGAH_PRIVATE_DEF_PATH]指定API定义文件存放位置")
	} else {
		if ok, _ := pathExists(app.privateDefPath); ok {
			log.Println("PRIVATE定义文件存放位置 ", app.privateDefPath)
		} else {
			log.Printf("通过环境变量[TGAH_PRIVATE_DEF_PATH]指定的API定义文件存放位置[%s]不存在\n", app.privateDefPath)
			app.privateDefPath = ""
		}
	}
	if app.privateDefPath == "" {
		app.privateDefPath = "./conf/privates"
		log.Println("使用默认PRIVATE定义文件存放位置 ", app.privateDefPath)
	}

	app.flowDefPath = os.Getenv("TGAH_FLOW_DEF_PATH")
	if app.flowDefPath == "" {
		log.Println("没有通过环境变量[TGAH_FLOW_DEF_PATH]指定FLOW定义文件存放位置")
	} else {
		if ok, _ := pathExists(app.flowDefPath); ok {
			log.Println("FLOW定义文件存放位置 ", app.flowDefPath)
		} else {
			log.Printf("通过环境变量[TGAH_FLOW_DEF_PATH]指定FLOW定义文件存放位置[%s]不存在\n", app.flowDefPath)
			app.flowDefPath = ""
		}
	}
	if app.flowDefPath == "" {
		app.flowDefPath = "./conf/flows"
		log.Println("使用默认PRIVATE定义文件存放位置 ", app.privateDefPath)
	}

	router := gin.Default()
	if app.bucketEnable {
		router.Any("/api/:bucket/:apiId", doRelay)
		router.Any("/flow:bucket/:flowId", runFlow)
	} else {
		router.Any("/api/:apiId", doRelay)
		router.Any("/flow/:flowId", runFlow)
	}

	if app.port > 0 {
		router.Run(fmt.Sprintf("%s:%d", app.host, app.port))
	} else {
		router.Run(app.host)
	}
}

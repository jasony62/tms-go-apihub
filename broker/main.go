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
	"github.com/jasony62/tms-go-apihub/schedule"
	"github.com/jasony62/tms-go-apihub/unit"
	"github.com/joho/godotenv"
)

// 1次请求的上下文
func newStack(c *gin.Context) *hub.Stack {
	// 收到的数据
	inReqData := new(interface{})
	c.ShouldBindJSON(&inReqData)

	stack := new(hub.Stack)
	stack.GinContext = c

	stack.StepResult = make(map[string]interface{})
	stack.StepResult["origin"] = *inReqData
	return stack
}

// 执行1个API调用
func runApi(c *gin.Context) {
	// 构造运行上下文
	var err error
	stack := newStack(c)

	stack.ApiDef, err = unit.FindApiDef(stack, c.Param(`apiId`))

	if stack.ApiDef == nil {
		log.Panic("获得API定义失败：", err)
		return
	}

	// 收到的数据
	//inReqData := new(interface{})
	//c.BindJSON(&inReqData)

	// 调用api
	result, status := api.Relay(stack, "")

	c.IndentedJSON(status, result)
}

// 执行一个调用流程
func runFlow(c *gin.Context) {
	// 构造运行上下文
	var err error
	stack := newStack(c)
	stack.FlowDef, err = unit.FindFlowDef(stack, c.Param(`flowId`))

	if stack.FlowDef == nil {
		log.Panic("获得Flow定义失败：", err)
		return
	}

	// 执行编排
	result, status := flow.Run(stack)

	c.IndentedJSON(status, result)
}

// 执行一个计划流程
func runSchedule(c *gin.Context) {
	var err error
	// 构造运行上下文
	stack := newStack(c)

	stack.ScheduleDef, err = unit.FindScheduleDef(stack, c.Param(`scheduleId`))
	if stack.ScheduleDef == nil {
		log.Panic("获得Schedule定义失败：", err)
		return
	}

	// 执行编排
	result, status := schedule.Run(stack)

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

func loadPath(env string, inDefault string) string {
	result := os.Getenv(env)
	if result == "" {
		log.Println("没有通过环境变量", env, "指定API定义文件存放位置")
	} else {
		if ok, _ := pathExists(result); ok {
			log.Println("API定义文件存放位置 ", result)
		} else {
			log.Printf("通过环境变量[TGAH_API_DEF_PATH]指定的API定义文件存放位置[%s]不存在\n", result)
			result = ""
		}
	}
	if result == "" {
		result = inDefault
		log.Println("使用默认API定义文件存放位置 ", result)
	}
	return result
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
		hub.DefaultApp.Host = "0.0.0.0"
	} else {
		hub.DefaultApp.Host = host
	}
	log.Println("host: ", hub.DefaultApp.Host)

	port := os.Getenv("TGAH_PORT")
	if port == "" {
		hub.DefaultApp.Port = 8080
	} else {
		hub.DefaultApp.Port, _ = strconv.Atoi(port)
	}
	log.Println("port ", hub.DefaultApp.Port)

	BucketEnable := os.Getenv("TGAH_BUCKET_ENABLE")
	re := regexp.MustCompile(`(?i)yes|true`)
	hub.DefaultApp.BucketEnable = re.MatchString(BucketEnable)
	log.Println("bucket enable ", hub.DefaultApp.BucketEnable)

	hub.DefaultApp.ApiDefPath = loadPath("TGAH_API_DEF_PATH", "./conf/apis")
	hub.DefaultApp.PrivateDefPath = loadPath("TGAH_PRIVATE_DEF_PATH", "./conf/privates")
	hub.DefaultApp.FlowDefPath = loadPath("TGAH_FLOW_DEF_PATH", "./conf/flows")
	hub.DefaultApp.ScheduleDefPath = loadPath("TGAH_SCHEDULE_DEF_PATH", "./conf/schedules")

	router := gin.Default()
	if hub.DefaultApp.BucketEnable {
		router.Any("/api/:bucket/:apiId", runApi)
		router.Any("/flow:bucket/:flowId", runFlow)
		router.Any("/schedule:bucket/:scheduleId", runSchedule)
	} else {
		router.Any("/api/:apiId", runApi)
		router.Any("/flow/:flowId", runFlow)
		router.Any("/schedule/:scheduleId", runSchedule)
	}

	if hub.DefaultApp.Port > 0 {
		router.Run(fmt.Sprintf("%s:%d", hub.DefaultApp.Host, hub.DefaultApp.Port))
	} else {
		router.Run(hub.DefaultApp.Host)
	}
}

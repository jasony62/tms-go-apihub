package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	klog "k8s.io/klog/v2"

	"github.com/gin-gonic/gin"
	"github.com/jasony62/tms-go-apihub/api"
	"github.com/jasony62/tms-go-apihub/flow"
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/schedule"
	"github.com/jasony62/tms-go-apihub/unit"
	"github.com/jasony62/tms-go-apihub/util"
	"github.com/joho/godotenv"
)

// 1次请求的上下文
func newStack(c *gin.Context) *hub.Stack {
	// 收到的数据
	var value interface{}
	inReqData := new(interface{})
	c.ShouldBindJSON(&inReqData)

	if *inReqData == nil {
		value = make(map[string]interface{})
	} else {
		value = *inReqData
	}

	name := c.Param(`Id`)
	return &hub.Stack{
		GinContext: c,
		StepResult: map[string]interface{}{hub.OriginName: value},
		RootName:   name,
		ChildName:  name,
	}
}

// 执行1个API调用
func callApi(c *gin.Context) {
	// 调用api
	result, status := api.Run(newStack(c))
	c.IndentedJSON(status, result)
}

// 执行一个调用流程
func callFlow(c *gin.Context) {
	// 执行编排
	result, status := flow.Run(newStack(c))
	c.IndentedJSON(status, result)
}

// 执行一个计划流程
func callSchedule(c *gin.Context) {
	// 执行编排
	result, status := schedule.Run(newStack(c))
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
		klog.Infoln("没有通过环境变量", env, "指定API定义文件存放位置")
	} else {
		if ok, _ := pathExists(result); ok {
			klog.Infoln("API定义文件存放位置 ", result)
		} else {
			klog.Infof("通过环境变量[TGAH_API_DEF_PATH]指定的API定义文件存放位置[%s]不存在\n", result)
			result = ""
		}
	}
	if result == "" {
		result = inDefault
		klog.Infoln("使用默认API定义文件存放位置 ", result)
	}
	return result
}

func loadConf() bool {
	downUrl := os.Getenv("TGAH_REMOTE_CONF_DOWNLOAD")

	if downUrl == "1" {
		//从远端下载conf
		confUrl := os.Getenv("TGAH_REMOTE_CONF_URL")
		filename, err := util.DownloadFile(confUrl)
		if err != nil {
			klog.Errorln("Download conf file err: ", err)
			return false
		} else {
			//解压缩
			//filename = os.Getenv("TGAH_REMOTE_CONF_NAME")
			confStoreFolder := os.Getenv("TGAH_REMOTE_CONF_STORE_FOLDER")
			confUnzipPwd := os.Getenv("TGAH_REMOTE_CONF_UNZIP_PWD")
			klog.Infoln("filename: ", filename)
			klog.Infoln("confStoreFolder: ", confStoreFolder)
			klog.Infoln("confUnzipPwd: ", confUnzipPwd)

			err = util.DeCompressZip(filename, confStoreFolder, confUnzipPwd, nil, 0)
			if err != nil {
				klog.Errorln(err)
				return false
			}
		}
	} else {
		klog.Warningln("TGAH_REMOTE_CONF_DOWNLOAD not 1, use local conf!")
		return false
	}

	return true
}

func main() {
	flag.Parse()

	if envfile != "" {
		err := godotenv.Load(envfile)
		if err != nil {
			klog.Fatal(err)
		}
	}

	host := os.Getenv("TGAH_HOST")
	if host == "" {
		hub.DefaultApp.Host = "0.0.0.0"
	} else {
		hub.DefaultApp.Host = host
	}
	klog.Infoln("host: ", hub.DefaultApp.Host)

	port := os.Getenv("TGAH_PORT")
	if port == "" {
		hub.DefaultApp.Port = 8080
	} else {
		hub.DefaultApp.Port, _ = strconv.Atoi(port)
	}
	klog.Infoln("port ", hub.DefaultApp.Port)

	BucketEnable := os.Getenv("TGAH_BUCKET_ENABLE")
	re := regexp.MustCompile(`(?i)yes|true`)
	hub.DefaultApp.BucketEnable = re.MatchString(BucketEnable)
	klog.Infoln("bucket enable ", hub.DefaultApp.BucketEnable)

	if loadConf() {
		klog.Infoln("Download conf zip package from remote url OK")
	}

	unit.LoadConfigJsonData([]string{loadPath("TGAH_API_DEF_PATH", "./conf/apis"),
		loadPath("TGAH_FLOW_DEF_PATH", "./conf/flows"),
		loadPath("TGAH_SCHEDULE_DEF_PATH", "./conf/schedules"),
		loadPath("TGAH_PRIVATE_DEF_PATH", "./conf/privates")})

	unit.LoadConfigPluginData(loadPath("TGAH_PLUGIN_DEF_PATH", "./plugins"))
	router := gin.Default()
	if hub.DefaultApp.BucketEnable {
		router.Any("/api/:bucket/:Id", callApi)
		router.Any("/flow:bucket/:Id", callFlow)
		router.Any("/schedule:bucket/:Id", callSchedule)
	} else {
		router.Any("/api/:Id", callApi)
		router.Any("/flow/:Id", callFlow)
		router.Any("/schedule/:Id", callSchedule)
	}

	if hub.DefaultApp.Port > 0 {
		router.Run(fmt.Sprintf("%s:%d", hub.DefaultApp.Host, hub.DefaultApp.Port))
	} else {
		router.Run(hub.DefaultApp.Host)
	}
}

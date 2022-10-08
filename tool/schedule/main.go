package main

import (
	"encoding/json"
	"flag"
	"os"
	"strings"

	"k8s.io/klog/v2"
)

type ScheduleConf struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	ConcurrentNum int     `json:"concurrentNum"`
	Steps         []Steps `json:"steps,omitempty"`
}

type Steps struct {
	Type string `json:"type"`
	Mode string `json:"mode"`
	Api  Api    `json:"api,omitempty"`
}

type Api struct {
	Name    string `json:"name"`
	Command string `json:"command"`
	Args    []Args `json:"args,omitempty"`
}

type Args struct {
	Name  string `json:"name"`
	Value Value  `json:"value,omitempty"`
}

type Value struct {
	From    string `json:"from"`
	Content string `json:"content"`
}

var httpapisPath string
var schedulePath string

var scheduleConf ScheduleConf

// 初始化
func init() {
	flag.StringVar(&httpapisPath, "from", "../../example/httpapis/__INTERNAL/", "指定httpapisPath文件路径")
	flag.StringVar(&schedulePath, "to", "./schedule/", "指定生成schedule文件路径")
}

func main() {
	initConf()
	readFileName(httpapisPath)
	generateSchedule(schedulePath)
}

func initConf() {
	scheduleConf.Name = "healthCheck"
	scheduleConf.Description = "healthCheck"
	scheduleConf.ConcurrentNum = 100
}

func readFileName(path string) error {
	fileInfoList, err := os.ReadDir(path)
	if err != nil {
		klog.Errorln(err)
		return err
	}
	var prefix string
	j := 0
	// 遍历文件
	for i := range fileInfoList {
		// 是否是一个子目录 若是子目录，进入子目录遍历文件
		if fileInfoList[i].IsDir() {
			klog.Infoln("__schedule 子目录名: ", fileInfoList[i].Name())
			prefix = fileInfoList[i].Name()
			readFileName(path + "/" + prefix)
		} else {
			// 判断文件是否curl类型
			if !strings.HasSuffix(fileInfoList[i].Name(), ".json") {
				continue
			}
			klog.Infoln("__Load schedule 文件: ", fileInfoList[i].Name())
			tempName := fileInfoList[i].Name()
			tempName = strings.Replace(tempName, ".json", "", -1)
			// fileInfoList[i].Name() = string.Replace(fileInfoList[i].Name(), ".json", "", -1)

			steps := Steps{Type: "api", Mode: "concurrent", Api: Api{Name: tempName, Command: "httpApi"}}
			scheduleConf.Steps = append(scheduleConf.Steps, steps)
			args := Args{Name: "name", Value: Value{From: "literal", Content: tempName}}
			scheduleConf.Steps[j].Api.Args = append(scheduleConf.Steps[j].Api.Args, args)
			j++
		}
	}
	return nil
}

func generateSchedule(schedulePath string) {

	fileName := schedulePath + "healthCheck.json"

	byteHttpApi, err := json.Marshal(scheduleConf)
	if err != nil {
		klog.Errorln("json.Marshal失败!", fileName)
		return
	}

	// ！！！os.Create无法自动创建文件路径中不存在的文件夹
	f, err := os.Create(fileName)
	if err != nil {
		klog.Errorln("创建文件失败!", fileName)
	} else {
		defer f.Close()
		_, err = f.Write(byteHttpApi)
		if err != nil {
			klog.Errorln("写入文件失败!", fileName)
		}
	}
}

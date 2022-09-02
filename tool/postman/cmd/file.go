package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/rbretecher/go-postman-collection"
	"k8s.io/klog/v2"
)

// 生成json文件，无法自动创建文件路径中不存在的文件夹
func generateApiHubJson(postmanBytes *postman.Collection, multipleName string) {
	if postmanBytes == nil {
		return
	}
	fileName := ""
	postmanBytes.Info.Name = strings.Replace(postmanBytes.Info.Name, "/", "_", -1)
	apiHubHttpConf.ID = strings.Replace(apiHubHttpConf.ID, "/", "_", -1)
	if multipleName == "" {
		fileName = apiHubJsonPath + postmanBytes.Info.Name + "_" + apiHubHttpConf.ID + ".json"
	} else {
		fileName = apiHubJsonPath + postmanBytes.Info.Name + "_" + multipleName + "_" + apiHubHttpConf.ID + ".json"
	}
	byteHttpApi, err := json.Marshal(apiHubHttpConf)
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

func generateApiHubPrivatesJson(postmanBytes *postman.Collection, privateName string) {

	fileName := apiHubPrivatesJsonPath + privateName + ".json"

	byteHttpApi, err := json.Marshal(apiHubHttpPrivates)
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

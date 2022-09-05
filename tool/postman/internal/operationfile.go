package postmaninternal

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rbretecher/go-postman-collection"
	"k8s.io/klog/v2"
)

// Postman文件转换函数
func convertPostmanFiles(path string) {
	// 读取指定目录下文件信息list
	fileInfoList, err := os.ReadDir(path)
	if err != nil {
		klog.Errorln(err)
		return
	}

	var prefix string
	// 遍历postman_collection文件
	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())
		// 是否是一个子目录 若是子目录，进入进入子目录遍历postman文件
		if fileInfoList[i].IsDir() {
			klog.Infoln("__postman_collections子目录名: ", fileName)
			prefix = fileInfoList[i].Name()
			convertPostmanFiles(path + "/" + prefix)
		} else {
			// 判断文件是否postman_collection类型
			if !strings.HasSuffix(fileName, ".postman_collection") && !strings.HasSuffix(fileName, ".json") {
				continue
			}
			klog.Infoln("__Load postman_collection文件: ", fileName)
			// 读取文件内容到fileBytes
			file, err := os.Open(fileName)
			if err != nil {
				klog.Errorln(err)
				panic(err)
			}
			defer file.Close()
			// Parse the contents
			postmanfileBytes, err := postman.ParseCollection(file)
			if err != nil {
				klog.Errorln(err)
				panic(err)
			}
			getPostmanFilesBytes(postmanfileBytes)
		}
	}
}

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

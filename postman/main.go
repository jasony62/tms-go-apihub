package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"k8s.io/klog/v2"
)

type ApiHubHttpConf struct {
	ID                 string `json:"id"`
	Description        string `json:"description"`
	URL                string `json:"url"`
	Method             string `json:"method"`
	Private            string `json:"private"`
	Requestcontenttype string `json:"requestContentType"`
	Args               []Args `json:"args"`
}

type Value struct {
	From    string `json:"from"`
	Content string `json:"content"`
}
type Args struct {
	In    string `json:"in"`
	Name  string `json:"name"`
	Value Value  `json:"value"`
}

var postmanPath string

// var postmanCollectionBytes postman

func init() {
	flag.StringVar(&postmanPath, "from", "./postman_collections/", "指定postman_collections文件路径")
	// flag.StringVar(&apiHubConfPath, "to", "./to/", "指定转换后的apiHubConf json文件路径")
}

func main() {
	convertPostmanFiles(postmanPath)
}

func convertPostmanFiles(path string) {

	// 读取指定目录下文件信息list
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		klog.Errorln(err)
		return
	}

	var prefix string
	// 遍历postman_collection文件
	for i := range fileInfoList {

		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())
		// 是否是一个子目录
		if fileInfoList[i].IsDir() {
			klog.Infoln("postman_collections子目录fileName: ", fileName)
			prefix = fileInfoList[i].Name()
			convertPostmanFiles(path + "/" + prefix)
		} else {
			// 判断文件是否postman_collection类型
			if !strings.HasSuffix(fileName, ".postman_collection") && !strings.HasSuffix(fileName, ".postman_collection.json") {
				continue
			}
			klog.Infoln("######加载postman_collection(*.postman_collection or .postman_collection.json)文件: ", fileName)
			// ReadFil读取文件内容到fileBytes
			// fileBytes, err := ioutil.ReadFile(fileName)
			// if err != nil {
			// 	klog.Errorln(err)
			// 	panic(err)
			// }
			// covertPostmanToApihubConf(fileInfoList[i])
		}
	}
}

// func covertPostmanToApihubConf(fileInfoList []fs.FileInfo) {
// 	var err error
// 	postmanCollectionBytes, err = postman.ParseCollection(fileInfoList)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	// getServerUrl(postmanCollectionBytes)
// 	// parsePaths(postmanCollectionBytes)
// 	// parseRequestBodies(postmanCollectionBytes.Components.RequestBodies) //?
// 	// generateApiHubConf(postmanCollectionBytes)
// }

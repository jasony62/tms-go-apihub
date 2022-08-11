package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"os"
	"strings"

	"github.com/rbretecher/go-postman-collection"
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

// postman文件路径
var postmanPath string

// json结构体
var apiHubHttpConf ApiHubHttpConf

// 初始化
func init() {
	flag.StringVar(&postmanPath, "from", "./postman_collections/", "指定postman_collections文件路径")
	// flag.StringVar(&apiHubConfPath, "to", "./to/", "指定转换后的apiHubConf json文件路径")
}

/*********************main主程序*****************************/
func main() {
	convertPostmanFiles(postmanPath)
}

// Postman文件转换函数
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
			// _ = postmanfileBytes

			for j := range postmanfileBytes.Items {
				covertOneRequest(postmanfileBytes.Items[j])
			}
		}
	}
}

func covertOneRequest(postmanItem *postman.Items) {

	// 	// _ = postmanfileBytes
	getHttpapiInfo(postmanItem)
	getHttpapiArgs(postmanItem.Request.URL)
}

func getHttpapiInfo(postmanItem *postman.Items) {

	apiHubHttpConf.ID = postmanItem.Name
	klog.Infoln("__request Name : ", apiHubHttpConf.Description)

	apiHubHttpConf.Description = postmanItem.Name
	klog.Infoln("__request Description : ", apiHubHttpConf.Description)

	apiHubHttpConf.URL = postmanItem.Request.URL.Raw
	// getdel := fmt.Sprintf("?")
	// for i := range postmanItem.Responses.URL.Path {
	// getdel := strings.Join(postmanItem.Request.URL.Path)
	// }
	// apiHubHttpConf.URL = strings.TrimSuffix(apiHubHttpConf.URL, getdel)
	klog.Infoln("__request URL : ", apiHubHttpConf.URL)

	apiHubHttpConf.Method = string(postmanItem.Request.Method)
	klog.Infoln("__request Method : ", apiHubHttpConf.Method)

	// parseEvent(postmanItem)
	// apiHubHttpConf.Requestcontenttype = "none" // default content
	// apiHubHttpConf.Private = "none" // default private content
}

func getHttpapiArgs(postmanURL *postman.URL) {
	// for i := range list(postmanURL.Query) {
	if postmanURL.Query == "query" {
		args := Args{In: "query", Name: "", Value: Value{From: "query", Content: ""}}
		apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
	}
	// }
}

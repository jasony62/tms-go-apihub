package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/rbretecher/go-postman-collection"
	"k8s.io/klog/v2"
)

// Postman文件转换函数
func ConvertPostmanFiles(path string) {
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
			ConvertPostmanFiles(path + "/" + prefix)
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

func getPostmanFilesBytes(postmanfileBytes *postman.Collection) {
	if postmanfileBytes != nil {
		for i := range postmanfileBytes.Items {
			if postmanfileBytes.Items[i].Items == nil {
				converOneRequest(postmanfileBytes.Items[i])
				if len(apiHubHttpPrivates.Privates) != 0 {
					apiHubHttpConf.Private = postmanfileBytes.Info.Name + "_" + postmanfileBytes.Items[i].Name + "_key"
					generateApiHubPrivatesJson(postmanfileBytes, apiHubHttpConf.Private)
				}
				generateApiHubJson(postmanfileBytes, "")
			} else {
				for j := range postmanfileBytes.Items[i].Items {
					converOneRequest(postmanfileBytes.Items[i].Items[j])
					if len(apiHubHttpPrivates.Privates) != 0 {
						apiHubHttpConf.Private = postmanfileBytes.Info.Name + "_" + postmanfileBytes.Items[i].Name + "_" + postmanfileBytes.Items[i].Items[j].Name + "_key"
						generateApiHubPrivatesJson(postmanfileBytes, apiHubHttpConf.Private)
					}
					generateApiHubJson(postmanfileBytes, postmanfileBytes.Items[i].Name)
				}
			}
		}
	}
}

// 转换postman collection中一个request
func converOneRequest(postmanItem *postman.Items) {
	if postmanItem == nil {
		return
	}
	httpapiArgsLen := len(apiHubHttpConf.Args)
	delHttpapiConfArgs(httpapiArgsLen)
	httpapiPrivatesLen := len(apiHubHttpPrivates.Privates)
	delHttpapiPrivates(httpapiPrivatesLen)
	coversionFuncMap = make(map[string]string)
	preGlobalKeyMap = make(map[string]string)
	preGlobalValueMap = make(map[string]string)
	keyWordGlobal := "postman.setEnvironmentVariable"

	getHttpapiInfo(postmanItem)
	getPostmanEventFunc(postmanItem, preEventFuncKeyMap, keyWordGlobal)
	getHttpapiArgs(postmanItem.Request)

}

// 获取Args
func getHttpapiArgs(postmanRequest *postman.Request) {
	if postmanRequest == nil {
		return
	}
	// 解析header
	if postmanRequest.Header != nil {
		for i := range postmanRequest.Header {
			if postmanRequest.Header[i].Key != "Content-Type" {
				args := Args{In: "header", Name: postmanRequest.Header[i].Key, Value: Value{From: "private", Content: postmanRequest.Header[i].Key}}
				apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
				privates := Privates{Name: postmanRequest.Header[i].Key, Value: postmanRequest.Header[i].Value}
				apiHubHttpPrivates.Privates = append(apiHubHttpPrivates.Privates, privates)
			} else if postmanRequest.Header[i].Key == "Content-Type" {
				headerindex := strings.Index(postmanRequest.Header[i].Value, "/")
				apiHubHttpConf.Requestcontenttype = postmanRequest.Header[i].Value[headerindex+1:]
			}
		}
	}

	switch apiHubHttpConf.Method {
	case "GET":
		// 解析qury
		if postmanRequest.URL.Query != nil {
			parseRequestUrlQuery(postmanRequest.URL.Query)
		}
	case "POST":
		// 解析body
		if postmanRequest.Body != nil {
			switch postmanRequest.Body.Mode {
			case "urlencoded":
				parseRequestBodyUrlencoded(postmanRequest.Body)
			case "raw":
				parseRequestBodyRaw(postmanRequest.Body)
			default:
			}
		}
	}
}

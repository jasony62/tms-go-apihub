// /////////////////////////////////////////////////////////////
//
//	                                                          //
//	                    _ooOoo_                               //
//	                   o8888888o                              //
//	                   88" . "88                              //
//	                   (| ^_^ |)                              //
//	                   O\  =  /O                              //
//	                ____/`---'\____                           //
//	              .'  \\|     |//  `.                         //
//	             /  \\|||  :  |||//  \                        //
//	            /  _||||| -:- |||||-  \                       //
//	            |   | \\\  -  /// |   |                       //
//	            | \_|  ''\---/''  |   |                       //
//	            \  .-\__  `-`  ___/-. /                       //
//	          ___`. .'  /--.--\  `. . ___                     //
//	        ."" '<  `.___\_<|>_/___.'  >'"".                  //
//	      | | :  `- \`.;`\ _ /`;.`/ - ` : | |                 //
//	      \  \ `-.   \_ __\ /__ _/   .-` /  /                 //
//	========`-.____`-.___\_____/___.-`____.-'========         //
//	                     `=---='                              //
//	^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^        //
//	         佛祖保佑       永不宕机     永无BUG                //
//
// /////////////////////////////////////////////////////////////
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

type Args struct {
	In    string `json:"in"`
	Name  string `json:"name"`
	Value Value  `json:"value"`
}

type Value struct {
	From    string `json:"from"`
	Content string `json:"content"`
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
	klog.Infoln("__request Name : ", apiHubHttpConf.ID)

	apiHubHttpConf.Description = postmanItem.Name
	klog.Infoln("__request Description : ", apiHubHttpConf.Description)

	apiHubHttpConf.URL = getPostmanURL(postmanItem.Request.URL)
	klog.Infoln("__request URL : ", apiHubHttpConf.URL)

	apiHubHttpConf.Method = string(postmanItem.Request.Method)
	klog.Infoln("__request Method : ", apiHubHttpConf.Method)

	// getPostmanEvent(postmanItem)
	apiHubHttpConf.Requestcontenttype = "none" // default content
	apiHubHttpConf.Private = "none"            // default private content
}

// 获取Request URL
func getPostmanURL(postmanUrl *postman.URL) string {

	httpapiUrl := postmanUrl.Protocol + "://"
	// Host IP
	for i := range postmanUrl.Host {
		if i > 0 {
			httpapiUrl = httpapiUrl + "." + postmanUrl.Host[i]
		} else {
			httpapiUrl = httpapiUrl + postmanUrl.Host[i]
		}
	}
	// Port number
	if postmanUrl.Port != "" {
		httpapiUrl = httpapiUrl + ":" + postmanUrl.Port + "/"
	} else {
		httpapiUrl = httpapiUrl + "/"
	}
	// Path
	for i := range postmanUrl.Path {
		if postmanUrl.Path[i] != "" {
			httpapiUrl = httpapiUrl + postmanUrl.Path[i] + "/"
		}
	}

	return httpapiUrl
}

// 获取Args
// func getHttpapiArgs(postmanURL *postman.URL) {

// 	// postmanURL.Query 是个type interface{}，坑！！！
// 	var list []string
// 	if reflect.TypeOf(postmanURL.Query).Kind() == reflect.Slice {
// 		s := reflect.ValueOf(postmanURL.Query)
// 		for i := 0; i < s.Len(); i++ {
// 			ele := s.Index(i)
// 			list = append(list, ele.Interface().(string))
// 		}
// 	}
// }

// func getHttpapiOneQuery(postmanQuery *postman.Query) string {
// 	if postmanURL.Query == "query" {
// 		args := Args{In: "query", Name: "", Value: Value{From: "query", Content: ""}}
// 		apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
// 	}
// }

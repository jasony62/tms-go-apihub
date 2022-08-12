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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mozillazg/go-pinyin"
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

// 导出json路径
var apiHubJsonPath string

// 初始化
func init() {
	flag.StringVar(&postmanPath, "from", "./postman_collections/", "指定postman_collections文件路径")
	flag.StringVar(&apiHubJsonPath, "to", "./jsonFiles/", "指定转换后的apiHub json文件路径")
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

			for i := range postmanfileBytes.Items {
				converOneRequest(postmanfileBytes.Items[i])
				generateApiHubJson(postmanfileBytes)
			}
		}
	}
}

func converOneRequest(postmanItem *postman.Items) {

	if postmanItem == nil {
		return
	}

	httpapiArgsLen := len(apiHubHttpConf.Args)
	delHttpapiQuery(httpapiArgsLen)

	getHttpapiInfo(postmanItem)
	getHttpapiArgs(postmanItem.Request)
}

// 删除上个request append到args的值
func delHttpapiQuery(httpapiQueryLen int) {
	apiHubHttpConf.Args = append(apiHubHttpConf.Args[:0], apiHubHttpConf.Args[httpapiQueryLen:]...)
}

func getHttpapiInfo(postmanItem *postman.Items) {

	if postmanItem == nil {
		return
	}

	apiHubHttpConf.ID = postmanItem.Name
	klog.Infoln("__request Name : ", apiHubHttpConf.ID)

	apiHubHttpConf.Description = postmanItem.Name
	klog.Infoln("__request Description : ", apiHubHttpConf.Description)

	apiHubHttpConf.URL = getPostmanURL(postmanItem.Request.URL)
	klog.Infoln("__request URL : ", apiHubHttpConf.URL)

	apiHubHttpConf.Method = string(postmanItem.Request.Method)
	klog.Infoln("__request Method : ", apiHubHttpConf.Method)

	// getPostmanEvent(postmanItem)
	apiHubHttpConf.Requestcontenttype = "json" // default content
	apiHubHttpConf.Private = ""                // default private content
}

// 获取Request URL
func getPostmanURL(postmanUrl *postman.URL) string {

	if postmanUrl == nil {
		return ""
	}

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
func getHttpapiArgs(postmanRequest *postman.Request) {

	if postmanRequest == nil {
		return
	}

	if postmanRequest.Header != nil {
		for i := range postmanRequest.Header {
			args := Args{In: "header", Name: postmanRequest.Header[i].Key, Value: Value{From: "header", Content: postmanRequest.Header[i].Key}}
			apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
		}
	}

	// postmanURL.Query 是个type interface{}，坑！！！
	if postmanRequest.URL.Query != nil {
		httpapiQuery := postmanRequest.URL.Query.([]interface{})
		for i := range httpapiQuery {
			httpapiQueryArg := httpapiQuery[i]
			valuename := httpapiQueryArg.(map[string]interface{})["key"]
			// valuecontent := httpapiQueryArg.(map[string]interface{})["value"]
			args := Args{In: "query", Name: valuename.(string), Value: Value{From: "query", Content: valuename.(string)}}
			apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
			// klog.Infoln("__httpapiQueryArgs valuename is : ", valuename.(string))
			// klog.Infoln("__httpapiQueryArgs valuecontent is : ", valuecontent.(string))
		}
	}

	// if postmanRequest.Body != nil {
	// 	for i := range postmanRequest.Body {

	// 	}
	// }
}

func generateApiHubJson(postmanBytes *postman.Collection) {
	if postmanBytes == nil {
		return
	}

	infoName := converCNtoPinyin(postmanBytes.Info.Name)
	infoID := converCNtoPinyin(apiHubHttpConf.ID)
	fileName := apiHubJsonPath + infoName + "_" + infoID + ".json"
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

// name中文转拼音
func converCNtoPinyin(postmanBytesInfoName string) string {
	var infoName []string = []string{}
	infoNameTemp := pinyin.LazyConvert(postmanBytesInfoName, nil)
	for infoNameTemp, v := range infoNameTemp {
		if infoNameTemp == 0 {
			infoName = append(infoName, v)
		} else {
			infoName = append(infoName, ",")
			infoName = append(infoName, v)
		}
	}
	resultInfoName := strings.Trim(fmt.Sprint(infoName), "[]")
	result2InfoName := strings.Replace(resultInfoName, " , ", "", -1)
	if result2InfoName == "" {
		result2InfoName = postmanBytesInfoName
	}
	return result2InfoName
}

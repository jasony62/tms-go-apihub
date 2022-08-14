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

// 创建list,提前预设
// 左：postman js关键字
// 右：apihub内部函数名称
var preEventFuncReferenceList = map[string]string{
	"getTime":      "utc",
	"CryptoJS.MD5": "md5",
}

// postman js关键字
var preEventFuncKeyList = []string{
	"getTime",
	"CryptoJS.MD5",
}

// 相当于postman脚本中全局变量转换的一个中间量，映射postman脚本requset中全局变量值到apihub内部函数名称
// coversionFuncMap[time] = preEventFuncReferenceList["getTime"] = "utc"
var coversionFuncMap map[string]string

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
	coversionFuncMap = make(map[string]string)
	getPostmanEventFunc(postmanItem, preEventFuncKeyList)
	getHttpapiArgs(postmanItem.Request)

}

// 删除上个request append到args的值
func delHttpapiQuery(httpapiArgsLen int) {
	apiHubHttpConf.Args = append(apiHubHttpConf.Args[:0], apiHubHttpConf.Args[httpapiArgsLen:]...)
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

	apiHubHttpConf.Private = "" // default private content
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
	// 解析header
	if postmanRequest.Header != nil {
		for i := range postmanRequest.Header {
			args := Args{In: "header", Name: postmanRequest.Header[i].Key, Value: Value{From: "header", Content: postmanRequest.Header[i].Key}}
			apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
		}
	}
	// 解析query
	if postmanRequest.URL.Query != nil {
		parseRequestUrlQuery(postmanRequest.URL.Query)
	}
	// 解析body中的header、func
	if postmanRequest.Body != nil {
		parseRequestBody(postmanRequest.Body)
	}
}

func parseRequestUrlQuery(postmanRequestURLQuery interface{}) {
	if postmanRequestURLQuery != nil { // postmanURL.Query 是个type interface{}，坑！！！
		httpapiQuery := postmanRequestURLQuery.([]interface{})
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
}

func parseRequestBody(postmanRequestBody *postman.Body) {
	if postmanRequestBody != nil {
		requestBody := postmanRequestBody
		switch requestBody.Mode {
		case "raw":
			apiHubHttpConf.Requestcontenttype = "jsonraw"
		case "x-www-form-urlencoded":
			apiHubHttpConf.Requestcontenttype = "form"
		case "application/json":
			apiHubHttpConf.Requestcontenttype = "json"
		default:
			apiHubHttpConf.Requestcontenttype = requestBody.Mode
		}
		if requestBody.Raw != "" {
			// klog.Infoln("__httpapirequestBody.Raw is : ", requestBody.Raw)
			nameString := ""
			contentString := ""
			backNameIndex := 0
			backContentIndex := 0
			for i := 0; i < len(requestBody.Raw); i++ {
				nameString, backNameIndex = getStringBetweenDoubleQuotationMarks(requestBody.Raw[i:])
				// klog.Infoln("test request raw is :", requestBody.Raw[backNameIndex+i+3:backNameIndex+i+4])
				if requestBody.Raw[backNameIndex+i+3:backNameIndex+i+5] == "{{" {
					contentString, backContentIndex = getStringBetweenDoubleBrackets(requestBody.Raw[backNameIndex+i+2:])
					args := Args{In: "header", Name: nameString, Value: Value{From: "func", Content: coversionFuncMap[contentString]}}
					apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
				} else if requestBody.Raw[backNameIndex+i+3:backNameIndex+i+4] == "\"" {
					contentString, backContentIndex = getStringBetweenDoubleQuotationMarks(requestBody.Raw[backNameIndex+i+2:])
					args := Args{In: "header", Name: nameString, Value: Value{From: "literal", Content: contentString}}
					apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
				} else {
					nameString = ""
					contentString = ""
					klog.Infoln("__parseRequestBodyRawError:Format error")
					break
				}
				i = backNameIndex + backContentIndex + i + 8
			}
		}
	} else {
		apiHubHttpConf.Requestcontenttype = ""
	}
}

func getPostmanEventFunc(postmanItem *postman.Items, preFuncKeyWord []string) {
	if (postmanItem == nil) || (preFuncKeyWord == nil) {
		return
	}
	for i := range postmanItem.Events { // 通常Events就一个
		if postmanItem.Events[i].Script.Type == "text/javascript" {
			for j := range postmanItem.Events[i].Script.Exec {
				for k := range preFuncKeyWord {
					keyWordIndex := strings.Index(postmanItem.Events[i].Script.Exec[j], preFuncKeyWord[k])
					if keyWordIndex != -1 {
						switch preFuncKeyWord[k] {
						case "getTime":
							keyWordString, index := getStringBetweenDoubleQuotationMarks(postmanItem.Events[i].Script.Exec[j])
							// 解析js代码中EnvironmentVariable，查找到getTime代码，确定本行生成time变量名称，赋值到coversionFuncMap
							keyWordString = strings.TrimSpace(keyWordString)
							if index != -1 && keyWordString != "" {
								coversionFuncMap[keyWordString] = preEventFuncReferenceList[preFuncKeyWord[k]]
							}
						case "CryptoJS.MD5":
							keyWordString, index := getStringBetweenSpecifySymbols(postmanItem.Events[i].Script.Exec[j], "var", "=")
							keyWordString = strings.TrimSpace(keyWordString)
							if index != -1 && keyWordString != "" {
								coversionFuncMap[keyWordString] = preEventFuncReferenceList[preFuncKeyWord[k]]
							}
							// klog.Infoln("__postmanItem.Events[i].Script.Exec[j]", keyWordString)
						default:
						}
					}
				}
			}
		} else {
			klog.Infoln("__postmanItem.Events[i].Script.Type not text/javascript")
			return
		}
	}
}

func getStringBetweenDoubleQuotationMarks(inputStrings string) (outputString string, outputIndex int) {
	return getStringBetweenSpecifySymbols(inputStrings, "\"", "\"")
}
func getStringBetweenDoubleBrackets(inputStrings string) (outputString string, outputIndex int) {
	return getStringBetweenSpecifySymbols(inputStrings, "{{", "}}")
}
func getStringBetweenSpecifySymbols(inputStrings string, specifySymbolBefore string, specifySymbolAfter string) (outputString string, outputIndex int) {
	currentIndex := strings.Index(inputStrings, specifySymbolBefore)
	if currentIndex != -1 {
		nextIndex := strings.Index(inputStrings[currentIndex+len(specifySymbolBefore):], specifySymbolAfter)
		if nextIndex != -1 {
			outputString = inputStrings[currentIndex+len(specifySymbolBefore) : currentIndex+len(specifySymbolBefore)+nextIndex]
			outputIndex = nextIndex + len(specifySymbolAfter) + currentIndex
		} else {
			outputString = ""
			outputIndex = -1
		}
	} else {
		outputString = ""
		outputIndex = -1
	}
	return outputString, outputIndex
}

// 生成json文件
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

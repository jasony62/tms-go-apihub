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
	Value Value  `json:"value,omitempty"`
}

type Value struct {
	From    string            `json:"from"`
	Content string            `json:"content",omitempty"`
	Args    string            `json:"args",omitempty"`
	Json    map[string]string `json:"json,omitempty"`
	// Json    *interface{} `json:"json,omitempty"`
}

// 创建list,提前预设
// 左：postman js关键字
// 右：apihub内部函数名称
var preEventFuncReferenceList = map[string]string{
	"getTime":      "utc",
	"CryptoJS.MD5": "md5",
}

// 创建list,提前预设
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

			for i := range postmanfileBytes.Items {
				if postmanfileBytes.Items[i].Items == nil {
					converOneRequest(postmanfileBytes.Items[i])
					generateApiHubJson(postmanfileBytes)
				} else {
					for j := range postmanfileBytes.Items[i].Items {
						converOneRequest(postmanfileBytes.Items[i].Items[j])
						generateApiHubJson(postmanfileBytes)
					}
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

	getHttpapiInfo(postmanItem)
	coversionFuncMap = make(map[string]string)
	getPostmanEventFunc(postmanItem, preEventFuncKeyList)
	getHttpapiArgs(postmanItem.Request)

}

// 删除上个request append到args的值
func delHttpapiConfArgs(httpapiArgsLen int) {
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

	apiHubHttpConf.Private = ""                // default private content
	apiHubHttpConf.Requestcontenttype = "json" //default content typeStr
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
			if postmanRequest.Header[i].Key != "Content-Type" {
				args := Args{In: "header", Name: postmanRequest.Header[i].Key, Value: Value{From: "header", Content: postmanRequest.Header[i].Value}}
				apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
			} else if postmanRequest.Header[i].Key == "Content-Type" {
				headerindex := strings.Index(postmanRequest.Header[i].Value, "/")
				apiHubHttpConf.Requestcontenttype = postmanRequest.Header[i].Value[headerindex+1:]
			}
		}
	}
	// 解析query
	if postmanRequest.URL.Query != nil {
		parseRequestUrlQuery(postmanRequest.URL.Query)
	}
	// 解析body
	if postmanRequest.Body != nil {
		parseRequestBody(postmanRequest.Body)
	}
}

// 解析query
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

// 解析body
func parseRequestBody(postmanRequestBody *postman.Body) {
	if postmanRequestBody != nil {
		requestBody := postmanRequestBody
		// switch requestBody.Mode {
		// case "raw":
		// case "x-www-form-urlencoded":
		// case "application/json":
		// default:
		// }
		if requestBody.Raw != "" {
			klog.Infoln("__httpapirequestBody.Raw is : ", requestBody.Raw)
			nameString := ""
			contentString := ""
			backNameIndex := 0
			backContentIndex := 0
			for i := 0; i < len(requestBody.Raw); i++ {
				nameString, backNameIndex = getStringBetweenDoubleQuotationMarks(requestBody.Raw[i:])
				// klog.Infoln("test request raw is :", requestBody.Raw[backNameIndex+i+3:backNameIndex+i+4])
				if backNameIndex != -1 {
					tempstring1, tempstringflag1 := getStringBetweenDoubleQuotationMarks(requestBody.Raw[backNameIndex+i:])
					tempstring2, tempstringflag2 := getStringBetweenSpecifySymbols(requestBody.Raw[backNameIndex+i:], "\"", "{{")
					tempstring1 = strings.TrimSpace(tempstring1)
					tempstring2 = strings.TrimSpace(tempstring2)

					if (tempstringflag1 != -1) && (tempstring1 == ":") {
						contentString, backContentIndex = getStringBetweenDoubleQuotationMarks(requestBody.Raw[backNameIndex+i+tempstringflag1:])
						if backContentIndex != -1 {
							args := Args{In: "vars", Name: nameString, Value: Value{From: "literal", Content: contentString}}
							apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
							tempbodyjson := Value.Json{
								Type: "json",
							}
							_ = tempbodyjson
							args = Args{In: "body", Name: "body", Value: Value{From: "json", Content: contentString}}
						}
						i = backNameIndex + backContentIndex + i + tempstringflag1
					} else if (tempstringflag2 != -1) && (tempstring2 == ":") {
						contentString, backContentIndex = getStringBetweenDoubleBrackets(requestBody.Raw[backNameIndex+i+tempstringflag2:])
						if backContentIndex != -1 {
							// if coversionFuncMap[contentString] == "md5" {
							// 	nameList = strings.TrimSpace(nameList)
							// 	args := Args{In: "header", Name: nameString, Value: Value{From: "func", Content: coversionFuncMap[contentString], Args: nameList}}
							// 	apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
							// } else {
							// 	args := Args{In: "vars", Name: nameString, Value: Value{From: "func", Content: coversionFuncMap[contentString]}}
							// 	apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
							// 	nameList = nameList + " " + nameString
							// }
							_ = contentString
							i = backNameIndex + backContentIndex + i + tempstringflag2
						}
					} else {
						nameString = ""
						contentString = ""
						klog.Infoln("__parseRequestBodyRawError:Format error")
						break
					}
				}
			}

		}
	} else {
		apiHubHttpConf.Requestcontenttype = ""
	}
}

// 获取postman Event中的全局变量和js函数
func getPostmanEventFunc(postmanItem *postman.Items, preFuncKeyWord []string) {
	if (postmanItem == nil) || (preFuncKeyWord == nil) {
		return
	}
	for i := range postmanItem.Events { // 通常Events就一个，但是个数组形式，所以使用遍历
		if postmanItem.Events[i].Script.Type == "text/javascript" {
			for j := range postmanItem.Events[i].Script.Exec { // Exec中js命令一行是一个数组元素
				for k := range preFuncKeyWord { // 查找js命令行中有无js常见命令的关键字
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

// 获取指定字符中间的字符串，并返回字符串最右的索引值，backIndex = -1表示错误
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

// 生成json文件，无法自动创建文件路径中不存在的文件夹
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

// name中文转拼音，解决apihub不识别中文文件名问题。转换仅支持纯中文，若纯英文则默认不修改
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

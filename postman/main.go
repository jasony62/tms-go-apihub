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
	"strconv"
	"strings"

	"github.com/rbretecher/go-postman-collection"
	"k8s.io/klog/v2"
)

type ApiHubHttpConf struct {
	ID                 string `json:"id"`
	Description        string `json:"description"`
	URL                string `json:"url"`
	Method             string `json:"method"`
	Private            string `json:"private,omitempty"`
	Requestcontenttype string `json:"requestContentType"`
	Args               []Args `json:"args,omitempty"`
}

type Args struct {
	In    string `json:"in"`
	Name  string `json:"name"`
	Value Value  `json:"value,omitempty"`
}

type Value struct {
	From    string            `json:"from"`
	Content string            `json:"content,omitempty"`
	Args    string            `json:"args,omitempty"`
	Json    map[string]string `json:"json,omitempty"`
	// Json    *interface{} `json:"json,omitempty"`
}

type ApiHubHttpPrivates struct {
	Privates []Privates `json:"privates"`
}

type Privates struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// 创建list,提前预设
// 左：postman js关键字
// 右：apihub内部函数名称
var preEventFuncReferenceList = map[string]string{
	"getTime":      "utc_ms",
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
var apiHubHttpPrivates ApiHubHttpPrivates

// 导出json路径
var apiHubJsonPath string
var apiHubPrivatesJsonPath string

// Event中MD5涉及的变量内容组成的字符串数组
var setEnvironmentVariableMD5Array []string

// 初始化
func init() {
	flag.StringVar(&postmanPath, "from", "./postman_collections/", "指定postman_collections文件路径")
	flag.StringVar(&apiHubJsonPath, "to", "./httpapis/", "指定转换后的apiHub json文件路径")
	flag.StringVar(&apiHubPrivatesJsonPath, "private", "./privates/", "指定转换后的apiHub privates json文件路径")
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

	getHttpapiInfo(postmanItem)
	coversionFuncMap = make(map[string]string)
	getPostmanEventFunc(postmanItem, preEventFuncKeyList)
	getHttpapiArgs(postmanItem.Request)

}

// 删除上个request append到args的值
func delHttpapiConfArgs(httpapiArgsLen int) {
	apiHubHttpConf.Args = append(apiHubHttpConf.Args[:0], apiHubHttpConf.Args[httpapiArgsLen:]...)
}

func delHttpapiPrivates(httpapiPrivatesLen int) {
	apiHubHttpPrivates.Privates = append(apiHubHttpPrivates.Privates[:0], apiHubHttpPrivates.Privates[httpapiPrivatesLen:]...)
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
			if i != (len(postmanUrl.Path) - 1) {
				httpapiUrl = httpapiUrl + postmanUrl.Path[i] + "/"
			} else {
				httpapiUrl = httpapiUrl + postmanUrl.Path[i]
			}
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
			parseRequestBody(postmanRequest.Body)
		}
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
		if requestBody.Raw != "" {
			klog.Infoln("__httpapirequestBody.Raw is : ", requestBody.Raw)
			nameString := ""
			backNameIndex := 0
			contentString := ""
			backContentIndex := 0
			contentStringFunc := ""
			backContentIndexFunc := 0
			tempbodyjson := make(map[string]string) // 组建body的json对象
			tempRequestRawArray := *new([]string)   // 暂存MD5需要对比的全部body元素
			// tempRequestRawArrayMap := make(map[string]string)
			nameList := ""
			for i := 0; i < len(requestBody.Raw); i++ {
				nameString, backNameIndex = getStringBetweenDoubleQuotationMarks(requestBody.Raw[i:])
				// klog.Infoln("test request raw is :", requestBody.Raw[backNameIndex+i+3:backNameIndex+i+4])
				if backNameIndex != -1 {
					_, tempstringflag := getStringBetweenSpecifySymbols(requestBody.Raw[backNameIndex+i:], "\"", ":")
					_, tempstringflag2 := getStringBetweenSpecifySymbols(requestBody.Raw[backNameIndex+i:], "\"", ",")
					if tempstringflag != -1 {
						// 考虑到字符串最后一组没有 ， 增加判断
						if tempstringflag2 != -1 {
							contentString, backContentIndex = getStringBetweenDoubleQuotationMarks(requestBody.Raw[backNameIndex+i+tempstringflag : backNameIndex+i+tempstringflag2])
						} else {
							contentString, backContentIndex = getStringBetweenDoubleQuotationMarks(requestBody.Raw[backNameIndex+i+tempstringflag:])
						}
						contentStringFunc, backContentIndexFunc = getStringBetweenDoubleBrackets(requestBody.Raw[backNameIndex+i+tempstringflag:])
						// "":""
						if backContentIndex != -1 {
							tempbodyjson[nameString] = contentString
							tempRequestRawArray = append(tempRequestRawArray, contentString)
							/* vars bad Code
							args := Args{In: "vars", Name: nameString, Value: Value{From: "private", Content: nameString}}
							apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
							tempbodyjson[nameString] = "{{" + ".var." + nameString + "}}"
							privates := Privates{Name: nameString, Value: contentString}
							apiHubHttpPrivates.Privates = append(apiHubHttpPrivates.Privates, privates)*/
							i = backNameIndex + backContentIndex + i + tempstringflag
						} else if backContentIndexFunc != -1 { // "":{{}} 全局变量或函数
							tempbodyjson[nameString] = "{{" + ".vars." + contentStringFunc + "}}"
							tempRequestRawArray = append(tempRequestRawArray, contentStringFunc)
							if coversionFuncMap[contentStringFunc] == "md5" {

								for a := range setEnvironmentVariableMD5Array {
									for n, v := range tempbodyjson {
										tempstring, tempindex := getStringBetweenSpecifySymbols(v, "vars.", "}}")
										if tempindex != -1 { //如果是func类型
											if (setEnvironmentVariableMD5Array[a] == tempstring) && (setEnvironmentVariableMD5Array[a] != contentStringFunc) {
												nameList = nameList + " " + tempstring
												break
											}
										} else {
											if (setEnvironmentVariableMD5Array[a] == v) && (v != "") {
												args := Args{In: "vars", Name: n, Value: Value{From: "literal", Content: v}}
												apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
												nameList = nameList + " " + n
												break
											}
										}
									}
								}
								tempLocalMD5Key, _ := Arrcmp(tempRequestRawArray, setEnvironmentVariableMD5Array)
								for x := range tempLocalMD5Key {
									args := Args{In: "vars", Name: ("key" + strconv.Itoa(x)), Value: Value{From: "literal", Content: tempLocalMD5Key[x]}}
									apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
									nameList = nameList + " " + "key" + strconv.Itoa(x)
								}

								nameList = strings.TrimSpace(nameList)
								args := Args{In: "vars", Name: nameString, Value: Value{From: "func", Content: coversionFuncMap[contentStringFunc], Args: nameList}}
								apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
							} else {
								args := Args{In: "vars", Name: contentStringFunc, Value: Value{From: "func", Content: coversionFuncMap[contentStringFunc]}}
								apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
								// nameList = nameList + " " + coversionFuncMap[contentStringFunc]
							}
							// _ = contentString
							i = backNameIndex + backContentIndex + i + tempstringflag
						} else {
							nameString = ""
							contentString = ""
							klog.Infoln("__parseRequestBodyRawError:Format error")
							break
						}
					}
				}
			}
			bodyArgs := Args{In: "body", Name: "body", Value: Value{From: "json", Json: tempbodyjson}}
			apiHubHttpConf.Args = append(apiHubHttpConf.Args, bodyArgs)
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
					if (strings.Index(postmanItem.Events[i].Script.Exec[j], preFuncKeyWord[k])) != -1 {
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
							getSetEnvironmentVariableMD5(postmanItem, postmanItem.Events[i].Script.Exec[j], keyWordString)
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

func getSetEnvironmentVariableMD5(postmanItem *postman.Items, postmanEventScriptExec string, keyWordString string) {
	if postmanItem == nil {
		return
	}

	tempEnvironmentVariable, _ := getStringBetweenSpecifySymbols(postmanEventScriptExec, "(", ")")
	tempEnvironmentVariable = strings.TrimSpace(tempEnvironmentVariable)
	setEnvironmentVariableMD5 := getStringFromEvent(postmanItem, tempEnvironmentVariable)

	i := 0
	setEnvironmentVariableMD5Array = *new([]string)
	if len(setEnvironmentVariableMD5) > 0 {
		for i = 0; i < len(setEnvironmentVariableMD5); i++ {
			if i < 3 {
				tempString, tempIndex := getStringBetweenDoubleQuotationMarks(setEnvironmentVariableMD5[i:])
				if tempIndex != -1 {
					setEnvironmentVariableMD5Array = append(setEnvironmentVariableMD5Array, tempString)
					i = i + tempIndex
				}
			} else {
				tempString, tempIndex := getStringBetweenSpecifySymbols(setEnvironmentVariableMD5[i:], "+", "+")
				tempString2, tempIndex2 := getStringBetweenSpecifySymbols(setEnvironmentVariableMD5[i:], "\"", "\"")
				if tempIndex != -1 {
					tempString1, tempIndex1 := getStringBetweenDoubleQuotationMarks(tempString)
					if tempIndex1 != -1 {
						setEnvironmentVariableMD5Array = append(setEnvironmentVariableMD5Array, tempString1)
						i = i + tempIndex1
						if tempString1 == "" {
							i = i + 4
						}
					} else {
						tempString = strings.TrimSpace(tempString)
						setEnvironmentVariableMD5Array = append(setEnvironmentVariableMD5Array, tempString)
						i = i + len(tempString)
					}
				} else if tempIndex2 != -1 {
					setEnvironmentVariableMD5Array = append(setEnvironmentVariableMD5Array, tempString2)
					i = i + tempIndex2
				}
			}

		}
	}

}

func getStringFromEvent(postmanItem *postman.Items, keyWordString string) string {
	if postmanItem == nil {
		return ""
	}
	for i := range postmanItem.Events { // 通常Events就一个，但是个数组形式，所以使用遍历
		if postmanItem.Events[i].Script.Type == "text/javascript" {
			for j := range postmanItem.Events[i].Script.Exec { // Exec中js命令一行是一个数组元素
				tempIndex := strings.Index(postmanItem.Events[i].Script.Exec[j], keyWordString)
				if tempIndex != -1 {
					backString := postmanItem.Events[i].Script.Exec[j][tempIndex+len(keyWordString):]
					return backString
				}
			}
		}
	}
	return ""
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
func generateApiHubJson(postmanBytes *postman.Collection, multipleName string) {
	if postmanBytes == nil {
		return
	}
	fileName := ""
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

// // name中文转拼音，解决apihub不识别中文文件名问题。转换仅支持纯中文，若纯英文则默认不修改
// func converCNtoPinyin(postmanBytesInfoName string) string {
// 	var infoName []string = []string{}
// 	infoNameTemp := pinyin.LazyConvert(postmanBytesInfoName, nil)
// 	for infoNameTemp, v := range infoNameTemp {
// 		if infoNameTemp == 0 {
// 			infoName = append(infoName, v)
// 		} else {
// 			infoName = append(infoName, ",")
// 			infoName = append(infoName, v)
// 		}
// 	}
// 	resultInfoName := strings.Trim(fmt.Sprint(infoName), "[]")
// 	result2InfoName := strings.Replace(resultInfoName, " , ", "", -1)
// 	if result2InfoName == "" {
// 		result2InfoName = postmanBytesInfoName
// 	}
// 	return result2InfoName
// }

func Arrcmp(src []string, dest []string) ([]string, []string) {
	msrc := make(map[string]byte) //按源数组建索引
	mall := make(map[string]byte) //源+目所有元素建索引

	var set []string //交集

	//1.源数组建立map
	for _, v := range src {
		msrc[v] = 0
		mall[v] = 0
	}
	//2.目数组中，存不进去，即重复元素，所有存不进去的集合就是并集
	for _, v := range dest {
		l := len(mall)
		mall[v] = 1
		if l != len(mall) { //长度变化，即可以存
			l = len(mall)
		} else { //存不了，进并集
			set = append(set, v)
		}
	}
	//3.遍历交集，在并集中找，找到就从并集中删，删完后就是补集（即并-交=所有变化的元素）
	for _, v := range set {
		delete(mall, v)
	}
	//4.此时，mall是补集，所有元素去源中找，找到就是删除的，找不到的必定能在目数组中找到，即新加的
	var added, deleted []string
	for v, _ := range mall {
		_, exist := msrc[v]
		if exist {
			deleted = append(deleted, v)
		} else {
			added = append(added, v)
		}
	}

	return added, deleted
}

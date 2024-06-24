package postmaninternal

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/rbretecher/go-postman-collection"
	"k8s.io/klog/v2"
)

// 解析query
func parseRequestUrlQuery(postmanRequestURLQuery interface{}) {
	if postmanRequestURLQuery != nil { // postmanURL.Query 是个type interface{}，坑！！！
		httpapiQuery := postmanRequestURLQuery.([]interface{})
		for i := range httpapiQuery {
			httpapiQueryArg := httpapiQuery[i]
			valuename := httpapiQueryArg.(map[string]interface{})["key"]
			valuecontent := httpapiQueryArg.(map[string]interface{})["value"]
			valuedisabled := httpapiQueryArg.(map[string]interface{})["disabled"]
			// args := Args{In: "query", Name: valuename.(string), Value: Value{From: "query", Content: valuename.(string)}}
			if valuedisabled != nil {
				if valuedisabled == true {
					if (valuename != "") && (valuecontent != "") {
						args := Args{In: "query", Name: valuename.(string), Value: Value{From: "literal", Content: valuecontent.(string)}}
						apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
					} else {
						args := Args{In: "query", Name: valuename.(string), Value: Value{From: "literal", Content: ""}}
						apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
					}
				}
			} else {

				if (valuename != "") && (valuecontent != "") {
					args := Args{In: "query", Name: valuename.(string), Value: Value{From: "literal", Content: valuecontent.(string)}}
					apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
				} else {
					args := Args{In: "query", Name: valuename.(string), Value: Value{From: "literal", Content: ""}}
					apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
				}
			}

			// klog.Infoln("__httpapiQueryArgs valuename is : ", valuename.(string))
			// klog.Infoln("__httpapiQueryArgs valuecontent is : ", valuecontent.(string))
		}
	}
}

func parseRequestBodyUrlencoded(postmanRequestBody *postman.Body) {
	if postmanRequestBody != nil {
		tempbodyjson := make(map[string]string) // 组建body的json对象
		requestBody := postmanRequestBody
		requestBodyURLEncoded := requestBody.URLEncoded.([]interface{})
		var requestBodyUrlencodedStruct RequestBodyUrlencodedStruct
		// 因为是[]interface{}接口类型，所以先转换成byteArr
		byteArr, err := json.Marshal(requestBodyURLEncoded)
		if err != nil {
			klog.Infoln("json.marshal failed, err:", err)
			return
		}
		bytestring := string(byteArr)
		// 构建成一个完整的json
		bytestring = "{\"urlencoded\": " + bytestring + "}"
		// 重新解析[]byte(bytestring)到结构体
		err = json.Unmarshal([]byte(bytestring), &requestBodyUrlencodedStruct)
		if err != nil {
			panic(err)
		}
		for i := range requestBodyUrlencodedStruct.Urlencoded {
			if requestBodyUrlencodedStruct.Urlencoded[i].Enable {
				tempbodyjson[requestBodyUrlencodedStruct.Urlencoded[i].Key] = requestBodyUrlencodedStruct.Urlencoded[i].Value
			}
		}
		bodyArgs := Args{In: "body", Name: "body", Value: Value{From: "json", Json: tempbodyjson}}
		apiHubHttpConf.Args = append(apiHubHttpConf.Args, bodyArgs)
	}
}

// 解析body Raw New
func parseRequestBodyRaw(postmanRequestBody *postman.Body) error {
	if postmanRequestBody != nil {
		requestBody := postmanRequestBody
		if requestBody.Raw != "" {
			klog.Infoln("requestBody.Raw: ", requestBody.Raw)
			if strings.Index(requestBody.Raw, "\r") > 0 { //基本判定为json格式
				// if !json.Valid([]byte(requestBody.Raw)) { // 含有 \r \n 的json
				requestBody.Raw = strings.Replace(requestBody.Raw, "\r", "", -1)
				requestBody.Raw = strings.Replace(requestBody.Raw, "\n", "", -1)
				requestBody.Raw = strings.Replace(requestBody.Raw, "\t", "", -1)
				requestBody.Raw = strings.Replace(requestBody.Raw, "\u00a0", "", -1)
				// }
				requestBody.Raw = strings.Replace(requestBody.Raw, " ", "", -1)
				requestBody.Raw = strings.Replace(requestBody.Raw, "{{", "\"{{.vars.", -1)
				requestBody.Raw = strings.Replace(requestBody.Raw, "}}", "}}\"", -1)

				requestBodyRawMap := make(map[string]string)
				err := json.Unmarshal([]byte(requestBody.Raw), &requestBodyRawMap)
				if err != nil {
					bodyArgs := Args{In: "body", Name: "body", Value: Value{From: "literal", Content: requestBody.Raw}}
					apiHubHttpConf.Args = append(apiHubHttpConf.Args, bodyArgs)
					return err
				}

				nameList := ""
				tempRequestRawArray := *new([]string) // 暂存MD5需要对比的全部body元素
				for _, value := range requestBodyRawMap {
					if strings.Index(value, "vars") > 0 {
						value = strings.Replace(value, "{{.vars.", "", -1)
						value = strings.Replace(value, "}}", "", -1)
						// requestBody.Raw = strings.Replace(requestBody.Raw, "}}", "}}\"", -1)
					}
					tempRequestRawArray = append(tempRequestRawArray, value)

				}
				tempbodyjson := requestBodyRawMap
				for key, value := range requestBodyRawMap {
					nameString := key
					contentStringFunc, tempIndex := GetStringBetweenSpecifySymbols(value, "vars.", "}}")
					if tempIndex != -1 && coversionFuncMap[contentStringFunc] == "md5" { // MD5 特殊，单独处理
						for a := range setEnvironmentVariableMD5Array { // 遍历MD5涉及变量
							for n, v := range tempbodyjson {
								tempstring, tempindex := GetStringBetweenSpecifySymbols(v, "vars.", "}}")
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
						tempLocalMD5Key, _ := arrcmp(tempRequestRawArray, setEnvironmentVariableMD5Array)
						for x := range tempLocalMD5Key {
							args := Args{In: "vars", Name: ("key" + strconv.Itoa(x)), Value: Value{From: "literal", Content: tempLocalMD5Key[x]}}
							apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
							nameList = nameList + " " + "key" + strconv.Itoa(x)
						}
						nameList = strings.TrimSpace(nameList)
						args := Args{In: "vars", Name: nameString, Value: Value{From: "func", Content: coversionFuncMap[contentStringFunc], Args: nameList}}
						apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
					} else if tempIndex != -1 {
						args := Args{In: "vars", Name: contentStringFunc, Value: Value{From: "func", Content: coversionFuncMap[contentStringFunc]}}
						apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
						// nameList = nameList + " " + coversionFuncMap[contentStringFunc]
					}
				}
				bodyArgs := Args{In: "body", Name: "body", Value: Value{From: "json", Json: requestBodyRawMap}}
				apiHubHttpConf.Args = append(apiHubHttpConf.Args, bodyArgs)
			} else { // 大概率字符串
				requestBody.Raw = strings.Replace(requestBody.Raw, "\r", "", -1)
				requestBody.Raw = strings.Replace(requestBody.Raw, "\n", "", -1)
				requestBody.Raw = strings.Replace(requestBody.Raw, "\t", "", -1)
				requestBody.Raw = strings.Replace(requestBody.Raw, "\u00a0", "", -1)
				bodyArgs := Args{In: "body", Name: "body", Value: Value{From: "literal", Content: requestBody.Raw}}
				apiHubHttpConf.Args = append(apiHubHttpConf.Args, bodyArgs)
			}
		} else {
			klog.Infoln("__parseRequestBodyRawError: Content null")
		}
	} else {
		apiHubHttpConf.Requestcontenttype = ""
	}
	return nil
}

// 获取postman Event中的全局变量和js函数对应的变量
func getPostmanEventFunc(postmanItem *postman.Items, preFuncKeyWord []string, keyWordGlobal string) {
	if (postmanItem == nil) || (preFuncKeyWord == nil) {
		return
	}
	for i := range postmanItem.Events { // 通常Events就一个，但是个数组形式，所以使用遍历
		if postmanItem.Events[i].Script.Type == "text/javascript" {
			for j := range postmanItem.Events[i].Script.Exec { // Exec中js命令一行是一个数组元素
				for k := range preFuncKeyWord { // 查找js命令行中有无js常见命令的关键字
					if ((strings.Index(postmanItem.Events[i].Script.Exec[j], preFuncKeyWord[k])) != -1) && ((strings.Index(postmanItem.Events[i].Script.Exec[j], "//")) == -1) {
						switch preFuncKeyWord[k] {
						case "getTime":
							keyWordString, index := getStringBetweenDoubleQuotationMarks(postmanItem.Events[i].Script.Exec[j])
							// 解析js代码中EnvironmentVariable，查找到getTime代码，确定本行生成time变量名称，赋值到coversionFuncMap
							keyWordString = strings.TrimSpace(keyWordString)
							if index != -1 && keyWordString != "" {
								coversionFuncMap[keyWordString] = preEventFuncReferenceMap[preFuncKeyWord[k]]
							}
						case "CryptoJS.MD5":
							keyWordString, index := GetStringBetweenSpecifySymbols(postmanItem.Events[i].Script.Exec[j], "var", "=")
							keyWordString = strings.TrimSpace(keyWordString)
							if index != -1 && keyWordString != "" {
								coversionFuncMap[keyWordString] = preEventFuncReferenceMap[preFuncKeyWord[k]]
							}
							getSetEnvironmentVariableMD5(postmanItem, postmanItem.Events[i].Script.Exec[j], keyWordString)
							// klog.Infoln("__postmanItem.Events[i].Script.Exec[j]", keyWordString)
						default:
						}
					}
					if ((strings.Index(postmanItem.Events[i].Script.Exec[j], keyWordGlobal)) != -1) && ((strings.Index(postmanItem.Events[i].Script.Exec[j], "//")) == -1) {
						keyWordStringKey, indexKey := getStringBetweenDoubleQuotationMarks(postmanItem.Events[i].Script.Exec[j])
						keyWordStringValue, indexValue := GetStringBetweenSpecifySymbols(postmanItem.Events[i].Script.Exec[j], ",", ")")
						if indexKey != -1 {
							preGlobalKeyMap[keyWordStringKey] = keyWordStringKey
						}
						if indexValue != -1 {
							preGlobalValueMap[keyWordStringKey] = keyWordStringValue
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

	tempEnvironmentVariable, _ := GetStringBetweenSpecifySymbols(postmanEventScriptExec, "(", ")")
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
				tempString, tempIndex := GetStringBetweenSpecifySymbols(setEnvironmentVariableMD5[i:], "+", "+")
				tempString2, tempIndex2 := GetStringBetweenSpecifySymbols(setEnvironmentVariableMD5[i:], "\"", "\"")
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

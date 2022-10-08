package postmaninternal

import (
	"strings"

	"github.com/rbretecher/go-postman-collection"
)

/*
func getPostmanFilesBytes_Old(postmanfileBytes *postman.Collection) {
	if postmanfileBytes != nil {
		for i := range postmanfileBytes.Items {
			if postmanfileBytes.Items[i].Items == nil {
				converOneRequest(postmanfileBytes.Items[i])
				if privatesExport {
					if len(apiHubHttpPrivates.Privates) != 0 {
						apiHubHttpConf.Private = postmanfileBytes.Info.Name + "_" + postmanfileBytes.Items[i].Name + "_key"
						generateApiHubPrivatesJson(postmanfileBytes, apiHubHttpConf.Private)
					}
				}
				generateApiHubJson(postmanfileBytes, "")
			} else {
				for j := range postmanfileBytes.Items[i].Items {
					converOneRequest(postmanfileBytes.Items[i].Items[j])
					if privatesExport {
						if len(apiHubHttpPrivates.Privates) != 0 {
							apiHubHttpConf.Private = postmanfileBytes.Info.Name + "_" + postmanfileBytes.Items[i].Name + "_" + postmanfileBytes.Items[i].Items[j].Name + "_key"
							generateApiHubPrivatesJson(postmanfileBytes, apiHubHttpConf.Private)
						}
					}
					generateApiHubJson(postmanfileBytes, postmanfileBytes.Items[i].Name)
				}
			}
		}
	}
}
*/

func getPostmanFilesBytes(postmanfileBytes *postman.Collection) string {
	if postmanfileBytes != nil && postmanfileBytes.Items != nil {
		for i := range postmanfileBytes.Items {
			if (postmanfileBytes.Items[i].Items == nil) && (converOneRequest(postmanfileBytes.Items[i]) == "") { // 若只有一级Items
				if privatesExport && (len(apiHubHttpPrivates.Privates) != 0) {
					apiHubHttpConf.Private = postmanfileBytes.Info.Name + "_" + postmanfileBytes.Items[i].Name + "_key"
					generateApiHubJson(apiHubPrivatesJsonPath, ReplaceName(apiHubHttpConf.Private))
				}
				tempName := postmanfileBytes.Info.Name + "_" + postmanfileBytes.Items[i].Name
				generateApiHubJson(apiHubJsonPath, ReplaceName(tempName))
			} else {
				for j := range postmanfileBytes.Items[i].Items {
					if (postmanfileBytes.Items[i].Items[j].Items == nil) && (converOneRequest(postmanfileBytes.Items[i].Items[j]) == "") { // 两级Items
						if privatesExport && (len(apiHubHttpPrivates.Privates) != 0) {
							apiHubHttpConf.Private = postmanfileBytes.Info.Name + "_" + postmanfileBytes.Items[i].Name + "_" + postmanfileBytes.Items[i].Items[j].Name + "_key"
							generateApiHubJson(apiHubPrivatesJsonPath, ReplaceName(apiHubHttpConf.Private))
						}
						tempName := postmanfileBytes.Info.Name + "_" + postmanfileBytes.Items[i].Name + "_" + postmanfileBytes.Items[i].Items[j].Name
						generateApiHubJson(apiHubJsonPath, ReplaceName(tempName))
					} else {
						for k := range postmanfileBytes.Items[i].Items[j].Items {
							if (postmanfileBytes.Items[i].Items[j].Items[k].Items == nil) && (converOneRequest(postmanfileBytes.Items[i].Items[j].Items[k]) == "") { // 三级Items
								if privatesExport && (len(apiHubHttpPrivates.Privates) != 0) {
									apiHubHttpConf.Private = postmanfileBytes.Info.Name + "_" + postmanfileBytes.Items[i].Name + "_" + postmanfileBytes.Items[i].Items[j].Name + "_" + postmanfileBytes.Items[i].Items[j].Items[k].Name + "_key"
									generateApiHubJson(apiHubPrivatesJsonPath, ReplaceName(apiHubHttpConf.Private))
								}
								tempName := postmanfileBytes.Info.Name + "_" + postmanfileBytes.Items[i].Name + "_" + postmanfileBytes.Items[i].Items[j].Name + "_" + postmanfileBytes.Items[i].Items[j].Items[k].Name
								generateApiHubJson(apiHubJsonPath, ReplaceName(tempName))
							} else {
								for x := range postmanfileBytes.Items[i].Items[j].Items[k].Items {
									if (postmanfileBytes.Items[i].Items[j].Items[k].Items[x].Items == nil) && (converOneRequest(postmanfileBytes.Items[i].Items[j].Items[k].Items[x]) == "") { // 三级Items
										if privatesExport && (len(apiHubHttpPrivates.Privates) != 0) {
											apiHubHttpConf.Private = postmanfileBytes.Info.Name + "_" + postmanfileBytes.Items[i].Name + "_" + postmanfileBytes.Items[i].Items[j].Name + "_" + postmanfileBytes.Items[i].Items[j].Items[k].Name + "_" + postmanfileBytes.Items[i].Items[j].Items[k].Items[x].Name + "_key"
											generateApiHubJson(apiHubPrivatesJsonPath, ReplaceName(apiHubHttpConf.Private))
										}
										tempName := postmanfileBytes.Info.Name + "_" + postmanfileBytes.Items[i].Name + "_" + postmanfileBytes.Items[i].Items[j].Name + "_" + postmanfileBytes.Items[i].Items[j].Items[k].Name + "_" + postmanfileBytes.Items[i].Items[j].Items[k].Items[x].Name
										generateApiHubJson(apiHubJsonPath, ReplaceName(tempName))
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return "postmanfileBytes.Items ERROR"
}

// 转换postman collection中一个request
func converOneRequest(postmanItem *postman.Items) string {
	if postmanItem == nil {
		return "postmanItem nil"
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

	if apiHubHttpConf.URL == "" {
		return "URL nil, invalid api"
	}

	return ""
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
				if privatesExport {
					args := Args{In: "header", Name: postmanRequest.Header[i].Key, Value: Value{From: "private", Content: postmanRequest.Header[i].Key}}
					apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
					privates := Privates{Name: postmanRequest.Header[i].Key, Value: postmanRequest.Header[i].Value}
					apiHubHttpPrivates.Privates = append(apiHubHttpPrivates.Privates, privates)
				} else {
					args := Args{In: "header", Name: postmanRequest.Header[i].Key, Value: Value{From: "literal", Content: postmanRequest.Header[i].Value}}
					apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
				}
			} else if postmanRequest.Header[i].Key == "Content-Type" {
				headerindex := strings.Index(postmanRequest.Header[i].Value, "/")
				if strings.Index("json,form,origin,none,text", postmanRequest.Header[i].Value[headerindex+1:]) > 0 {
					apiHubHttpConf.Requestcontenttype = postmanRequest.Header[i].Value[headerindex+1:]
				}
			}
		}
	}
	if postmanRequest.URL != nil {
		switch apiHubHttpConf.Method {
		case "GET":
			// 解析qury
			if postmanRequest.URL.Query != nil {
				parseRequestUrlQuery(postmanRequest.URL.Query)
			}
		case "POST":
			// 解析qury
			if postmanRequest.URL.Query != nil {
				parseRequestUrlQuery(postmanRequest.URL.Query)
			}
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
}

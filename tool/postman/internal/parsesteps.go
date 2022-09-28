package postmaninternal

import (
	"strings"

	"github.com/rbretecher/go-postman-collection"
)

func getPostmanFilesBytes(postmanfileBytes *postman.Collection) {
	if postmanfileBytes != nil {
		for i := range postmanfileBytes.Items {
			if postmanfileBytes.Items[i].Items == nil {
				if converOneRequest(postmanfileBytes.Items[i]) == "" {
					if privatesExport {
						if len(apiHubHttpPrivates.Privates) != 0 {
							apiHubHttpConf.Private = postmanfileBytes.Info.Name + "_" + postmanfileBytes.Items[i].Name + "_key"
							apiHubHttpConf.Private = strings.Replace(apiHubHttpConf.Private, " ", "_", -1)
							generateApiHubPrivatesJson(postmanfileBytes, apiHubHttpConf.Private)
						}
					}
					generateApiHubJson(postmanfileBytes, "")
				}
			} else {
				for j := range postmanfileBytes.Items[i].Items {
					if converOneRequest(postmanfileBytes.Items[i].Items[j]) == "" {
						if privatesExport {
							if len(apiHubHttpPrivates.Privates) != 0 {
								apiHubHttpConf.Private = postmanfileBytes.Info.Name + "_" + postmanfileBytes.Items[i].Name + "_" + postmanfileBytes.Items[i].Items[j].Name + "_key"
								apiHubHttpConf.Private = strings.Replace(apiHubHttpConf.Private, " ", "_", -1)
								generateApiHubPrivatesJson(postmanfileBytes, apiHubHttpConf.Private)
							}
						}
						postmanfileBytes.Items[i].Name = strings.Replace(postmanfileBytes.Items[i].Name, " ", "_", -1)
						generateApiHubJson(postmanfileBytes, postmanfileBytes.Items[i].Name)
					}
				}
			}
		}
	}
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
					parseRequestBodyRawNew(postmanRequest.Body)
				default:
				}
			}
		}
	}
}

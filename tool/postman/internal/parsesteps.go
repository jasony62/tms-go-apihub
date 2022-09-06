package postmaninternal

import (
	"strings"

	"github.com/rbretecher/go-postman-collection"
)

func getPostmanFilesBytes(postmanfileBytes *postman.Collection) {
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
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/neotoolkit/openapi"
	"k8s.io/klog"
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

type Value struct {
	From    string `json:"from"`
	Content string `json:"content"`
}
type Args struct {
	In    string `json:"in"`
	Name  string `json:"name"`
	Value Value  `json:"value"`
}

var apiHubHttpConf ApiHubHttpConf
var oapi openapi.OpenAPI

var swaggerPath string
var apiHubConfPath string

func init() {
	flag.StringVar(&swaggerPath, "from", "./from/", "指定swagger文件路径")
	flag.StringVar(&apiHubConfPath, "to", "./to/", "指定转换后的apiHubConf json文件路径")
}

func main() {
	convertSwaggerFiles(swaggerPath)
}

func convertSwaggerFiles(path string) {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		klog.Errorln(err)
		return
	}
	var prefix string
	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())

		if fileInfoList[i].IsDir() {
			klog.Infoln("Swagger子目录: ", fileName)
			prefix = fileInfoList[i].Name()
			convertSwaggerFiles(path + "/" + prefix)
		} else {
			if !strings.HasSuffix(fileName, ".yaml") && !strings.HasSuffix(fileName, ".json") {
				continue
			}
			klog.Infoln("######加载Swagger(*.yaml or .json)文件: ", fileName)
			fileBytes, err := ioutil.ReadFile(fileName)
			if err != nil {
				klog.Errorln(err)
				panic(err)
			}
			covertSwaggerToApihubConf(fileBytes)
		}
	}
}

func getServerUrl(api openapi.OpenAPI) {

	for _, s := range api.Servers {
		if strings.Contains(s.URL, "http:") {
			apiHubHttpConf.URL = s.URL
			klog.Infof("Servers URL is %s\n", s.URL)
			return
		}
	}
}

func parseParameters(params openapi.Parameters) {
	for _, param := range params {
		switch param.In {
		case "query":
			args := Args{In: "query", Name: param.Name, Value: Value{From: "query", Content: param.Name}}
			apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
		case "header":
			args := Args{In: "header", Name: param.Name, Value: Value{From: "header", Content: param.Name}}
			apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
		case "path":
			//替换路径，RESTful接口，暂时不支持
		}
	}
}

//Operation-RequestBody
func parseRequestBody(body *openapi.RequestBody) {
	if body.IsRef() {
		klog.Infoln("Operation-RequestBody has ref, Not supported!")
		return
	}

	for typeStr, media := range body.Content {
		if media.Schema.Ref != "" {
			var err error
			klog.Infoln("Operation-RequestBody.Content.Schema has ref, Start parse!")
			media.Schema, err = oapi.LookupSchemaByReference(media.Schema.Ref)
			if err != nil {
				klog.Infoln("Operation-RequestBody.Content.Schema ref parse error!")
				return
			}
		}
		parseRequestContentType(typeStr)
		parseSchema(media.Schema)
	}
}
func parseRequestContentType(typeStr string) {
	switch typeStr {
	case "application/x-www-form-urlencoded":
		apiHubHttpConf.Requestcontenttype = "form"
	case "application/json":
		apiHubHttpConf.Requestcontenttype = "json"
	default:
		apiHubHttpConf.Requestcontenttype = typeStr
	}
}
func parseSchema(schema openapi.Schema) {
	if schema.IsRef() {
		klog.Infoln("schema has a ref recursively, will NOT parse!")
		return
	}
	for param, property := range schema.Properties {
		if property != nil && property.Type == "string" {
			args := Args{In: "body", Name: param, Value: Value{From: "body", Content: param}}
			apiHubHttpConf.Args = append(apiHubHttpConf.Args, args)
		}
	}
}

// func parseResponses(responses openapi.Responses) {
// }

//Components-RequestBodies
func parseRequestBodies(bodies openapi.RequestBodies) {
	if bodies == nil {
		klog.Infoln("Components-RequestBodies is nil")
		return
	}
	for _, body := range bodies {
		parseRequestBody(body)
	}
}

func parsePathOperation(oper *openapi.Path) {
	if oper.Post != nil {
		apiHubHttpConf.Method = "post"
		if oper.Post.Parameters != nil {
			parseParameters(oper.Post.Parameters)
			parseRequestBody(&oper.Post.RequestBody)
			// parseResponses(oper.Post.Responses)
		}
	}
	if oper.Get != nil {
		apiHubHttpConf.Method = "get"
		if oper.Get.Parameters != nil {
			parseParameters(oper.Get.Parameters)
			parseRequestBody(&oper.Get.RequestBody)
		}
	}
	if oper.Delete != nil {
		apiHubHttpConf.Method = "delete"
	}
	if oper.Put != nil {
		apiHubHttpConf.Method = "put"
		if oper.Put.Parameters != nil {
			parseParameters(oper.Put.Parameters)
			parseRequestBody(&oper.Put.RequestBody)
		}
	}
}
func parsePaths(api openapi.OpenAPI) {
	for p, oper := range api.Paths {
		apiHubHttpConf.URL += p
		klog.Infof("Swagger Paths is %s\n", apiHubHttpConf.URL)
		parsePathOperation(oper)
	}
}

func covertSwaggerToApihubConf(fileBytes []byte) {
	var err error
	oapi, err = openapi.Parse(fileBytes)
	if err != nil {
		log.Fatalln(err)
	}

	getServerUrl(oapi)
	parsePaths(oapi)
	parseRequestBodies(oapi.Components.RequestBodies)
	generateApiHubConf(oapi)
	apiHubHttpConf = ApiHubHttpConf{}
}

func generateApiHubConf(api openapi.OpenAPI) {
	apiHubHttpConf.ID = api.Info.Title
	apiHubHttpConf.Description = api.Info.Description
	fileName := apiHubConfPath + apiHubHttpConf.ID
	byteHttpApi, err := json.Marshal(apiHubHttpConf)
	if err != nil {
		return
	}
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
	klog.Infoln("%%文件转换完成:", fileName)
}

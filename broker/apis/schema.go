package apis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
	"github.com/xeipuuv/gojsonschema"
	klog "k8s.io/klog/v2"
)

func schemaChecker(path string, schema *gojsonschema.Schema) int {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		klog.Errorln(err)
		return 500
	}

	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())

		if fileInfoList[i].IsDir() {
			klog.Infoln("Plugin子目录: ", fileName)
			schemaChecker(path+"/"+fileName, schema)
		} else {
			jsonContent, err := ioutil.ReadFile(fileName)
			if err != nil {
				panic(err.Error())
			}

			if !json.Valid(jsonContent) {
				str := "Json文件无效：" + fileName
				klog.Errorln(str)
				panic(str)
			}

			documentLoader := gojsonschema.NewStringLoader(string(jsonContent))
			result, err := schema.Validate(documentLoader)
			if err != nil {
				klog.Errorln(err)
				return 500
			}

			if !result.Valid() {
				fmt.Printf("%s is not valid. see errors :		\r\n", fileName)
				for _, desc := range result.Errors() {
					fmt.Printf("- %s		\r\n", desc)
				}
				fmt.Printf("\r\n")
				return 500
			}

		}
	}
	return 200
}

func confValidator(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		return nil, http.StatusInternalServerError
	}

	path, OK := params["schema"]
	if !OK {
		return nil, http.StatusInternalServerError
	}
	fileName := path + "/httpapi.json"
	schemaContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	if !json.Valid(schemaContent) {
		str := "Schema文件无效：" + fileName
		klog.Errorln(str)
		panic(str)
	}

	loader1 := gojsonschema.NewStringLoader(string(schemaContent))
	schema, err := gojsonschema.NewSchema(loader1)
	if err != nil {
		panic(err.Error())
	}

	return nil, schemaChecker(util.GetBasePath()+"/httpapis/", schema)
}

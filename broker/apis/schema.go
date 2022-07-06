package apis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
			klog.Infoln("Schema检查Json子目录: ", fileName)
			if schemaChecker(fileName, schema) != 200 {
				return 500
			}
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
					klog.Errorln("- %s		", desc)
				}
				klog.Errorln("")
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

	return loadSchemaDefData(path)
}

func loadSchemaDefData(path string) (interface{}, int) {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		klog.Errorln(err)
		return nil, http.StatusInternalServerError
	}

	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())

		klog.Infoln("Schema校验文件: ", fileName)
		if fileInfoList[i].IsDir() {
			loadSchemaDefData(fileName)
		} else {
			var apipath string
			if strings.Contains(fileInfoList[i].Name(), "httpapi") {
				apipath = "httpapis"
			} else if strings.Contains(fileInfoList[i].Name(), "flow") {
				apipath = "flows"
			} else if strings.Contains(fileInfoList[i].Name(), "right") {
				apipath = "rights"
			} else if strings.Contains(fileInfoList[i].Name(), "schedule") {
				apipath = "schedules"
			}

			schemaContent, err := ioutil.ReadFile(fileName)
			if err != nil {
				panic(err.Error())
			}
			if !json.Valid(schemaContent) {
				str := "Schema文件无效：" + fileName
				klog.Errorln(str)
				panic(str)
			}

			loader := gojsonschema.NewStringLoader(string(schemaContent))
			schema, err := gojsonschema.NewSchema(loader)
			if err != nil {
				panic(err.Error())
			}

			if schemaChecker(util.GetBasePath()+apipath, schema) != 200 {
				klog.Errorln("Schema检查json文件不合法，目录: ", util.GetBasePath()+apipath)
				return nil, http.StatusInternalServerError
			}
			klog.Infoln("Schema检查json文件合法，目录: ", util.GetBasePath()+apipath)
		}
	}
	return nil, 200
}

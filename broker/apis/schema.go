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
	"go.uber.org/zap"
)

func schemaChecker(path string, schema *gojsonschema.Schema) int {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		zap.S().Errorln(err.Error())
		return 500
	}

	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())

		if fileInfoList[i].IsDir() {
			//	zap.S().Infoln("Schema检查Json子目录: ", fileName)
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
				zap.S().Errorln(str)
				panic(str)
			}

			documentLoader := gojsonschema.NewStringLoader(string(jsonContent))
			result, err := schema.Validate(documentLoader)
			if err != nil {
				zap.S().Errorln(err.Error())
				return 500
			}

			if !result.Valid() {
				fmt.Printf("%s is not valid. see errors :		\r\n", fileName)
				for _, desc := range result.Errors() {
					zap.S().Errorln("- %s		", desc)
				}
				zap.S().Errorln("")
				return 500
			}

		}
	}
	return 200
}

func confValidator(stack *hub.Stack, params map[string]string) (interface{}, int) {
	if len(params) == 0 {
		str := "confValidator 参数为空"
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	path, OK := params["schema"]
	if !OK {
		str := "confValidator schema为空"
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
	}

	return loadSchemaDefData(path)
}

func loadSchemaDefData(path string) (interface{}, int) {
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		zap.S().Errorln(err.Error())
		return util.CreateTmsError(hub.TmsErrorApisId, err.Error(), nil), http.StatusInternalServerError
	}

	zap.S().Infoln("校验Schema文件...")
	for i := range fileInfoList {
		fileName := fmt.Sprintf("%s/%s", path, fileInfoList[i].Name())

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
				zap.S().Errorln(str)
				panic(str)
			}

			loader := gojsonschema.NewStringLoader(string(schemaContent))
			schema, err := gojsonschema.NewSchema(loader)
			if err != nil {
				panic(err.Error())
			}

			if schemaChecker(util.GetBasePath()+apipath, schema) != 200 {
				str := "Schema检查json文件不合法，目录: " + util.GetBasePath() + apipath
				return util.CreateTmsError(hub.TmsErrorApisId, str, nil), http.StatusInternalServerError
			}
		}
	}
	return nil, http.StatusOK
}

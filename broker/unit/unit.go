package unit

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/plugin"
)

func loadPrivateData(path string, bucket string, name string) (*hub.PrivateArray, error) {
	var filePath string
	if bucket != "" {
		filePath = fmt.Sprintf("%s/%s/%s.json", path, bucket, name)
	} else {
		filePath = fmt.Sprintf("%s/%s.json", path, name)
	}
	filePtr, err := os.Open(filePath)
	if err != nil {
		log.Panic("获得API定义失败：", err)
		return nil, err
	}
	defer filePtr.Close()

	result := new(hub.PrivateArray)
	decoder := json.NewDecoder(filePtr)
	decoder.Decode(result)
	return result, nil
}

func FindApiDef(stack *hub.Stack, bucket string, id string) (*hub.ApiDef, error) {
	var filePath string
	if bucket != "" {
		filePath = fmt.Sprintf("%s/%s/%s.json", stack.ApiDefPath, bucket, id)
	} else {
		filePath = fmt.Sprintf("%s/%s.json", stack.ApiDefPath, id)
	}
	filePtr, err := os.Open(filePath)
	if err != nil {
		log.Panic("获得API定义失败：", err)
		return nil, err
	}
	defer filePtr.Close()

	apiDef := new(hub.ApiDef)
	decoder := json.NewDecoder(filePtr)
	decoder.Decode(apiDef)

	if len(apiDef.PrivateName) > 0 {
		//需要load秘钥
		apiDef.Privates, err = loadPrivateData(stack.PrivateDefPath, bucket, apiDef.PrivateName)
		if err != nil {
			log.Panic("获得Private数据失败：", err)
			return nil, err
		}
	}

	// 通过插件改写API定义
	if apiDef.Plugins != nil && len(*apiDef.Plugins) > 0 {
		for _, pluginDef := range *apiDef.Plugins {
			if len(pluginDef.Path) > 0 {
				pluginFn, _ := plugin.RewriteApiDef(pluginDef.Path)
				if pluginFn != nil {
					pluginFn(apiDef)
				}
			}
		}
	}

	return apiDef, nil
}

func FindFlowDef(stack *hub.Stack, bucket string, id string) (*hub.FlowDef, error) {
	var filePath string
	if bucket != "" {
		filePath = fmt.Sprintf("%s/%s/%s.json", stack.FlowDefPath, bucket, id)
	} else {
		filePath = fmt.Sprintf("%s/%s.json", stack.FlowDefPath, id)
	}
	filePtr, err := os.Open(filePath)
	if err != nil {
		log.Panic("获得Flow定义失败", err)
		return nil, err
	}
	defer filePtr.Close()

	flowDef := new(hub.FlowDef)
	decoder := json.NewDecoder(filePtr)
	decoder.Decode(&flowDef)

	return flowDef, nil
}

func FindPrivateValue(api *hub.ApiDef, name string) string {
	for _, pair := range *api.Privates.Pairs {
		if pair.Name == name {
			return pair.Value
		}
	}
	return ""
}

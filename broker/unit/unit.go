package unit

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/plugin"
)

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

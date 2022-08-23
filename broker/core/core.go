package core

import (
	"time"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/logger"
	"github.com/jasony62/tms-go-apihub/util"
)

func init() {
	logger.LogS().Infoln("Core register apis\n")
	RegisterApis(map[string]hub.ApiHandler{"flowApi": runFlowApi,
		"scheduleApi": runScheduleApi,
	})
}

func ApiHubStartMainFlow(path string) {
	util.LoadMainFlow(path)
	var stack hub.Stack
	stack.BaseString = " base: main. "
	stack.StartTime = time.Now()

	ApiRun(&stack, &hub.ApiDef{Name: "main", Command: "flowApi",
		Args: &[]hub.BaseParamDef{{Name: "name", Value: hub.BaseValueDef{From: "literal", Content: "main"}}}}, "", false)
}

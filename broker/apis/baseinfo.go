package apis

import (
	"github.com/google/uuid"
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/util"
)

func fillBaseInfo(stack *hub.Stack, params map[string]string) (interface{}, int) {

	/* 获取基本信息到base map，例如：
	 * "root":"amap_city_weather"
	 * "type":"flow"
	 * "start":"1659432682"
	 */
	base := stack.Heap[hub.HeapBaseName].(map[string]interface{})

	if base == nil {
		str := "获得Base map失败" + hub.HeapBaseName
		return util.CreateTmsError(hub.TmsErrorApisId, str, nil), 500
	}
	/* 添加信息，例如：
	 * "user":""
	 * "uuid":""
	 */
	for k, v := range params {
		base[k] = v
	}

	// 若json未提供uuid，创建uuid字符串，作为唯一请求标识符，例如："e1b86e64-b26c-4b7f-bdd1-7ef9492b8780"
	if len(params["uuid"]) == 0 {
		base["uuid"] = uuid.New().String()
	}

	stack.BaseString = util.CreateBaseString(base)
	return nil, 200
}
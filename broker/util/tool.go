package util

import "encoding/json"

func MapToString(param map[string]interface{}) string {
	//json转map
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

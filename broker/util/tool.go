package util

import "encoding/json"

func CreateBaseString(param map[string]interface{}) string {
	//json转map
	result := []byte(" base :")
	dataType, _ := json.Marshal(param)
	dataType = append(dataType, '.', ' ')
	result = append(result, dataType...)
	return string(result)
}

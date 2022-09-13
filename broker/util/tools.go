package util

import (
	"encoding/json"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateBaseString(param map[string]interface{}) string {
	//json转map
	result := []byte(" base :")
	dataType, _ := json.Marshal(param)
	dataType = append(dataType, '.', ' ')
	result = append(result, dataType...)
	return string(result)
}

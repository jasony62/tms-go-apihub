package util

// import (
// 	"strconv"
// )

// func GetInterfaceToInt(in interface{}) int {
// 	var value int
// 	switch in.(type) {
// 	case uint:
// 		value = int(in.(uint))
// 	case int8:
// 		value = int(in.(int8))
// 	case uint8:
// 		value = int(in.(uint8))
// 	case int16:
// 		value = int(in.(int16))
// 	case uint16:
// 		value = int(in.(uint16))
// 	case int32:
// 		value = int(in.(int32))
// 	case uint32:
// 		value = int(in.(uint32))
// 	case int64:
// 		value = int(in.(int64))
// 	case uint64:
// 		value = int(in.(uint64))
// 	case float32:
// 		value = int(in.(float32))
// 	case float64:
// 		value = int(in.(float64))
// 	case string:
// 		value, _ = strconv.Atoi(in.(string))
// 	default:
// 		value = in.(int)
// 	}
// 	return value
// }

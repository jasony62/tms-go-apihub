package main

import (
	"errors"
	"fmt"
	postmaninternal "postmaninternal"
	"sync"
	"syscall/js"
)

var wg sync.WaitGroup // number of working goroutines

/* wasm对外暴露的接口，供nodejs调用
   入参：postman collection 的json字符串；
		 nodejs提供的回调函数，当wasm转换完数据格式后，调用回调把结果输出
*/
func postmanToHttpapis(this js.Value, args []js.Value) interface{} {
	callback := args[len(args)-1]
	go func() {
		json, _ := convertJson(args[0].String())
		callback.Invoke(json)
		defer wg.Done()
	}()
	return nil
}

func main() {
	wg.Add(1)
	js.Global().Set("postmanToHttpapis", js.FuncOf(postmanToHttpapis))
	fmt.Println("postmanToHttpapis registered")
	wg.Wait()
}

func convertJson(str string) (string, error) {
	jsonArr, err := postmaninternal.HandleStringFlows(str, false)
	if err != nil {
		fmt.Println("convertJson failed: ", err)
		return "", err
	}
	count := len(jsonArr)
	fmt.Println("convertJson OK, total: ", count)
	if count == 0 {
		fmt.Println("convert count is zero")
		return "", errors.New("convert count is zero")
	}
	//转成json类型的数组
	var json string
	json += "["
	for i := 0; i < count; i++ {
		json += jsonArr[i]
		if i == count-1 {
			json += "]"
		} else {
			json += ","
		}
	}
	return json, nil
}

package postmaninternal

import (
	"encoding/json"
	"strings"

	"github.com/rbretecher/go-postman-collection"
	"k8s.io/klog/v2"
)

func readWasmPostman(stringflows string) ([]string, error) {

	// tempString, _ := ioutil.ReadFile("./postman_collections//5G新增手机终端画像.postman_collection.json")
	// postmanString := string(stringflows)

	// var tempIoReaderString io.Reader
	tempIoReaderString := strings.NewReader(stringflows)
	// Parse the contents
	postmanfileBytes, err := postman.ParseCollection(tempIoReaderString)
	if err != nil {
		klog.Errorln(err)
		// panic(err)
		return nil, err
	}
	return getPostmanFilesBytesWasm(postmanfileBytes)
}

func getPostmanFilesBytesWasm(postmanfileBytes *postman.Collection) ([]string, error) {
	var outputJsonArray []string
	if postmanfileBytes != nil {
		for i := range postmanfileBytes.Items {
			if postmanfileBytes.Items[i].Items == nil {
				converOneRequest(postmanfileBytes.Items[i])
				tempString, err := outputJsonString()
				if err != nil {
					klog.Errorln("outputJsonString failed:", err)
				}
				outputJsonArray = append(outputJsonArray, tempString)
			} else {
				for j := range postmanfileBytes.Items[i].Items {
					converOneRequest(postmanfileBytes.Items[i].Items[j])
					tempString, err := outputJsonString()
					if err != nil {
						klog.Errorln("outputJsonString failed:", err)
					}
					outputJsonArray = append(outputJsonArray, tempString)
				}
			}
		}
	}
	return outputJsonArray, nil
}

func outputJsonString() (string, error) {
	byteHttpApi, err := json.Marshal(apiHubHttpConf)
	if err != nil {
		klog.Errorln("json.Marshal失败!")
		return "", err
	}
	return string(byteHttpApi), nil
}

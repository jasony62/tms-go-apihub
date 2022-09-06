package postmaninternal

import (
	"encoding/json"
	"strings"

	"github.com/rbretecher/go-postman-collection"
	"k8s.io/klog/v2"
)

func readWasmPostman(stringflows string) (string, error) {

	// tempString, _ := ioutil.ReadFile("./postman_collections//5G新增手机终端画像.postman_collection.json")
	// postmanString := string(stringflows)

	// var tempIoReaderString io.Reader
	tempIoReaderString := strings.NewReader(stringflows)
	// Parse the contents
	postmanfileBytes, err := postman.ParseCollection(tempIoReaderString)
	if err != nil {
		klog.Errorln(err)
		// panic(err)
		return "", err
	}
	return getPostmanFilesBytesWasm(postmanfileBytes)
}

func getPostmanFilesBytesWasm(postmanfileBytes *postman.Collection) (string, error) {
	if postmanfileBytes != nil {
		for i := range postmanfileBytes.Items {
			if postmanfileBytes.Items[i].Items == nil {
				converOneRequest(postmanfileBytes.Items[i])
				return outputJsonString()
			} else {
				for j := range postmanfileBytes.Items[i].Items {
					converOneRequest(postmanfileBytes.Items[i].Items[j])
					return outputJsonString()
				}
			}
		}
	}
	return "", nil
}

func outputJsonString() (string, error) {
	byteHttpApi, err := json.Marshal(apiHubHttpConf)
	if err != nil {
		klog.Errorln("json.Marshal失败!")
		return "", err
	}
	return string(byteHttpApi), nil
}

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
	if postmanfileBytes != nil && postmanfileBytes.Items != nil {
		for i := range postmanfileBytes.Items {
			if (postmanfileBytes.Items[i].Items == nil) && (converOneRequest(postmanfileBytes.Items[i]) == "") { // 若只有一级Items

				apiHubHttpConf.ID = strings.Replace(postmanfileBytes.Info.Name+"_"+postmanfileBytes.Items[i].Name, " ", "_", -1)
				apiHubHttpConf.ID = ReplaceName(apiHubHttpConf.ID)
				tempString, err := outputJsonString()
				if err != nil {
					klog.Errorln("outputJsonString failed:", err)
				}
				outputJsonArray = append(outputJsonArray, tempString)

			} else {
				for j := range postmanfileBytes.Items[i].Items {
					if (postmanfileBytes.Items[i].Items[j].Items == nil) && (converOneRequest(postmanfileBytes.Items[i].Items[j]) == "") { // 两级Items
						apiHubHttpConf.ID = strings.Replace(postmanfileBytes.Items[i].Name+"_"+postmanfileBytes.Items[i].Items[j].Name, " ", "_", -1)
						apiHubHttpConf.ID = ReplaceName(apiHubHttpConf.ID)
						tempString, err := outputJsonString()
						if err != nil {
							klog.Errorln("outputJsonString failed:", err)
						}
						outputJsonArray = append(outputJsonArray, tempString)
					} else {
						for k := range postmanfileBytes.Items[i].Items[j].Items {
							if (postmanfileBytes.Items[i].Items[j].Items[k].Items == nil) && (converOneRequest(postmanfileBytes.Items[i].Items[j].Items[k]) == "") { // 三级Items
								apiHubHttpConf.ID = strings.Replace(postmanfileBytes.Items[i].Items[j].Name+"_"+postmanfileBytes.Items[i].Items[j].Items[k].Name, " ", "_", -1)
								apiHubHttpConf.ID = ReplaceName(apiHubHttpConf.ID)
								tempString, err := outputJsonString()
								if err != nil {
									klog.Errorln("outputJsonString failed:", err)
								}
								outputJsonArray = append(outputJsonArray, tempString)
							} else {
								for x := range postmanfileBytes.Items[i].Items[j].Items[k].Items {
									if (postmanfileBytes.Items[i].Items[j].Items[k].Items[x].Items == nil) && (converOneRequest(postmanfileBytes.Items[i].Items[j].Items[k].Items[x]) == "") { // 三级Items
										apiHubHttpConf.ID = strings.Replace(postmanfileBytes.Items[i].Items[j].Items[k].Name+"_"+postmanfileBytes.Items[i].Items[j].Items[k].Items[x].Name, " ", "_", -1)
										apiHubHttpConf.ID = ReplaceName(apiHubHttpConf.ID)
										tempString, err := outputJsonString()
										if err != nil {
											klog.Errorln("outputJsonString failed:", err)
										}
										outputJsonArray = append(outputJsonArray, tempString)
									} else {
										for y := range postmanfileBytes.Items[i].Items[j].Items[k].Items[x].Items {
											if (postmanfileBytes.Items[i].Items[j].Items[k].Items[x].Items[y].Items == nil) && (converOneRequest(postmanfileBytes.Items[i].Items[j].Items[k].Items[x].Items[y]) == "") { // 四级Items
												apiHubHttpConf.ID = strings.Replace(postmanfileBytes.Items[i].Items[j].Items[k].Items[x].Name+"_"+postmanfileBytes.Items[i].Items[j].Items[k].Items[x].Items[y].Name, " ", "_", -1)
												apiHubHttpConf.ID = ReplaceName(apiHubHttpConf.ID)
												tempString, err := outputJsonString()
												if err != nil {
													klog.Errorln("outputJsonString failed:", err)
												}
												outputJsonArray = append(outputJsonArray, tempString)
											}
										}
									}
								}
							}
						}
					}
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

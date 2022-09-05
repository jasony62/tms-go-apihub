package postmaninternal

func HandleStringFlows(stringflows string, privatesSwitch bool) (string, error) {
	privatesExport = privatesSwitch
	return readWasmPostman(stringflows)
}

// func main() {
// 	tempString, _ := ioutil.ReadFile("./postman_collections//5G新增手机终端画像.postman_collection.json")
// 	postmanString := string(tempString)
// 	m, err := postmaninternal.InputStringFlows(postmanString)
// 	_ = m
// 	_ = err
// 	println("Reading", m)
// }

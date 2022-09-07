1. 如何生成wasm文件
GOARCH=wasm GOOS=js go build -o postman.wasm main.go

2. 如何使用nodejs调用wasm文件

node wasm_exec_node.js ./postman.wasm

3. wasm_exec_node.js 为调用wasm文件的示例代码，里面读取postman json文件：“5G新增手机终端画像.postman_collection.json”， 返回生成的httpapi的json数组。
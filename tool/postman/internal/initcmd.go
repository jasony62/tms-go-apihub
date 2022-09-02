package postmaninternal

import (
	"flag"
)

// 初始化
func init() {
	flag.StringVar(&postmanPath, "from", "./postman_collections/", "指定postman_collections文件路径")
	flag.StringVar(&apiHubJsonPath, "to", "./httpapis/", "指定转换后的apiHub json文件路径")
	flag.StringVar(&apiHubPrivatesJsonPath, "private", "./privates/", "指定转换后的apiHub privates json文件路径")
}
func ConvertPostman() {

	privatesExport = true
	convertPostmanFiles(postmanPath)
}

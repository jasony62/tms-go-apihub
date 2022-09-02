package postmaninternal

type ApiHubHttpConf struct {
	ID                 string `json:"id"`
	Description        string `json:"description"`
	URL                string `json:"url"`
	Method             string `json:"method"`
	Private            string `json:"private,omitempty"`
	Requestcontenttype string `json:"requestContentType"`
	Args               []Args `json:"args,omitempty"`
}

type Args struct {
	In    string `json:"in"`
	Name  string `json:"name"`
	Value Value  `json:"value,omitempty"`
}

type Value struct {
	From    string            `json:"from"`
	Content string            `json:"content,omitempty"`
	Args    string            `json:"args,omitempty"`
	Json    map[string]string `json:"json,omitempty"`
	// Json    *interface{} `json:"json,omitempty"`
}

type ApiHubHttpPrivates struct {
	Privates []Privates `json:"privates"`
}

type Privates struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// 创建list,提前预设
// 左：postman js关键字
// 右：apihub内部函数名称
var preEventFuncReferenceMap = map[string]string{
	"getTime":      "utc_ms",
	"CryptoJS.MD5": "md5",
}

// 创建list,提前预设
// postman js关键字
var preEventFuncKeyMap = []string{
	"getTime",
	"CryptoJS.MD5",
}

// 相当于postman脚本中全局变量转换的一个中间量，映射postman脚本requset中全局变量值到apihub内部函数名称
// coversionFuncMap[time] = preEventFuncReferenceMap["getTime"] = "utc"
var coversionFuncMap map[string]string

// pr部分的全部全局变量Key Map
var preGlobalKeyMap map[string]string

// pr部分的全部全局变量Value Map
var preGlobalValueMap map[string]string

// postman文件路径
var PostmanPath string

// 导出json路径
var ApiHubJsonPath string
var ApiHubPrivatesJsonPath string

// json结构体
var apiHubHttpConf ApiHubHttpConf
var apiHubHttpPrivates ApiHubHttpPrivates

// Event中MD5涉及的变量内容组成的字符串数组
var setEnvironmentVariableMD5Array []string

type RequestBodyUrlencoded struct {
	Enable bool   `json:"enable"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type RequestBodyUrlencodedStruct struct {
	Urlencoded []RequestBodyUrlencoded `json:"urlencoded"`
}

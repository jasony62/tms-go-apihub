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
	From    string `json:"from"`
	Content string `json:"content,omitempty"`
	Args    string `json:"args,omitempty"`
	// Json    map[string]string `json:"json,omitempty"`
	Json interface{} `json:"json,omitempty"`
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

var privatesExport bool

// pr部分的全部全局变量Key Map
var preGlobalKeyMap map[string]string

// pr部分的全部全局变量Value Map
var preGlobalValueMap map[string]string

// postman文件路径
var postmanPath string

// 导出json路径
var apiHubJsonPath string
var apiHubPrivatesJsonPath string

// json结构体
var apiHubHttpConf ApiHubHttpConf
var apiHubHttpPrivates ApiHubHttpPrivates

// Event js中MD5涉及的变量内容组成的字符串数组
/*例如
[0]:"105**6"
[1]:"time"
[2]:"182******47"
[3]:"【签名测试】这是一条测试短信"
[4]:"234"
[5]:""
[6]:"04d6**********5c"
*/
var setEnvironmentVariableMD5Array []string

// 相当于postman脚本中全局变量转换的一个中间量，映射postman脚本requset中全局变量值到apihub内部函数名称
// coversionFuncMap[time] = preEventFuncReferenceMap["getTime"] = "utc"
/* 例如
"time":"utc_ms"
"sign":"md5"
*/
var coversionFuncMap map[string]string

type RequestBodyUrlencoded struct {
	Enable bool   `json:"enable"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type RequestBodyUrlencodedStruct struct {
	Urlencoded []RequestBodyUrlencoded `json:"urlencoded"`
}

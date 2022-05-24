package hub

// 应用的基本信息
type App struct {
	Host         string
	Port         int
	BucketEnable bool
	ApiMap       map[string]*ApiDef
	PrivateMap   map[string]*PrivateArray
	FlowMap      map[string]*FlowDef
	ScheduleMap  map[string]*ScheduleDef
}

type BaseDefParamValue struct {
	From    string       `json:"from"`
	Content string       `json:"content"`
	Json    *interface{} `json:"json"`
	Args    string       `json:"args"`
}

var DefaultApp App

//当from.from为"funcs"时，直接调用函数
var FuncMap map[string](interface{})

//用于template调用Funcs时，解析函数并调用
var FuncMapForTemplate map[string](interface{})

const OriginName = "origin"
const VarsName = "vars"
const LoopName = "loop"
const ResultName = "result"

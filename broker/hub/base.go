package hub

// 应用的基本信息
type App struct {
	Host         string
	Port         int
	BucketEnable bool
	ApiMap       map[string]*HttpApiDef
	PrivateMap   map[string]*PrivateArray
	FlowMap      map[string]*FlowDef
	ScheduleMap  map[string]*ScheduleDef
	SourceMap    map[string]string
}

var DefaultApp App

type BaseValueDef struct {
	From    string       `json:"from"`
	Content string       `json:"content"`
	Json    *interface{} `json:"json"`
	Args    string       `json:"args"`
}

type BaseParamDef struct {
	Name  string       `json:"name"`
	Value BaseValueDef `json:"value,omitempty"`
}
type TaskHandler func(*Stack, map[string]string) (interface{}, int)
type FuncHandler func(params []string) string
type TemplateHandler func(args ...interface{}) string

//当from.from为"funcs"时，直接调用函数
var FuncMap map[string]FuncHandler

//用于template调用Funcs时，解析函数并调用
var FuncMapForTemplate map[string](interface{})

const OriginName = "origin"
const VarsName = "vars"
const LoopName = "loop"
const ResultName = "result"

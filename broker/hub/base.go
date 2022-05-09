package hub

// 应用的基本信息
type App struct {
	Host         string
	Port         int
	BucketEnable bool
	ApiMap       map[string]ApiDef
	PrivateMap   map[string]PrivateArray
	FlowMap      map[string]FlowDef
	ScheduleMap  map[string]ScheduleDef
}

var DefaultApp App

const OriginName = "origin"
const VarsName = "vars"

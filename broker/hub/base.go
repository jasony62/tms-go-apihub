package hub

// 应用的基本信息
type App struct {
	Host            string
	Port            int
	BucketEnable    bool
	ApiDefPath      string
	PrivateDefPath  string
	FlowDefPath     string
	ScheduleDefPath string
}

var DefaultApp App

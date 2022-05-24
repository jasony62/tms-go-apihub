package hub

type ScheduleSwitchCaseDef struct {
	Value         string             `json:"value"`
	ConcurrentNum int                `json:"concurrentNum,omitempty"`
	Tasks         *[]ScheduleTaskDef `json:"tasks"`
}

type ScheduleTaskDef struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ResultKey   string `json:"resultKey"`
	//用于switch和loop
	Key BaseDefParamValue `json:"key"`
	//用于switch和loop
	ConcurrentNum int `json:"concurrentNum,omitempty"`
	//用于loop
	ConcurrentLoopNum int `json:"concurrentLoopNum,omitempty"`
	//用于switch，loop，api，flow
	Concurrent bool `json:"concurrent,omitempty"`
	//用于switch
	Cases *[]ScheduleSwitchCaseDef `json:"cases,omitempty"`
	//用于loop
	Tasks *[]ScheduleTaskDef `json:"tasks,omitempty"`
	//用于api，flow
	Parameters *[]OriginDefParam `json:"parameters,omitempty"`
	//api
	PrivateName string `json:"private"`
}

type ScheduleDef struct {
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	ConcurrentNum int                `json:"concurrentNum"`
	Tasks         *[]ScheduleTaskDef `json:"tasks"`
}

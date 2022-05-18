package hub

type ScheduleSwitchCaseDef struct {
	Value         string             `json:"value"`
	ConcurrentNum int                `json:"concurrentNum,omitempty"`
	Tasks         *[]ScheduleTaskDef `json:"tasks"`
}

type ScheduleTaskDef struct {
	Type              string          `json:"type"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	ResultKey         string          `json:"resultKey"`
	Key               ApiDefParamFrom `json:"key"`
	ConcurrentNum     int             `json:"concurrentNum,omitempty"`
	ConcurrentLoopNum int             `json:"concurrentLoopNum,omitempty"`
	Concurrent        bool            `json:"concurrent,omitempty"`
	//用于switch
	Cases      *[]ScheduleSwitchCaseDef `json:"cases,omitempty"`
	Tasks      *[]ScheduleTaskDef       `json:"tasks,omitempty"`
	Parameters *[]OriginDefParam        `json:"parameters,omitempty"`
}

type ScheduleDef struct {
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	ConcurrentNum int                `json:"concurrentNum"`
	Tasks         *[]ScheduleTaskDef `json:"tasks"`
}

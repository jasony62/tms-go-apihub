package hub

type ScheduleSwitchCaseDef struct {
	Value string             `json:"value"`
	Tasks *[]ScheduleTaskDef `json:"tasks"`
}

type ScheduleTaskDef struct {
	Type        string          `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	ResultKey   string          `json:"resultKey"`
	Commond     string          `json:"command"`
	Key         ApiDefParamFrom `json:"key"`
	Concurrent  int             `json:"Concurrent,omitempty"`
	//用于switch
	Cases      *[]ScheduleSwitchCaseDef `json:"cases,omitempty"`
	Tasks      *[]ScheduleTaskDef       `json:"tasks,omitempty"`
	Parameters *[]OriginDefParam        `json:"parameters,omitempty"`
}

type ScheduleDef struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Tasks       *[]ScheduleTaskDef `json:"tasks"`
}

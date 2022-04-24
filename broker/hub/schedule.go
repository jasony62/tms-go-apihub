package hub

type ScheduleSwitchCaseDef struct {
	Value string             `json:"value"`
	Tasks *[]ScheduleTaskDef `json:"tasks"`
}

type ScheduleDefParam struct {
	Name  string           `json:"name"`
	Value string           `json:"value,omitempty"`
	From  *ApiDefParamFrom `json:"from,omitempty"`
}

type ScheduleTaskDef struct {
	Type        string          `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	ResultKey   string          `json:"resultKey"`
	Commond     string          `json:"command"`
	Key         ApiDefParamFrom `json:"key"`
	//用于switch
	Cases      *[]ScheduleSwitchCaseDef `json:"cases"`
	Parameters *[]ScheduleDefParam
}

type ScheduleDef struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Tasks       []ScheduleTaskDef `json:"tasks"`
}

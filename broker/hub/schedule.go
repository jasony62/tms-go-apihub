package hub

type ScheduleSwitchCaseDef struct {
	Value string             `json:"value"`
	Tasks *[]ScheduleTaskDef `json:"tasks"`
}

type ScheduleTaskDef struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ResultKey   string `json:"resultKey"`
	Commond     string `json:"command"`
	Times       string
	Key         ApiDefParamFrom `json:"key"`
	//用于switch
	Cases      *[]ScheduleSwitchCaseDef `json:"cases"`
	Tasks      *[]ScheduleTaskDef       `json:"tasks"`
	Parameters *[]OriginDefParam
}

type ScheduleDef struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Tasks       *[]ScheduleTaskDef `json:"tasks"`
}

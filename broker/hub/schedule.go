package hub

type ScheduleSwitchCaseDef struct {
	Value         string            `json:"value"`
	ConcurrentNum int               `json:"concurrentNum,omitempty"`
	Steps         *[]ScheduleApiDef `json:"steps"`
}

type ScheduleControlDef struct {
	Name              string                   `json:"name"`
	Private           string                   `json:"private"`
	Description       string                   `json:"description"`
	ResultKey         string                   `json:"resultKey"`
	Key               BaseValueDef             `json:"key"`
	ConcurrentNum     int                      `json:"concurrentNum,omitempty"`
	ConcurrentLoopNum int                      `json:"concurrentLoopNum,omitempty"`
	Cases             *[]ScheduleSwitchCaseDef `json:"cases,omitempty"`
	Steps             *[]ScheduleApiDef        `json:"steps,omitempty"`
}
type ScheduleApiDef struct {
	Type string `json:"type"`
	Mode string `json:"mode"`
	/*只用于Api*/
	Private string              `json:"private"`
	Api     *ApiDef             `json:"api"`
	Control *ScheduleControlDef `json:"control"`
}

type ScheduleDef struct {
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	ConcurrentNum int               `json:"concurrentNum"`
	Steps         *[]ScheduleApiDef `json:"steps"`
}

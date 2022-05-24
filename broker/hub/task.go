package hub

type FlowStepResponseDef struct {
	Type string             `json:"type"`
	From *BaseDefParamValue `json:"from,omitempty"`
}

type TaskDef struct {
	Name             string               `json:"name"`
	Command          string               `json:"command"`
	Description      string               `json:"description"`
	ResultKey        string               `json:"resultKey"`
	Parameters       *[]OriginDefParam    `json:"parameters"`
	OriginParameters *[]OriginDefParam    `json:"origin"`
	Response         *FlowStepResponseDef `json:"response,omitempty"`
	Concurrent       bool                 `json:"concurrent"`
}

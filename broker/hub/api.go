package hub

type FlowStepResponseDef struct {
	Type  string       `json:"type"`
	Value BaseValueDef `json:"value,omitempty"`
}

type ApiDef struct {
	Name             string               `json:"name"`
	Command          string               `json:"command"`
	Private          string               `json:"private"`
	Description      string               `json:"description"`
	ResultKey        string               `json:"resultKey"`
	Parameters       *[]BaseParamDef      `json:"parameters"`
	OriginParameters *[]BaseParamDef      `json:"origin"`
	Response         *FlowStepResponseDef `json:"response,omitempty"`
	Concurrent       bool                 `json:"concurrent"`
}

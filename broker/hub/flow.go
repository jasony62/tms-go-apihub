package hub

type FlowStepApiDef struct {
	Id         string        `json:"id"`
	Parameters []ApiDefParam `json:"parameters"`
}

type FlowStepResponseDef struct {
	Json interface{} `json:"json"`
}

type FlowStepDef struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	ResultKey   string               `json:"resultKey"`
	Api         *FlowStepApiDef      `json:"api,omitempty"`
	Response    *FlowStepResponseDef `json:"response,omitempty"`
}

type FlowDef struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Steps       []FlowStepDef `json:"steps"`
}

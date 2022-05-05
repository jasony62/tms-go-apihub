package hub

type OriginDefParam struct {
	Name  string           `json:"name"`
	Value string           `json:"value,omitempty"`
	From  *ApiDefParamFrom `json:"from,omitempty"`
}

type FlowStepApiDef struct {
	Id         string            `json:"id"`
	Parameters *[]OriginDefParam `json:"parameters"`
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

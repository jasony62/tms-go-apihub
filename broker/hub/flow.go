package hub

type OriginDefParam struct {
	Name  string            `json:"name"`
	Value string            `json:"value,omitempty"`
	From  *BaseDefParamFrom `json:"from,omitempty"`
}

type FlowStepApiDef struct {
	Id          string            `json:"id"`
	Parameters  *[]OriginDefParam `json:"parameters"`
	PrivateName string            `json:"private"`
}

type FlowStepResponseDef struct {
	Type string            `json:"type"`
	From *BaseDefParamFrom `json:"from,omitempty"`
}

type FlowStepDef struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	ResultKey   string               `json:"resultKey"`
	Api         *FlowStepApiDef      `json:"api,omitempty"`
	Response    *FlowStepResponseDef `json:"response,omitempty"`
	Concurrent  bool                 `json:"concurrent"`
}

type FlowDef struct {
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	ConcurrentNum int           `json:"concurrentNum"`
	Steps         []FlowStepDef `json:"steps"`
}

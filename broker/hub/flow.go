package hub

const (
	ORIGIN_SRC_API = iota
	ORIGIN_SRC_RESPONSE
)

type OriginDefParam struct {
	In    string             `json:"in"`
	Name  string             `json:"name"`
	Value string             `json:"value,omitempty"`
	From  *BaseDefParamValue `json:"from,omitempty"`
}

type FlowStepApiDef struct {
	Id          string            `json:"id"`
	Parameters  *[]OriginDefParam `json:"parameters"`
	PrivateName string            `json:"private"`
}

type FlowStepResponseDef struct {
	Type        string             `json:"type"`
	From        *BaseDefParamValue `json:"from,omitempty"`
	PrivateName string             `json:"private,omitempty"`
	Parameters  *[]OriginDefParam  `json:"parameters,omitempty"`
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

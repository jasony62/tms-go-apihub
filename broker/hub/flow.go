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

type FlowDef struct {
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ConcurrentNum int       `json:"concurrentNum"`
	Tasks         []TaskDef `json:"tasks"`
}

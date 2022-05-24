package hub

type OriginDefParam struct {
	Name string             `json:"name"`
	From *BaseDefParamValue `json:"from,omitempty"`
}

type FlowDef struct {
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ConcurrentNum int       `json:"concurrentNum"`
	Tasks         []TaskDef `json:"tasks"`
}

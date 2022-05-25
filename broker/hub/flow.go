package hub

type FlowDef struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Steps       []ApiDef `json:"steps"`
}

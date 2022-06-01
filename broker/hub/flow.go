package hub

type FlowDef struct {
	Name        string   `json:"name"`
	Private     string   `json:"private"`
	Description string   `json:"description"`
	Steps       []ApiDef `json:"steps"`
}

package hub

type FlowDef struct {
	Name         string   `json:"name"`
	Private      string   `json:"private"`
	Description  string   `json:"description"`
	DefaultRight string   `json:"defaultRight"`
	Steps        []ApiDef `json:"steps"`
}

package hub

type RightList struct {
	User string `json:"user"`
}

type RightArray struct {
	Type    string       `json:"type"`
	Right   string       `json:"right"`
	Default string       `json:"default"`
	List    *[]RightList `json:"list"`
}

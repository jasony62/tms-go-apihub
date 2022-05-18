package hub

type PrivatePairs struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PrivateArray struct {
	Pairs *[]PrivatePairs `json:"privates"`
}

package hub

type ApiDefParamFrom struct {
	In   string `json:"in"`
	Name string `json:"name"`
}

type ApiDefParam struct {
	In    string           `json:"in"`
	Name  string           `json:"name"`
	Value string           `json:"value,omitempty"`
	From  *ApiDefParamFrom `json:"from,omitempty"`
}

type ApiDefRequestBody struct {
	ContentType string       `json:"contentType"`
	Content     *interface{} `json:"content"`
}

type ApiDefResponse struct {
	Json interface{} `json:"json"`
}

type ApiDefPlugin struct {
	Path string `json:"path"`
}

type PrivatePairs struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PrivateArray struct {
	Pairs *[]PrivatePairs `json:"privates"`
}

type ApiDef struct {
	Id          string             `json:"id"`
	Url         string             `json:"url"`
	Description string             `json:"description"`
	Method      string             `json:"method"`
	PrivateName string             `json:"private"`
	Parameters  *[]ApiDefParam     `json:"parameters"`
	RequestBody *ApiDefRequestBody `json:"requestBody,omitempty"`
	Response    *ApiDefResponse    `json:"response"`
	Plugins     *[]ApiDefPlugin    `json:"plugins,omitempty"`
	Privates    *PrivateArray
}

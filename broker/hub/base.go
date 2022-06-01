package hub

type BaseValueDef struct {
	From    string       `json:"from"`
	Content string       `json:"content"`
	Json    *interface{} `json:"json"`
	Args    string       `json:"args"`
}

type BaseParamDef struct {
	Name  string       `json:"name"`
	Value BaseValueDef `json:"value,omitempty"`
}
type ApiHandler func(*Stack, map[string]string) (interface{}, int)
type FuncHandler func(params []string) string
type TemplateHandler func(args ...interface{}) string

const OriginName = "origin"
const VarsName = "vars"
const LoopName = "loop"
const ResultName = "result"

const Right_Access = "access"
const Right_Deny = "deny"

const Right_Pulbic = "public"
const Right_Internal = "internal"
const Right_Whitelist = "whitelist"
const Right_Blacklist = "blacklist"

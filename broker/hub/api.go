package hub

import (
	"sync"
	"time"
)

type ApiDefParamFrom struct {
	From     string       `json:"from"`
	Name     string       `json:"name"`
	Template *interface{} `json:"template"`
	Args     string       `json:"args"`
}

type ApiDefParam struct {
	In    string           `json:"in"`
	Name  string           `json:"name"`
	Value string           `json:"value,omitempty"`
	From  *ApiDefParamFrom `json:"from,omitempty"`
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
	Id                 string          `json:"id"`
	Url                string          `json:"url"`
	Description        string          `json:"description"`
	Method             string          `json:"method"`
	PrivateName        string          `json:"private"`
	Parameters         *[]ApiDefParam  `json:"parameters"`
	RequestContentType string          `json:"requestContentType"`
	Response           *ApiDefResponse `json:"response"`
	Privates           *PrivateArray
	Cache              *ApiCache      `json:"cache"`
	RespStatus         *ApiRespStatus `json:"respStatus"`
}

type ApiCache struct {
	From    *ApiDefParamFrom `json:"from,omitempty"`
	Format  string           `json:"format,omitempty"`
	Expires time.Time        //expires time
	Resp    interface{}      //response content
	Locker  sync.RWMutex     //api rw lock
}

type ApiRespStatus struct {
	From     *ApiDefParamFrom `json:"from,omitempty"`
	Format   string           `json:"format,omitempty"`   //number or string
	Expected string           `json:"expected,omitempty"` //expected correct code
}

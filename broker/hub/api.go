package hub

import (
	"sync"
	"time"
)

type ApiDefParam struct {
	In   string            `json:"in"`
	Name string            `json:"name"`
	From BaseDefParamValue `json:"from,omitempty"`
}

type ApiDefResponse struct {
	Json interface{} `json:"json"`
}

type ApiDefPlugin struct {
	Path string `json:"path"`
}

type ApiDef struct {
	Id                 string             `json:"id"`
	Url                string             `json:"url"`
	DynamicUrl         *BaseDefParamValue `json:"dynamicUrl"`
	Description        string             `json:"description"`
	Method             string             `json:"method"`
	PrivateName        string             `json:"private"`
	Parameters         *[]ApiDefParam     `json:"parameters"`
	RequestContentType string             `json:"requestContentType"`
	Response           *ApiDefResponse    `json:"response"`
	Cache              *ApiCache          `json:"cache"`
	RespStatus         *ApiRespStatus     `json:"respStatus"`
}

type ApiRespStatus struct {
	From     *BaseDefParamValue `json:"from,omitempty"`
	Format   string             `json:"format,omitempty"`   //number or string
	Expected string             `json:"expected,omitempty"` //expected correct code
}

type ApiCache struct {
	From    *BaseDefParamValue `json:"from,omitempty"`
	Format  string             `json:"format,omitempty"`
	Expires time.Time          //expires time
	Resp    interface{}        //response content
	Locker  sync.RWMutex       //api rw lock
}

package hub

import (
	"sync"
	"time"
)

type HttpApiDefParam struct {
	In    string       `json:"in"`
	Name  string       `json:"name"`
	Value BaseValueDef `json:"value,omitempty"`
}

type ApiCache struct {
	Expire  *BaseValueDef `json:"expire,omitempty"`
	Format  string        `json:"format,omitempty"`
	Expires time.Time     //expires time
	Resp    interface{}   //response content
	Locker  sync.RWMutex  //api rw lock
}

type HttpApiDef struct {
	Id                 string             `json:"id"`
	Url                string             `json:"url"`
	DynamicUrl         *BaseValueDef      `json:"dynamicUrl"`
	Description        string             `json:"description"`
	PrivateName        string             `json:"private"`
	Method             string             `json:"method"`
	RequestContentType string             `json:"requestContentType"`
	Args               *[]HttpApiDefParam `json:"args"`
	Cache              *ApiCache          `json:"cache"`
}

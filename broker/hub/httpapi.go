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
	Method             string             `json:"method"`
	PrivateName        string             `json:"private"`
	Parameters         *[]HttpApiDefParam `json:"parameters"`
	RequestContentType string             `json:"requestContentType"`
	Cache              *ApiCache          `json:"cache"`
}

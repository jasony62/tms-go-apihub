package hub

import (
	"github.com/gin-gonic/gin"
)

type Stack struct {
	GinContext *gin.Context
	StepResult map[string]interface{}
	RootName   string
	ChildName  string
}

// 从请求参数中获取查询参数
func (stack Stack) Query(name string) string {
	/*默认从请求的查询参数中获得*/
	return stack.GinContext.Query(name)
}

package hub

import (
	"bytes"
	"html/template"

	"github.com/gin-gonic/gin"
)

type Stack struct {
	BucketEnable   bool
	ApiDefPath     string
	PrivateDefPath string
	FlowDefPath    string
	ApiDef         *ApiDef
	FlowDef        *FlowDef
	GinContext     *gin.Context
	StepResult     map[string]interface{}
	CurrentStep    *FlowStepDef
}

// 从请求参数中获取查询参数
func (stack Stack) Query(name string) string {
	/*默认从请求的查询参数中获得*/
	return stack.GinContext.Query(name)
}

// 从执行结果中获取查询参数
func (stack Stack) QueryFromStepResult(name string) string {
	tmpl, _ := template.New("key").Parse(name)
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, stack.StepResult)
	return buf.String()
}

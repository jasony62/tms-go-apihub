package hub

import (
	"bytes"
	"html/template"

	"github.com/gin-gonic/gin"
	klog "k8s.io/klog/v2"
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

// 从执行结果中获取查询参数
func (stack Stack) QueryFromStepResult(name string) string {
	tmpl, err := template.New("key").Funcs(FuncMapForTemplate).Parse(name)
	if err != nil {
		klog.Infoln("QueryFromStepResult 创建并解析template失败:", err)
		return ""
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, stack.StepResult)
	if err != nil {
		klog.Infoln("渲染template失败:", err)
	}
	return buf.String()
}

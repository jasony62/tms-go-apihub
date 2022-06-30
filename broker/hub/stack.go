package hub

import (
	"github.com/gin-gonic/gin"
)

type Stack struct {
	GinContext *gin.Context
	Heap       map[string]interface{}
	Base       map[string]string
}

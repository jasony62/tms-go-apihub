package hub

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Stack struct {
	GinContext *gin.Context
	Heap       map[string]interface{}
	BaseString string
	StartTime  time.Time
}

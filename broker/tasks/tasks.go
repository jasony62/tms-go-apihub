package tasks

import (
	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/task"
)

func Init() {
	task.RegisterTasks(map[string]hub.TaskHandler{"checkStringsEqual": checkStringsEqual, "checkStringsNotEqual": checkStringsNotEqual})
}

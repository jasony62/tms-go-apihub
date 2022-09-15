package apis

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jasony62/tms-go-apihub/hub"
	"github.com/jasony62/tms-go-apihub/logger"
)

func pressureSleep(stack *hub.Stack, params map[string]string) (interface{}, int) {
	hour := params["hour"]
	minute := params["minute"]
	second := params["second"]
	if len(hour) == 0 {
		hour = "0"
	}
	if len(minute) == 0 {
		minute = "1"
	}
	if len(second) == 0 {
		second = "0"
	}
	logger.LogS().Infoln("sleep time interval", hour, " h: ", minute, " m: ", second, " s")
	hourInt, _ := strconv.Atoi(hour)
	minuteInt, _ := strconv.Atoi(minute)
	secondInt, _ := strconv.Atoi(second)
	time.Sleep(time.Duration(hourInt)*time.Hour + time.Duration(minuteInt)*time.Minute + time.Duration(secondInt)*time.Second)
	return nil, http.StatusOK
}

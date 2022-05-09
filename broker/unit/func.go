package unit

import (
	"crypto/md5"
	"fmt"
	"github.com/jasony62/tms-go-apihub/hub"
	"io"
	"strconv"
	"time"
)

func init() {
	hub.FuncMap = map[string](interface{}){
		"utc":      utc,
		"checkSum": checkSum,
	}
}

func utc() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func checkSumNoMinusParam(args ...interface{}) string {

	if len(args) == 0 {
		return ""
	}
	str := fmt.Sprint(args...)
	w := md5.New()
	io.WriteString(w, str)
	checksum := fmt.Sprintf("%x", w.Sum(nil))
	return checksum
}

func checkSum(params []string) string {

	if len(params) == 0 {
		return ""
	}
	var str string
	for _, v := range params {
		str = str + v
	}
	w := md5.New()
	io.WriteString(w, str)
	checksum := fmt.Sprintf("%x", w.Sum(nil))
	return checksum
}

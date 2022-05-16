package unit

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/jasony62/tms-go-apihub/hub"
)

func init() {
	hub.FuncMap = map[string](interface{}){
		"utc": utc,
		"md5": md5Func,
	}
	hub.FuncMapForTemplate = map[string](interface{}){
		"utc": utc,
		"md5": md5Template,
	}
}

func utc() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func md5Template(args ...interface{}) string {

	if len(args) == 0 {
		return ""
	}
	str := fmt.Sprint(args...)
	w := md5.New()
	io.WriteString(w, str)
	checksum := fmt.Sprintf("%x", w.Sum(nil))
	return checksum
}

func md5Func(params []string) string {

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

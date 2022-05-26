package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
)

func utc(params []string) string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func utcTemplate(args ...interface{}) string {
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

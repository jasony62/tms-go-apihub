package main

import (
	"encoding/json"
	"flag"
	"log"
	"strings"

	klog "k8s.io/klog/v2"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

var (
	addr = flag.String("addr", "127.0.0.1:6060", "TCP address to listen to")
)

func main() {
	flag.Parse()

	r := router.New()
	r.GET("/", Echo)
	r.GET("/joint", Joint)
	r.GET("/split", Split)
	// r.GET("/hello/{name}", Hello)v3/config/district
	r.GET("/v3/config/district", AmapDistrict)
	r.GET("/v3/weather/weatherInfo", AmapWeather)

	if err := fasthttp.ListenAndServe(*addr, r.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe: %v", err)
	}
}

func AmapDistrict(ctx *fasthttp.RequestCtx) {
	klog.Infof("%%%%%%%%Connection has been established at %s\n", ctx.ConnTime())
	klog.Infof("%%%%%%%%RequestURI is %q\n", ctx.RequestURI())
	city := string(ctx.Request.URI().QueryArgs().Peek("keywords"))
	klog.Infoln("%%%%%%%%City is ", city)
	// klog.Infoln("%%%%%%%%Header is ", string(ctx.Request.Header.Peek("key")))
	// if strings.Compare(city, "北京") == 0 {
	// 	klog.Infoln("%%%%%%%%City is 北京！ ", city)
	// }

	ctx.SetContentType("application/json; charset=utf-8")
	district := District{
		Citycode: "010",
		Adcode:   "110100",
		Name:     "北京城区",
		Center:   "116.405285,39.904989",
		Level:    "city"}
	reqDistrict := Amap{
		Status:    "1",
		Count:     "1",
		Info:      "OK",
		Infocode:  "10000",
		Districts: []District{district}}
	reqDistrictBytes, _ := json.Marshal(reqDistrict)
	ctx.Response.SetBody(reqDistrictBytes)
}

func AmapWeather(ctx *fasthttp.RequestCtx) {
	klog.Infof("%%%%%%%%Connection has been established at %s\n", ctx.ConnTime())
	ctx.SetContentType("application/json; charset=utf-8")

	klog.Infof("%%%%%%%%RequestURI is %q\n", ctx.RequestURI())
	cityCode := string(ctx.Request.URI().QueryArgs().Peek("city"))
	klog.Infoln("%%%%%%%%cityCode is ", cityCode)
	if strings.Compare(cityCode, "110100") == 0 {
		reqWeatherBytes, _ := json.Marshal(apiMap[cityCode])
		ctx.Response.SetBody(reqWeatherBytes)
	}
	// reqWeatherBytes, _ := json.Marshal(apiMap[cityCode])
	// ctx.Response.SetBody(reqWeatherBytes)
}

func Echo(ctx *fasthttp.RequestCtx) {
	klog.Infof("%%%%%%%%Connection has been established at %s\n", ctx.ConnTime())
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.Request.Header.Cookie()
	klog.Infof("DEBUG Request: %s\n", ctx.Request.Body())
	ctx.Response.SetBody(ctx.Request.Body())
}

func Joint(ctx *fasthttp.RequestCtx) {
	klog.Infof("%%%%%%%%Connection has been established at %s\n", ctx.ConnTime())
	ctx.SetContentType("application/json; charset=utf-8")
	klog.Infof("DEBUG Request: %s\n", ctx.Request.Body())
	reqBody := &Params{}
	json.Unmarshal(ctx.Request.Body(), reqBody)

	content := &Content{
		Content: reqBody.Param1 + reqBody.Param2}
	contentBytes, _ := json.Marshal(content)
	ctx.Response.SetBody(contentBytes)
}

func Split(ctx *fasthttp.RequestCtx) {
	klog.Infof("%%%%%%%%Connection has been established at %s\n", ctx.ConnTime())
	ctx.SetContentType("application/json; charset=utf-8")
	klog.Infof("DEBUG Request: %s\n", ctx.Request.Body())
	reqBody := &Content{}
	json.Unmarshal(ctx.Request.Body(), reqBody)
	reqStr := strings.Fields(reqBody.Content)
	strJoint := strings.Join(reqStr, ",")
	respContent := &Content{
		Content: strJoint}
	reqEntityBytes, _ := json.Marshal(respContent)
	ctx.Response.SetBody(reqEntityBytes)

}

// Set cookies
//var c fasthttp.Cookie
//c.SetKey("cookie-name")
//c.SetValue("cookie-value")
//ctx.Response.Header.SetCookie(&c)

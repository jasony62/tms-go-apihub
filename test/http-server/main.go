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
	r.GET("/", RequestHandler)
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

func RequestHandler(ctx *fasthttp.RequestCtx) {
	klog.Infof("%%%%%%%%Connection has been established at %s\n", ctx.ConnTime())
	ctx.SetContentType("text/plain; charset=utf8")
	// Set arbitrary headers
	ctx.Response.Header.Set("X-My-Header", "my-header-value")
	ctx.Response.Header.SetContentType("application/json")
	reqEntity := &Entity{
		Name: "test",
	}
	reqEntityBytes, _ := json.Marshal(reqEntity)
	ctx.Response.SetBody(reqEntityBytes)

}

// Set cookies
//var c fasthttp.Cookie
//c.SetKey("cookie-name")
//c.SetValue("cookie-value")
//ctx.Response.Header.SetCookie(&c)

package main

import (
	"encoding/json"
	"flag"
	"strings"
	"time"

	klog "k8s.io/klog/v2"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

var (
	addr = flag.String("addr", "127.0.0.1:6060", "TCP address to listen to")
	conf = flag.String("f", "./conf", "config file name")
)

func main() {
	flag.Parse()
	readConfig(*conf, &userInfo)
	hRouter := createRouter()

	klog.Infoln("Listen and Serve addr:", *addr)
	if err := fasthttp.ListenAndServe(*addr, hRouter.Handler); err != nil {
		klog.Fatal("Error in ListenAndServe: %v", err)
	}
}

func createRouter() *router.Router {
	r := router.New()
	r.GET("/echo", Echo)
	r.GET("/joint", Joint)
	r.GET("/split", Split)
	r.GET("/register", Register)
	// r.GET("/hello/{name}", Hello)v3/config/district
	r.GET("/v3/config/district", AmapDistrict)
	r.GET("/v3/weather/weatherInfo", AmapWeather)
	return r
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

func checkAccountAndTokenValid(ctx *fasthttp.RequestCtx, userInfo UserInfo) bool {
	usr := AppUser(ctx.Request.URI().QueryArgs().Peek("app"))
	userVerif, ok := userInfo.UsrSecurity[usr]
	if !ok || userVerif.Pwd == "" {
		klog.Warningln("This app is NOT recognized!")
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		return false
	}
	reqToken := string(ctx.Request.Header.Peek("Authorization"))
	if reqToken != userVerif.Token {
		klog.Warningln("This app is Unauthorized!")
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		return false
	}
	return true
}

func Echo(ctx *fasthttp.RequestCtx) {
	klog.Infof("%%%%%%%%Connection has been established at %s\n", ctx.ConnTime())
	klog.Infof("DEBUG Request: %s\n", ctx.Request.Body())
	ctx.SetContentType("application/json; charset=utf-8")

	//检查该usr是否已经配置了key,并验证token
	if !checkAccountAndTokenValid(ctx, userInfo) {
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}

	ctx.Response.SetBody(ctx.Request.Body())
}

func Joint(ctx *fasthttp.RequestCtx) {
	klog.Infof("%%%%%%%%Connection has been established at %s\n", ctx.ConnTime())
	ctx.SetContentType("application/json; charset=utf-8")
	klog.Infof("DEBUG Request: %s\n", ctx.Request.Body())

	//检查该usr是否已经配置了key,并验证token
	if !checkAccountAndTokenValid(ctx, userInfo) {
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}
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
	//检查该usr是否已经配置了key,并验证token
	if !checkAccountAndTokenValid(ctx, userInfo) {
		ctx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}
	reqBody := &Content{}
	json.Unmarshal(ctx.Request.Body(), reqBody)
	reqStr := strings.Fields(reqBody.Content)
	strJoint := strings.Join(reqStr, ",")
	respContent := &Content{
		Content: strJoint}
	reqEntityBytes, _ := json.Marshal(respContent)
	ctx.Response.SetBody(reqEntityBytes)

}

func Register(ctx *fasthttp.RequestCtx) {
	klog.Infof("%%%%%%%%Connection has been established at %s\n", ctx.ConnTime())
	ctx.SetContentType("application/json; charset=utf-8")
	klog.Infof("DEBUG Request: %s\n", ctx.Request.Body())
	reqBody := &RegistEntity{}
	json.Unmarshal(ctx.Request.Body(), reqBody)
	usr := AppUser(ctx.Request.URI().QueryArgs().Peek("app"))
	curTime := time.Now().Unix()
	//检查该usr是否已经配置了key
	userVerif, ok := userInfo.UsrSecurity[usr]
	if !ok || userVerif.Pwd == "" {
		klog.Warningln("This app is NOT recognized!")
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	//检查utc是否超时
	if !checkUtcTimeValid(curTime, reqBody.Utc) {
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	//checkSum校验
	if !checkSumValid(userVerif.Pwd, reqBody) {
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	//若已有的token的未失效
	if userVerif.ExpireUtc >= curTime {
		expire := userVerif.ExpireUtc - curTime
		respContent := &RegistResp{
			Token:   userVerif.Token,
			Expires: expire}
		reqEntityBytes, _ := json.Marshal(respContent)
		ctx.Response.SetBody(reqEntityBytes)
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		klog.Infoln("%%%%%%%%Exist valid Token: ", *respContent)
		return
	}

	//注册成功,更新当前用户的过期时间和token
	userVerif.ExpireUtc = curTime + userVerif.Expires
	userVerif.Token = generateToken(ctx.Path(), ctx.Method(), userVerif.Expires)
	respContent := &RegistResp{
		Token:   userVerif.Token,
		Expires: userVerif.Expires}
	reqEntityBytes, _ := json.Marshal(respContent)
	ctx.Response.SetBody(reqEntityBytes)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	klog.Infoln("%%%%%%%%Register succeed!")
}

// Set cookies
//var c fasthttp.Cookie
//c.SetKey("cookie-name")
//c.SetValue("cookie-value")
//ctx.Response.Header.SetCookie(&c)

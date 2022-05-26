package main

type AppUser string
type Entity struct {
	Id   int
	Name string
}
type UserInfo struct {
	UsrSecurity map[AppUser]*Verification `json:"usrSecurity"`
}
type Verification struct {
	Pwd       string `json:"pwd"`
	Token     string `json:"token"`
	Expires   int64  `json:"expires"`
	ExpireUtc int64  `json:"expireUtc"`
}

type RegistEntity struct {
	Nonce    string `json:"nonce"`
	Utc      string `json:"utc"`
	Checksum string `json:"checksum"`
}

type RegistResp struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

type Content struct {
	Content string `json:"content"`
}

type Params struct {
	Param1 string `json:"param1"`
	Param2 string `json:"param2"`
}

type Amap struct {
	Status     string     `json:"status"`
	Count      string     `json:"count"`
	Info       string     `json:"info"`
	Infocode   string     `json:"infocode"`
	Suggestion Suggestion `json:"suggestion"`
	Districts  []District `json:"districts"`
	Lives      []Lives    `json:"lives"`
}

type Suggestion struct {
	Keywords []interface{} `json:"keywords"`
	Cities   []interface{} `json:"cities"`
}
type District struct {
	Citycode  string        `json:"citycode"`
	Adcode    string        `json:"adcode"`
	Name      string        `json:"name"`
	Center    string        `json:"center"`
	Level     string        `json:"level"`
	Districts []interface{} `json:"districts"`
}
type Lives struct {
	Province      string `json:"province"`
	City          string `json:"city"`
	Adcode        string `json:"adcode"`
	Weather       string `json:"weather"`
	Temperature   string `json:"temperature"`
	Winddirection string `json:"winddirection"`
	Windpower     string `json:"windpower"`
	Humidity      string `json:"humidity"`
	Reporttime    string `json:"reporttime"`
}

var reqWeather Amap = Amap{
	Status:   "1",
	Count:    "1",
	Info:     "OK",
	Infocode: "10000",
	Lives: []Lives{
		Lives{
			Province:      "北京",
			City:          "北京城区",
			Adcode:        "110100",
			Weather:       "晴",
			Temperature:   "29",
			Winddirection: "南",
			Windpower:     "4",
			Humidity:      "31",
			Reporttime:    "2022-05-19 15:41:42"}}}

var apiMap map[string](interface{}) = map[string](interface{}){
	"110100": reqWeather,
}

var userInfo UserInfo

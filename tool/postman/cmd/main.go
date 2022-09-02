// /////////////////////////////////////////////////////////////
//
//	                                                          //
//	                    _ooOoo_                               //
//	                   o8888888o                              //
//	                   88" . "88                              //
//	                   (| ^_^ |)                              //
//	                   O\  =  /O                              //
//	                ____/`---'\____                           //
//	              .'  \\|     |//  `.                         //
//	             /  \\|||  :  |||//  \                        //
//	            /  _||||| -:- |||||-  \                       //
//	            |   | \\\  -  /// |   |                       //
//	            | \_|  ''\---/''  |   |                       //
//	            \  .-\__  `-`  ___/-. /                       //
//	          ___`. .'  /--.--\  `. . ___                     //
//	        ."" '<  `.___\_<|>_/___.'  >'"".                  //
//	      | | :  `- \`.;`\ _ /`;.`/ - ` : | |                 //
//	      \  \ `-.   \_ __\ /__ _/   .-` /  /                 //
//	========`-.____`-.___\_____/___.-`____.-'========         //
//	                     `=---='                              //
//	^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^        //
//	         佛祖保佑       永不宕机     永无BUG                //
//
// /////////////////////////////////////////////////////////////
package main

import (
	"flag"

	postmaninternal "github.com/Sheng-ZM/tms-go-apihub/tree/shengzm-local/internal"
)

// 初始化
func init() {
	flag.StringVar(&postmaninternal.PostmanPath, "from", "./cmd/postman_collections/", "指定postman_collections文件路径")
	flag.StringVar(&postmaninternal.ApiHubJsonPath, "to", "./cmd/httpapis/", "指定转换后的apiHub json文件路径")
	flag.StringVar(&postmaninternal.ApiHubPrivatesJsonPath, "private", "./cmd/privates/", "指定转换后的apiHub privates json文件路径")
}

/*********************main主程序*****************************/
func main() {
	postmaninternal.ConvertPostmanFiles(postmaninternal.PostmanPath)
}

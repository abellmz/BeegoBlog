package routers

import (
	"BeegoBlog/AllSrcCode/controllers"
	"BeegoBlog/AllSrcCode/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	//beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.AdminController{}, "get:WebSocketHandler")
	beego.Router("/users", &controllers.AdminController{}, "get:GetUsers")
	beego.Router("/login", &controllers.AuthController{}, "get,post:Login")
	//自动路由，比如Login()方法直接转换为admin/login
	//beego.Controller是ControllerInterface下众多方法的接收者，是ControllerInterface的实现者
	beego.AutoRouter(&controllers.AdminController{})
	beego.AutoRouter(&controllers.ApiController{})
	beego.ErrorController(&controllers.ErrorController{})
	/**
	插入一个过滤器（middleware）
	"/* :" 匹配/*的路由，
	beego.BeforeRouter 中间件在路由处理前执行，
	beego.BeforeExec 在路由正确匹配后执行
	JWTAuthMiddleware 中间件函数本身，用于实现自定义逻辑
	*/
	beego.InsertFilter("/admin/*", beego.BeforeExec, JWTAuthMiddleware)

}

func JWTAuthMiddleware(context *context.Context) {
	util.JWTAuthMiddleware(context)
}

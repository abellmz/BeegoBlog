package routers

import (
	"BeegoBlog/AllSrcCode/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	//beego.Router("/ws", &controllers.AdminController{}, "get:WebSocketHandler")
	//beego.Router("/users", &controllers.AdminController{}, "get:GetUsers")
	//自动路由，比如Login()方法直接转换为admin/login
	//beego.Controller是ControllerInterface下众多方法的接收者，是ControllerInterface的实现者
	beego.AutoRouter(&controllers.AdminController{})
	beego.AutoRouter(&controllers.ApiController{})
}

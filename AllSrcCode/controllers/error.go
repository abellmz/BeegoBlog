package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

/*
*
命名规则是固定的:Error+错误的状态码,
若想要明确的状态码可以手动设置：c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
*/
func (c *ErrorController) Error404() {
	c.Data["content"] = "page not found"
	c.TplName = "admin/404.html"
}

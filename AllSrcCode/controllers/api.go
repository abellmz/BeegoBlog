package controllers

import (
	"BeegoBlog/AllSrcCode/util"
	"github.com/astaxie/beego/logs"
)

type ApiController struct {
	baseController
}

func (c *ApiController) RedisTest() {
	val, err := util.TTL("mykey")
	if err != nil {
		logs.Error("Error getting cache:", err)
		return
	}
	c.Data["json"] = map[string]interface{}{
		"code":    "200",
		"message": "Redis test",
		"data":    val,
	}
	c.ServeJSON()
	//if c.Ctx.Request.Method == "POST" {}
}

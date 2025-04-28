package controllers

import (
	"BeegoBlog/AllSrcCode/models"
	"BeegoBlog/AllSrcCode/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"net/http"
	"strings"
	"time"
)

type AuthController struct {
	beego.Controller
	o orm.Ormer
}

func (c *AuthController) Login() {

	if c.Ctx.Request.Method == "POST" {
		//从前台的请求中获取账号密码信息
		username := c.GetString("username")
		password := c.GetString("password")
		// 获取 ORM 实例
		c.o = orm.NewOrm()
		// 根据前面的信息从数据库中查出 username 的数据
		user := models.User{Username: username}
		ormerr := c.o.Read(&user, "username")
		if ormerr != nil {
			logs.Error("orm错误：", ormerr)
		}

		if user.Password == "" {
			c.History("账号不存在", "")
		}

		//判断登录信息是否正确
		if util.Md5(password) != strings.Trim(user.Password, " ") {
			c.History("密码错误", "")
		}

		//如果账号密码无误则更新admin信息
		user.LastIp = c.getClientIp()
		user.LoginCount = user.LoginCount + 1

		if _, err := c.o.Update(&user); err != nil {
			c.History("登录异常", "")
		} else {
			// 生成token
			token, err := util.GenerateToken(user.Id, user.Username, "123456")
			if err != nil {
				logs.Error("生成 Token 失败", err)
				c.History("登录异常", "")
			}
			// 将 token 存储在 Cookie 中
			http.SetCookie(c.Ctx.ResponseWriter, &http.Cookie{
				Name:     "token",
				Value:    token,
				Expires:  time.Now().Add(24 * time.Hour),
				Path:     "/",  // Cookie 的作用路径,这里表示对网站的所有路径都有效
				HttpOnly: true, //标记为HTTP Only
				Secure:   true, // 确保通过 HTTPS 传输
			})
			c.Ctx.SetCookie("token", token, 3600, "/", "", true, true)
			//-----
			// 将 token 存储在 HTTP 头部
			//c.Ctx.ResponseWriter.Header().Set("X-Auth-Token", token)
			logs.Info("生成 Token 成功!")
			c.History("登录成功", "/admin/main.html")
		}
		c.SetSession("user", user)
	}
	//注册登录模板
	c.TplName = "admin/login.html"
}

func (p *AuthController) History(msg string, url string) {
	if url == "" {
		p.Ctx.WriteString("<script>alert('" + msg + "');window.history.go(-1);</script>")
		p.StopRun()
	} else {
		p.Redirect(url, 302)
	}
}

// 获取用户IP地址
func (p *AuthController) getClientIp() string {
	s := strings.Split(p.Ctx.Request.RemoteAddr, ":")
	return s[0]
}

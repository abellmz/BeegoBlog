package controllers

import (
	"BeegoBlog/AllSrcCode/models"
	"BeegoBlog/AllSrcCode/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type AuthController struct {
	beego.Controller
	o orm.Ormer
}

func (c *AuthController) Login() {

	if c.Ctx.Request.Method == "POST" {
		// Beego 会自动验证 CSRF 令牌
		//从前台的请求中获取账号密码信息
		username := c.GetString("username")
		password := c.GetString("password")
		// 清理输入 防XSS
		cleanedUsername := cleanInput(username)
		cleanedPassword := cleanInput(password)
		//格式检查
		isValid := isValidInput(cleanedUsername, cleanedPassword)
		if !isValid {
			logs.Error("账号或密码格式错误！")
			c.History("账号或密码格式错误", "")
		}

		// 获取 ORM 实例
		c.o = orm.NewOrm()
		// 根据前面的信息从数据库中查出 username 的数据
		user := models.User{Username: cleanedUsername}
		ormerr := c.o.Read(&user, "username")
		if ormerr != nil {
			logs.Error("orm错误：", ormerr)
			return
		}

		if user.Password == "" {
			logs.Error("账号不存在：", "")
			c.History("账号不存在", "")
		}

		//判断登录信息是否正确
		if util.Md5(cleanedPassword) != strings.Trim(user.Password, " ") {
			c.History("密码错误", "")
		}

		//如果账号密码无误则更新admin信息
		user.LastIp = c.getClientIp()
		user.LoginCount = user.LoginCount + 1

		if _, err := c.o.Update(&user); err != nil {
			c.History("登录异常", "")
		} else {
			jwtSecretKey := util.GetJWTSecretKey()
			// 生成token
			token, tokenErr := util.GenerateToken(user.Id, user.Username, jwtSecretKey)
			if tokenErr != nil {
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
	// 获取 CSRF 令牌
	csrfToken := c.XSRFToken()
	// 将 CSRF 令牌传递给模板
	c.Data["xsrfdata"] = csrfToken
	//注册登录模板
	c.TplName = "admin/login.html"
}

// 清理输入，移除潜在的恶意代码
func cleanInput(input string) string {
	// 移除HTML标签
	noHTML := regexp.MustCompile(`</?[^>]+>`).ReplaceAllString(input, "")
	// 转义特殊字符
	escaped := template.HTMLEscapeString(noHTML)
	// 移除多余的空白字符
	cleaned := strings.TrimSpace(escaped)
	return cleaned
}

// 正则验证输入是否符合格式
func isValidInput(username, password string) bool {
	//5-20个字符，只能包含字母、数字和下划线
	usernamePattern := regexp.MustCompile(`^[a-zA-Z0-9_]{5,20}$`)
	passwordPattern := regexp.MustCompile(`^[a-zA-Z0-9_]{6,20}$`)

	return usernamePattern.MatchString(username) && passwordPattern.MatchString(password)
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

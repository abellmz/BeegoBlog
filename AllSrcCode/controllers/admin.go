package controllers

import (
	"BeegoBlog/AllSrcCode/models"
	"BeegoBlog/AllSrcCode/util"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	//"syscall"
	"time"
)

var (
	// 连接池
	socketClients = make(map[*websocket.Conn]string, 10)
	// 锁
	//mu sync.Mutex
	sendToUserMu       sync.Mutex
	socketClientsMu1   sync.Mutex
	socketClientsMu2   sync.Mutex
	broadcastMessageMu sync.Mutex
	//getUsersMu sync.Mutex
	// 消息通道
	messageChannel = make(chan Message, 20)
)

// WebSocket升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //直接返回 true 是为了允许所有来源的连接，但在实际项目里，最好还是加点限制，别让坏人钻了空子。如何防范呢？
	},
}

type AdminController struct {
	baseController
}

// 消息结构体
type Message struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Text       string `json:"text"`
}

func init() {
	// 启动消息处理协程
	go handleMessages()
}

// WebSocketHandler 处理 WebSocket 连接
func (c *AdminController) WebSocketHandler() {

	//返回 升级器升级后的连接和错误
	ws, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		logs.Error("Upgrade error:", err)
		return
	}
	defer ws.Close() //连接关闭

	//获取用户ID，token解析
	// 从 Cookie 中获取 token
	token := c.Ctx.GetCookie("token")
	if token == "" {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		logs.Error("Unauthorized:")
		return
	}
	claims, err := util.ParseToken(token, "123456")
	if err != nil {
		logs.Error("Invalid token:", err, claims)
		return
	}

	// 添加到连接池
	socketClientsMu1.Lock()
	//一个map 存储连接-userID
	socketClients[ws] = claims.UserID
	socketClientsMu1.Unlock()
	// 连接中 发送欢迎消息
	err = ws.WriteMessage(websocket.TextMessage, []byte("Welcome to the chat room!"))
	if err != nil {
		logs.Error("Write error:", err)
		return
	}
	for {
		_, msg, readMessageErr := ws.ReadMessage()
		if readMessageErr != nil { //若在此处打断点，在websocket未关闭前，刷新聊天室，会导致socket.CloseError
			logs.Error("Read error"+claims.UserID+":", readMessageErr)
			break
		}
		var message Message
		if unmarshalErr := json.Unmarshal(msg, &message); unmarshalErr != nil {
			logs.Error("Unmarshal error:", unmarshalErr)
			continue
		}
		message.SenderID = claims.UserID
		messageChannel <- message
	}
	//移除池中连接-userid
	socketClientsMu2.Lock()
	delete(socketClients, ws)
	socketClientsMu2.Unlock()
	//redisErr := util.Del(claims.UserID)
	//if redisErr != nil {
	//	logs.Error("Redis error:",err)
	//	return
	//}
}

// handleMessages 处理消息
func handleMessages() {
	for msg := range messageChannel {
		if msg.ReceiverID == "all" {
			broadcastMessage(msg)
		} else {
			sendToUser(msg)
		}
	}
}

// broadcastMessage 广播消息
func broadcastMessage(msg Message) {
	broadcastMessageMu.Lock()
	defer broadcastMessageMu.Unlock()
	for client := range socketClients {
		err := client.WriteMessage(websocket.TextMessage, []byte(msg.SenderID+":"+msg.Text))
		if err != nil {
			logs.Error("Broadcast error:", err)
			client.Close() //关闭连接，移除池中连接-userid
			delete(socketClients, client)
		}
	}
}

// sendToUser 单独发送消息
func sendToUser(msg Message) {
	sendToUserMu.Lock()
	defer sendToUserMu.Unlock()
	for client, clientID := range socketClients {
		if clientID == msg.ReceiverID {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg.SenderID+":"+msg.Text))
			if err != nil {
				logs.Error("Send error:", err)
				client.Close() //关闭连接，移除池中连接-userid
				delete(socketClients, client)
			}
			break
		}
	}
}

// GetUsers 获取在线用户列表
func (c *AdminController) GetUsers() {
	//长度为0，容量len(socketClients)
	users := make([]string, 0, len(socketClients))
	//getUsersMu.Lock()
	for _, userID := range socketClients {
		users = append(users, userID)
	}
	//getUsersMu.Unlock()
	c.Data["json"] = users
	c.ServeJSON()
}
func (c *AdminController) Chat() {
	//获取用户ID，token解析
	// 从 Cookie 中获取 token
	token := c.Ctx.GetCookie("token")
	if token == "" {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		logs.Error("Unauthorized")
		return
	}
	claims, err := util.ParseToken(token, "123456")
	if err != nil {
		logs.Error("Invalid token:", err, claims)
		return
	}
	c.Data["userID"] = claims.UserID
	c.TplName = c.controllerName + "/chat.html"

	// 根据前面的信息从数据库中查出 username 的数据
	//user := models.User{Username: username}
	//c.o.Read(&user, "username")
}

/********************************/

// 后台用户登录
func (c *AdminController) Login() {

	if c.Ctx.Request.Method == "POST" {
		//从前台的请求中获取账号密码信息
		username := c.GetString("username")
		password := c.GetString("password")

		// 根据前面的信息从数据库中查出 username 的数据
		user := models.User{Username: username}
		c.o.Read(&user, "username")

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
			logs.Info("生成 Token 成功！")
			c.History("登录成功", "/admin/main.html")
		}
		c.SetSession("user", user)
	}

	//注册登录模板
	c.TplName = c.controllerName + "/login.html"

}

// 退出登录
func (c *AdminController) Logout() {
	c.DestroySession()
	c.History("退出登录", "/admin/login.html")
}

// 配置后台
func (c *AdminController) Config() {
	//读取后台信息，后台信息由ID,name,value三列组成
	//key-value对可以是QQ, URL, 时区等内容
	var result []*models.Config
	c.o.QueryTable(new(models.Config).TableName()).All(&result)

	//建两个map
	//option用于向模板中传递参数，和数据库中的表示方法一样：name-value，所以是string-string
	//mp则用与检查，表示方法是：name-Config{}结构体
	options := make(map[string]string)
	mp := make(map[string]*models.Config)
	for _, v := range result {
		options[v.Name] = v.Value
		mp[v.Name] = v //v是一个Config结构体
	}
	if c.Ctx.Request.Method == "POST" {
		keys := []string{"url", "title", "keywords", "description", "email", "start", "qq"}
		for _, key := range keys {
			val := c.GetString(key) //从我们post的请求中读取参数
			//先检查一下数据库中有没有我们读取到的参数
			if _, ok := mp[key]; !ok { //没有
				c.o.Insert(&models.Config{Name: key, Value: val})
			} else { //如果已经有了，则需要更新
				opt := mp[key]
				if _, err := c.o.Update(&models.Config{Id: opt.Id, Name: opt.Name, Value: val /*此处使用POST数据*/}); err != nil {
					panic(err)
				}
			}
		}
		c.History("设置数据成功", "")
	}
	c.Data["config"] = options
	c.TplName = c.controllerName + "/config.html"
}

// 后台管理的main页面
func (c *AdminController) Main() { //不要忘了注册后台主页面
	c.TplName = c.controllerName + "/main.tpl"
}

/*************Category*******************/

// 后台管理的Category页面
func (c *AdminController) Category() {
	//categories := []*models.Category{} 不要这么声明切片
	//var categories [] *models.Category 这样也不好
	categories := make([]*models.Category, 0)
	c.o.QueryTable(new(models.Category).TableName()).All(&categories)
	c.Data["categorys"] = categories //虽然categorys拼写错误，但是最好不要修改，因为模板里面用的这个名字。:)
	c.TplName = c.controllerName + "/category.tpl"
}

// 添加category
func (c *AdminController) Categoryadd() {
	//一般来说可以采用c.GetInt等方法，但如果不是Int64则需要采用下面的办法
	id := c.GetString("id")
	//id:=c.Input().Get("id")
	if id != "" { //id为空表示是新添加，不为空则表示修改
		intId, _ := strconv.Atoi(id)
		cate := models.Category{Id: intId}
		c.o.Read(&cate)
		c.Data["cate"] = cate //引导
	}
	c.TplName = c.controllerName + "/category_add.tpl"
}

// 类别的插入和更新
func (c *AdminController) CategorySave() {
	name := c.Input().Get("name") //c.GetString也可以
	id := c.Input().Get("id")
	category := models.Category{
		Name:    name,
		Created: time.Now(),
		Updated: time.Now(),
	}
	if id == "" { //插入
		if _, err := c.o.Insert(&category); err != nil {
			c.History("插入数据错误", "")
		} else {
			c.History("插入数据成功", "/admin/category.html")
		}
	} else { //更新
		intId, err := strconv.Atoi(id)
		if err != nil {
			c.History("参数错误", "")
		}
		category.Id = intId
		if _, err := c.o.Update(&category); err != nil {
			c.History("更新出错", "")
		} else {
			c.History("插入数据成功", "/admin/category.html")
		}
	}
}

// 类别的删除
func (c *AdminController) CategoryDel() {
	strId := c.Input().Get("id")
	intId, err := strconv.Atoi(strId)
	if err != nil {
		c.History("参数错误", "")
	} else {
		if _, err := c.o.Delete(&models.Category{Id: intId}); err != nil {
			c.History("未能成功删除", "")
		} else {
			c.History("删除成功", "/admin/category.html")
		}
	}
}

/************Post*********************/
//博客列表
func (c *AdminController) Index() {
	categories := make([]*models.Category, 0)
	c.o.QueryTable(new(models.Category).TableName()).All(&categories)
	c.Data["categorys"] = categories

	var (
		page     int
		pagesize int = 8
		offset   int
		list     []*models.Post
		keyword  string
		cateId   int
	)
	keyword = c.GetString("title")
	cateId, _ = c.GetInt("cate_id")
	//第一次访问时，page未初始化，为0
	if page, _ = c.GetInt("page"); page < 1 {
		page = 1
	}
	offset = (page - 1) * pagesize
	q := c.o.QueryTable(new(models.Post).TableName())
	if keyword != "" {
		q = q.Filter("title__contains", keyword) //包含关键词的title
	}
	count, _ := q.Count()
	if count > 0 {
		//第一个参数是指定获取几条数据，第二个参数指定从哪里获取start
		//如果页码是1，则访问8篇，从0开始
		//如果页码是2，则访问8篇，从9开始
		q.OrderBy("-created").Limit(pagesize, offset).All(&list)
	}
	//查询过滤数据完毕，向模板写入
	c.Data["keyword"] = keyword
	c.Data["count"] = count
	c.Data["list"] = list
	c.Data["cate_id"] = cateId
	path := fmt.Sprintf("/admin/chat.html?keyword=%s", keyword)
	c.Data["pagebar"] = util.NewPager(page, int(count), pagesize, path, true).ToString()
	c.TplName = c.controllerName + "/list.tpl"
}

// 删除博文
func (c *AdminController) Delete() {
	id, err := strconv.Atoi(c.Input().Get("id"))
	if err != nil {
		c.History("参数错误", "")
	} else {
		if _, err := c.o.Delete(&models.Post{Id: id}); err != nil {
			c.History("未能成功删除", "")
		} else {
			c.History("删除成功", "/admin/chat.html")
		}
	}
}

// 添加博文
func (c *AdminController) Article() {
	categorys := make([]*models.Category, 0)
	c.o.QueryTable(new(models.Category).TableName()).All(&categorys)
	id, _ := c.GetInt("id")
	if id != 0 {
		post := models.Post{Id: id}
		c.o.Read(&post)
		c.Data["post"] = post
	}
	c.Data["categorys"] = categorys
	c.TplName = c.controllerName + "/_form.tpl"
}

// 更新博文
func (c *AdminController) Save() {
	//formData := c.Ctx.Request.Form
	//for key, value := range formData {
	//	fmt.Println(key, ":", value)
	//}//获取表单信息
	post := models.Post{}
	post.UserId = 1
	//post.Id = c.Input().Get("Id")
	post.Title = c.Input().Get("title")
	post.Content = c.Input().Get("content")
	post.IsTop, _ = c.GetInt8("is_top")
	post.Types, _ = c.GetInt8("types")
	post.Tags = c.Input().Get("tags")
	post.Url = c.Input().Get("url")
	post.CategoryId, _ = c.GetInt("cate_id")
	post.Info = c.Input().Get("info")
	post.Image = c.Input().Get("Image")
	post.Created = time.Now()
	post.Updated = time.Now()

	id, _ := c.GetInt("Id")
	if id == 0 {
		if _, err := c.o.Insert(&post); err != nil {
			c.History("插入数据错误"+err.Error(), "")
		} else {
			c.History("插入数据成功", "/admin/chat.html")
		}
	} else {
		post.Id = id
		if _, err := c.o.Update(&post); err != nil {
			c.History("更新数据出错"+err.Error(), "")
		} else {
			c.History("插入数据成功", "/admin/chat.html")
		}
	}
}

/************Post*********************/

// upload images
func (c *AdminController) Upload() {
	// 获取上传的文件
	file, header, err := c.GetFile("uploadname")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    1,
			"message": "获取文件失败",
		}
		c.ServeJSON()
		return
	}
	// 设置上传目录
	uploadDir := "static/uploads"
	// 创建目录（如果不存在）
	err = os.MkdirAll(uploadDir, 0777)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    1,
			"message": "创建目录失败",
		}
		c.ServeJSON()
		return
	}

	// 生成保存的文件名
	filename := filepath.Join(uploadDir, header.Filename)

	// 保存文件
	dst, err := os.Create(filename)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    1,
			"message": "创建文件失败",
		}
		c.ServeJSON()
		return
	}
	defer dst.Close()

	// 复制文件内容
	_, err = io.Copy(dst, file)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code":    1,
			"message": "保存文件失败",
		}
		c.ServeJSON()
		return
	}

	// 返回 JSON 响应
	response := map[string]interface{}{
		"code":    0,
		"message": filepath.Join("/", uploadDir, header.Filename),
	}
	c.Data["json"] = response
	c.ServeJSON()
}

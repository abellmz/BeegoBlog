package models

import (
	"BeegoBlog/AllSrcCode/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

func Init() {

	//从配置文件中获取连接参数
	host := beego.AppConfig.String("mysqlHost")
	port := beego.AppConfig.String("mysqlPort")
	user := beego.AppConfig.String("mysqlUser")
	password := beego.AppConfig.String("mysqlPassword")
	dbName := beego.AppConfig.String("mysqlDb")
	//Data Source Name=DSN
	connectInfo := []string{user, ":", password, "@tcp(", host, ":", port, ")/", dbName, "?charset=utf8"}
	DSN := strings.Join(connectInfo, "")
	err := orm.RegisterDataBase("default", "mysql", DSN)
	if err != nil {
		logs.Error("Failed to connect to mysql: %v", err)
		return
	}
	logs.Info("Mysql connected successfully!")
	orm.RegisterModel(new(User), new(Config), new(Category), new(Post), new(Role), new(Permission))
}

func GetTableName(str string) string {
	return beego.AppConfig.String("mysqlPrefix") + str
}

/*
*
初始化管理员用户
*/
func InitAdmin() {
	o := orm.NewOrm()

	// 检查管理员用户是否已存在
	exists := o.QueryTable(new(User)).Filter("Username", "admin").Exist()
	if !exists {
		ID := util.GenerateSnowflakeId()
		// 对密码进行加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Lmuyangren@123"), bcrypt.DefaultCost)
		if err != nil {
			logs.Error("encryption error:", err)
		}

		// 如果不存在，则插入管理员用户
		user := User{
			Id:       ID,
			Username: "admin",
			//Password: util.Md5("123456"),
			Password: string(hashedPassword),
			LastTime: time.Now(),
			Created:  time.Now(),
			Updated:  time.Now(),
		}
		_, InsertErr := o.Insert(&user)
		if InsertErr != nil {
			// 处理错误
			logs.Error("Failed to insert admin user:", InsertErr)
		} else {
			logs.Info("Admin user initialized successfully")
		}
	}
}

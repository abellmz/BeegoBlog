package main

import (
	"BeegoBlog/AllSrcCode/models"
	_ "BeegoBlog/AllSrcCode/routers"
	"BeegoBlog/AllSrcCode/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

func main() {
	//workPath, _ := os.Getwd()
	//fmt.Println(workPath)
	//appConfigPath := filepath.Join(workPath, "conf", "app.conf")
	//fmt.Println(appConfigPath)

	beego.LoadAppConfig("ini", "conf/app.conf")
	util.InitLogs() //初始化日志
	//beego.LoadAppConfig("ini", "conf/app.conf")
	// 打印加载的配置项（用于调试）
	//fmt.Println("RunMode:", beego.AppConfig.String("runmode"))
	//fmt.Println("HttpPort:", beego.AppConfig.String("httpport"))
	//fmt.Println("MysqlHost:", beego.AppConfig.String("mysqlHost"))
	//fmt.Println("MysqlPort:", beego.AppConfig.String("mysqlPort"))
	//fmt.Println("MysqlUser:", beego.AppConfig.String("mysqlUser"))
	//fmt.Println("MysqlPassword:", beego.AppConfig.String("mysqlPassword"))
	//fmt.Println("MysqlDb:", beego.AppConfig.String("mysqlDb"))
	//beego.BConfig.RunMode = "dev"
	models.Init() // 初始化模型和数据库连接
	//aliasName数据库别名,isDrop是否在创建表之前先删除现有的表,isForce是否强制创建表
	orm.RunSyncdb("default", false, true) // 自动建表 有这个debug直接报错并停止运行
	models.InitAdmin()                    //初始化管理员数据
	util.InitRedis()                      //初始化redis
	err := util.Set("mykey", "bhu", 10*time.Second)
	if err != nil {
		logs.Error("Error getting cache:", err)
		return
	}
	logs.Info("")
	beego.Run()

	//fmt.Println("----------")
	//fmt.Println(val)
	//fmt.Println("----------2")
}

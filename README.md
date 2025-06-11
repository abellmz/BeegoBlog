# BeegoBlog

#### 介绍
来源于Github用户JiangRRRen 改编的lightTrace 的开源博客项目，选取其中有关于后台管理的部分。
项目本身适合于新手练习使用，是基于Beego框架开发的一个Blog后台系统，包含了基本的博客增删改查、图片上传等。

原项目功能：包括用户登录登出、博客栏目和文章的增删改查、文件上传等。
新增功能：

1.新增监控日志功能，增强系统可维护性；

2.新增JWT，提高系统安全性；

3.新增聊天室功能，提升用户互动体验；

4.添加相关安全措施，进一步保障用户数据安全。

（功能的添加为实验性质，不代表一定适配系统）

#### 软件架构


> AllSrcCode 代码所在
>
> > conf 配置项目录
> >
> > > app.conf 配置文件
> >
> > controllers 控制器目录
> >
> > logs 日志目录
> >
> > models 数据表的目录
> >
> > routers 路由目录
> >
> > static 静态文件目录
> >
> > tests 测试目录
> >
> > util 辅助文件目录
> >
> > views 前端目录



#### 安装教程
1.项目安装
```
> git clone XXX
>cd BeegoBlog
>go mod init
>go mod tidy
>cd AllSrcCode
>go run main.go
```
2.项目配置：

添加AllSrcCode/conf/app.conf配置文件，内容为：
```bigquery
appname       = Blog
httpport      = 8080
runmode       = dev
# sessionon     = true
# sessionname   = beegosessionID
# [mysql]
mysqlHost     = localhost
mysqlPort     = 3306
mysqlUser     = root
mysqlPassword = 1111
mysqlDb       = myBlogWeb
mysqlPrefix   = tb_
# [redis]
redisAddress  = localhost:6379
redisPassword = 1111
redisDb       = 0
# 指定缓存的前缀
redisKey      = beego_cache
# [log]
logsfilename  = logs/app.log
level         = 3
jwtSecretKey  = c5006b96-71c1-4f78-9701-b3ace88498ac
# CSRF
enablexsrf    = false
xsrfkey       = 61oETzKXQAGaYdkL5gEmGeJJFuYh7EQnp2XdTP1o
xsrfexpire    = 3600

```
若要开启csrf,请将enablexsrf设置为true

日志监控:

若项目运行时未自动建立logs文件夹，可以手动添加AllSrcCode/logs/app.log

文件上传:

若项目运行时无法上传文件，可在自行建立AllSrcCode/static/uploads

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  Gitee 官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解 Gitee 上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
5.  Gitee 官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)

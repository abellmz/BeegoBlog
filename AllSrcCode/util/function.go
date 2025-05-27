package util

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/bwmarrin/snowflake"
	"time"
)

func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// snowflake
func GenerateSnowflakeId() int64 {
	// 设置起始时间和机器 ID
	startTime := "2025-05-20"
	machineID := 1

	// 解析起始时间 布局字符串 "2006-01-02" 是一个固定的格式模板，用于表示日期的格式
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		logs.Error("Snowflake time.Parse error", err)
		return 0
	}

	// 设置雪花算法的起始时间
	snowflake.Epoch = st.UnixNano() / 1000000

	// 创建节点实例
	node, err := snowflake.NewNode(int64(machineID))
	if err != nil {
		logs.Error("Snowflake NewNode error", err)
		return 0
	}

	// 生成并输出 ID
	return node.Generate().Int64()
}

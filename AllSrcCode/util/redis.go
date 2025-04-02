package util

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

// 创建带有超时功能的上下文
func WithTimeout(d int) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d)*time.Second)
	defer cancel()
	return ctx
}

/*
*
redis初始化
*/
func InitRedis() {
	redisHost := beego.AppConfig.String("redisAddress")
	redisPassword := beego.AppConfig.String("redisPassword")
	redisDb, err := beego.AppConfig.Int("redisDb")
	if err != nil {
		logs.Error("Failed to read redisDb: %v", err)
		return
	}
	//配置 Redis 连接选项
	opt := &redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       redisDb,
		PoolSize: 10,
	}
	//创建一个新的 Redis 客户端实例
	rdb = redis.NewClient(opt)
	// 使用默认的上下文
	err = rdb.Ping(ctx).Err()
	if err != nil {
		logs.Error("Failed to connect to Redis: %v", err)
		return
	}
	logs.Info("Redis connected successfully!")
}

// Set 设置缓存
func Set(key string, value interface{}, expiration time.Duration) error {
	_, err := rdb.Set(ctx, key, value, expiration).Result()
	if err != nil {
		logs.Error("Failed to set cache: %v", err)
		return err
	}
	fmt.Println("---------12213-----")
	return nil
}

// Get 获取缓存
func Get(key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		//if err == redis.Nil {
		//	return "", fmt.Errorf("key not found")
		//}
		logs.Error("Failed to get cache: %v", err)
		return "", err
	}
	return val, nil
}

// Del 删除缓存
func Del(key string) error {
	_, err := rdb.Del(ctx, key).Result()
	if err != nil {
		logs.Error("Failed to delete cache: %v", err)
		return err
	}
	return nil
}

// Exists 检查键是否存在
func Exists(key string) (bool, error) {
	val, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		logs.Error("Failed to check existence: %v", err)
		return false, err
	}
	return val > 0, nil
}

// TTL 获取键的剩余生存时间
func TTL(key string) (time.Duration, error) {
	return rdb.TTL(ctx, key).Result()
}

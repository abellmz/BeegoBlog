package util

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// 定义 JWT 参数，cuetomer claims
type CustomClaims struct {
	jwt.StandardClaims        // 嵌入标准声明
	UserID             string `json:"userId"`
	UserName           string `json:"userName"`
}

// 生成 Token
func GenerateToken(userId, userName string, secretKey string) (string, error) {
	//自定义声明
	claims := CustomClaims{
		UserID:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), //设置过期时间
		},
	}
	//创建了一个新的 JWT 实例，指定了签名算法（HMAC-SHA256）和负载（claims），
	//此时，JWT 的头部和负载部分已经准备好，但签名部分尚未生成。token是未签名的中间状态。
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//将密钥转换为字节数组,使用指定的签名算法和密钥对头部和负载进行签名，生成签名部分。
	//再将头部、负载和签名部分组合成一个完整的 JWT 字符串
	return token.SignedString([]byte(secretKey))
}

// 解析 Token
func ParseToken(tokenString string, secretKey string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	//负载的类型正确，如果签名无效，则说明 JWT 可能在传输过程中被篡改。
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil //若负载部分Claims类型正确，且签名正确，直接返回负载信息
	}
	return nil, err
}

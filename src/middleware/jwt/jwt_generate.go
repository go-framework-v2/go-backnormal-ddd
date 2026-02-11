package jwt

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JwtSecret 定义用于签名 token 的密钥
var JwtSecret = []byte("huanlema916")

// MyClaims 定义一个 Claims 结构体，它将嵌入 jwt.StandardClaims
type MyClaims struct {
	UserID             int64 `json:"userId"` // 自定义声明：用户 ID
	jwt.StandardClaims       // 嵌入 StandardClaims，它包含签发者、过期时间等信息
}

// GenerateToken 生成 JWT token
func GenerateToken(userID int64) (string, error) {
	// // 创建 Claims 实例，并设置过期时间为 24 小时
	// expirationTime := time.Now().Add(constants.TokenExpirationTime)
	// claims := &MyClaims{
	// 	UserID: userID,
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: expirationTime.Unix(),
	// 	},
	// }
	// 创建 Claims 实例，并设置过期时间为 7 天
	claims := &MyClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	}
	// 创建 token，指定签名算法为 HS256，并传入 Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用密钥签名 token
	signedToken, err := token.SignedString(JwtSecret)
	if err != nil {
		log.Fatalf("GenerateToken failed: %v", err)
		return "", err
	}

	return signedToken, nil
}

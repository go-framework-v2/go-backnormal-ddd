package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-framework-v2/go-access/access"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 中间件函数用于校验JWT
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		emptyInterfaceData := interface{}(nil)
		// 从Authorization头部获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			out := access.GetErrorResult(http.StatusUnauthorized, emptyInterfaceData, "Authorization header missing")
			c.JSON(http.StatusOK, out)
			return
		}

		// 处理Bearer Token（兼容带/不带Bearer前缀的情况）
		tokenString := ""
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString = parts[1]
		} else if len(parts) == 1 {
			// 兼容不带Bearer前缀的情况
			tokenString = parts[0]
		} else {
			out := access.GetErrorResult(http.StatusUnauthorized, emptyInterfaceData, "Authorization格式错误, 应为'Bearer [token]'或直接使用token")
			c.Abort()
			c.JSON(http.StatusOK, out)
			return
		}

		// 解析和验证 token
		token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			// 检查签名方法是否是你期望的（例如 HMAC SHA256）
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				msg := fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"])
				out := access.GetErrorResult(http.StatusOK, emptyInterfaceData, msg)
				c.JSON(http.StatusOK, out)
				c.Abort()
				return nil, fmt.Errorf(msg)
			}
			return JwtSecret, nil
		})

		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					out := access.GetErrorResult(http.StatusOK, emptyInterfaceData, "Malformed token")
					c.JSON(http.StatusOK, out)
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					out := access.GetErrorResult(http.StatusOK, emptyInterfaceData, "Token is expired")
					c.JSON(http.StatusOK, out)
				} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
					out := access.GetErrorResult(http.StatusOK, emptyInterfaceData, "Invalid token signature")
					c.JSON(http.StatusOK, out)
				} else {
					out := access.GetErrorResult(http.StatusOK, emptyInterfaceData, "Invalid token")
					c.JSON(http.StatusOK, out)
				}
			} else {
				out := access.GetErrorResult(http.StatusOK, emptyInterfaceData, "Invalid token")
				c.JSON(http.StatusOK, out)
			}
			c.Abort()
			return
		}

		// 如果到这里，说明 token 是有效的
		// 从 token 中提取用户信息，并将其附加到上下文中
		claims, ok := token.Claims.(*MyClaims)
		if !ok || !token.Valid {
			out := access.GetErrorResult(http.StatusOK, emptyInterfaceData, "Invalid token")
			c.JSON(http.StatusOK, out)
			c.Abort()
			return
		}
		c.Set("userId", claims.UserID)
		// 你可以在这里打印userID，以便调试
		// if userIdPrint, exists := c.Get("userId"); exists {
		// 	fmt.Println("userId: ", userIdPrint)
		// } else {
		// 	fmt.Println("userId is nil")
		// }

		// 继续处理请求
		c.Next()
	}
}

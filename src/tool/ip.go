package tool

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetIpByGinContext 获取请求IP地址
func GetIpByGinContext(c *gin.Context) string {
	ip := c.ClientIP()

	// 检查 X-Forwarded-For 头以获取真实 IP
	forwarded := c.Request.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			ip = strings.TrimSpace(ips[0])
		}
	} else if realIP := c.Request.Header.Get("X-Real-IP"); realIP != "" {
		ip = realIP
	}

	// 如果仍然没有获取到，使用默认的 ClientIP
	if ip == "" {
		ip = c.ClientIP()
	}

	// 仅当 IP 是回环地址时，不覆盖真实 IP 地址
	if net.ParseIP(ip).IsLoopback() {
		// 这里不需要覆盖 IP，保持原始 IP
	}

	return ip
}

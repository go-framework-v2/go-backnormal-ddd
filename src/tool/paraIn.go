package tool

import (
	"github.com/gin-gonic/gin"
	"github.com/go-framework-v2/go-access/access"
	"github.com/go-framework-v2/go-validate/validate"
	"go.uber.org/zap"
)

// ShouldBindJSON_Validate 封装了 ShouldBindJSON 和参数校验的逻辑
func ShouldBindJSON_Validate[T any](c *gin.Context) (*access.ParaIn[T], error) {
	var in access.ParaIn[T]

	// 1. 解析请求参数
	if err := c.ShouldBindJSON(&in); err != nil {
		zap.S().Error(err)
		return nil, err
	}

	// 2. 参数校验
	if err := validate.Validate(in); err != nil {
		zap.S().Error(err)
		return nil, err
	}

	// 3. 返回解析后的参数
	return &in, nil
}

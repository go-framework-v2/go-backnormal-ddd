package tool

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-framework-v2/go-access/access"
)

// serviceFunc模板为func XX(in dto.XXPraIn) (*dto.XXParaOut, error)
type HandlerConfig[T any, R any] struct {
	EmptyData    R
	ServiceFunc  func(T) (*R, error)
	NeedBindJSON bool
}

func GenericHandler[T any, R any](c *gin.Context, config HandlerConfig[T, R]) {
	var in access.ParaIn[T]
	var err error

	// 如果需要绑定JSON
	if config.NeedBindJSON {
		inPtr, bindErr := ShouldBindJSON_Validate[T](c)
		if bindErr != nil {
			out := access.GetErrorResult(http.StatusInternalServerError, config.EmptyData, bindErr.Error())
			c.JSON(http.StatusOK, out)
			return
		}
		in = *inPtr
	}

	// 执行业务逻辑
	outData, err := config.ServiceFunc(in.Data)
	if err != nil {
		out := access.GetErrorResult(http.StatusInternalServerError, config.EmptyData, err.Error())
		c.JSON(http.StatusOK, out)
		return
	}

	// 返回成功响应
	out := access.GetSuccessResult(outData, "success")
	c.JSON(http.StatusOK, out)
}

func HandleWithBind[In any, Out any](c *gin.Context, serviceFunc func(In) (*Out, error), emptyData Out) {
	config := HandlerConfig[In, Out]{
		EmptyData:    emptyData,
		ServiceFunc:  serviceFunc,
		NeedBindJSON: true,
	}
	GenericHandler(c, config)
}

func HandleWithoutBind[In any, Out any](c *gin.Context, serviceFunc func(In) (*Out, error), emptyData Out) {
	config := HandlerConfig[In, Out]{
		EmptyData:    emptyData,
		ServiceFunc:  serviceFunc,
		NeedBindJSON: false,
	}
	GenericHandler(c, config)
}

// serviceFunc模板为func XX(c *gin.Context, in dto.XXPraIn) (*dto.XXParaOut, error)
type HandlerConfigWithC[T any, R any] struct {
	EmptyData        R
	ServiceFuncWithC func(*gin.Context, T) (*R, error)
	NeedBindJSON     bool
}

func GenericHandlerWithC[T any, R any](c *gin.Context, config HandlerConfigWithC[T, R]) {
	var in access.ParaIn[T]
	var err error

	// 如果需要绑定JSON
	if config.NeedBindJSON {
		inPtr, bindErr := ShouldBindJSON_Validate[T](c)
		if bindErr != nil {
			out := access.GetErrorResult(http.StatusInternalServerError, config.EmptyData, bindErr.Error())
			c.JSON(http.StatusOK, out)
			return
		}
		in = *inPtr
	}

	// 执行业务逻辑
	outData, err := config.ServiceFuncWithC(c, in.Data)
	if err != nil {
		out := access.GetErrorResult(http.StatusInternalServerError, config.EmptyData, err.Error())
		c.JSON(http.StatusOK, out)
		return
	}

	// 返回成功响应
	out := access.GetSuccessResult(outData, "success")
	c.JSON(http.StatusOK, out)
}

func HandleWithBindWithC[In any, Out any](c *gin.Context, serviceFuncWithC func(*gin.Context, In) (*Out, error), emptyData Out) {
	config := HandlerConfigWithC[In, Out]{
		EmptyData:        emptyData,
		ServiceFuncWithC: serviceFuncWithC,
		NeedBindJSON:     true,
	}
	GenericHandlerWithC(c, config)
}

func HandleWithoutBindWithC[In any, Out any](c *gin.Context, serviceFuncWithC func(*gin.Context, In) (*Out, error), emptyData Out) {
	config := HandlerConfigWithC[In, Out]{
		EmptyData:        emptyData,
		ServiceFuncWithC: serviceFuncWithC,
		NeedBindJSON:     false,
	}
	GenericHandlerWithC(c, config)
}

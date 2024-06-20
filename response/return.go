package response

import (
	"github.com/xenochrony/xylitol"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TagReturn 定义日志标签，用于标记返回相关的日志条目。
const TagReturn = "[Return]"

// DefaultReturn 是一个默认的响应实例，作为Error和OK函数中响应对象的初始模板。
var DefaultReturn = &Response{}

// Error 用于处理并返回错误响应。
// 设置错误代码、消息，并记录错误日志，最后向客户端发送错误响应。
// @param c *gin.Context: Gin框架的上下文，包含HTTP请求和响应信息。
// @param code int: 错误状态码。
// @param err error: 错误对象，用于获取错误信息。
// @param msg string: 自定义错误消息。
func Error(c *gin.Context, code int, err error, msg string) {
	// 克隆默认响应对象以复用
	res := DefaultReturn.Clone()
	res.Success(false)                                // 标记响应为失败
	res.SetTraceID(xylitol.Default.GenerateTraceId()) // 设置追踪ID
	res.SetCode(int32(code))                          // 设置错误代码
	res.SetMsg(msg)                                   // 设置错误消息
	res.SetInfo(msg)                                  // 设置附加信息（与msg相同，可根据实际情况调整）
	if err != nil {                                   // 如果有具体的错误对象，则设置错误信息
		res.SetInfo(err.Error())
	}
	// 记录错误日志
	xylitol.Default.Error(TagReturn+c.FullPath(), res.GetFields()...)
	// 将响应对象放入上下文中
	c.Set("result", res)
	// 向客户端发送错误响应并终止后续中间件处理
	c.AbortWithStatusJSON(http.StatusOK, res)
}

// OK 用于处理并返回成功响应。
// 设置成功标志、数据及消息，并记录日志，最后向客户端发送成功响应。
// @param c *gin.Context: Gin框架的上下文，包含HTTP请求和响应信息。
// @param data any: 成功响应携带的数据。
// @param msg string: 成功消息。
func OK(c *gin.Context, data any, msg string) {
	// 克隆默认响应对象
	res := DefaultReturn.Clone()
	res.Success(true)                                 // 标记响应为成功
	res.SetTraceID(xylitol.Default.GenerateTraceId()) // 设置追踪ID
	res.SetCode(http.StatusOK)                        // 设置状态码为200
	res.SetMsg(msg)                                   // 设置成功消息
	res.SetInfo(msg)                                  // 设置附加信息（与msg相同，可根据实际情况调整）
	res.SetData(data)                                 // 设置响应数据
	// 记录成功日志
	xylitol.Default.Info(TagReturn+c.FullPath(), res.GetFields()...)
	// 将响应对象放入上下文中
	c.Set("result", res)
	// 向客户端发送成功响应并终止后续中间件处理
	c.AbortWithStatusJSON(http.StatusOK, res)
}

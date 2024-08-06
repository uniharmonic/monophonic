package middleware

import (
	"github.com/uniharmonic/monophonic"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinRecovery 是一个 Gin 中间件函数，用于捕获并恢复项目中可能出现的 panic 错误，
// 确保服务在遇到运行时错误时仍能保持稳定运行。它还提供了日志记录功能，并可选地记录调用栈信息。
//
// Parameters:
// - logger (*zap.Logger): Zap 日志库的 Logger 实例，用于记录恢复过程中的日志信息。
// - stack (bool): 表示是否在日志中包含调用栈信息。如果为 true，则会记录发生 panic 时的完整调用栈。
//
// Returns:
// - gin.HandlerFunc: 返回一个 Gin 处理函数，符合中间件的定义。
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 defer-recover 机制捕获 panic
		defer func() {
			if err := recover(); err != nil {
				// 检查错误是否由于断开的连接引起（如"broken pipe"或"connection reset by peer"）
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						brokenPipe = strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer")
					}
				}

				// 记录请求详情以供调试
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				if brokenPipe {
					// 对于断开的连接，仅记录错误和请求信息，不尝试写入响应状态
					monophonic.Default.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error)) // 记录错误但不检查错误，因为连接已断开
					c.Abort()            // 终止请求处理
					return               // 从 defer 中返回，避免执行后续的 AbortWithStatus
				}

				// 根据配置决定是否记录调用栈信息
				logFields := []zap.Field{
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
				}
				if stack {
					logFields = append(logFields, zap.String("stack", string(debug.Stack())))
				}
				monophonic.Default.Error("[Recovery from panic]", logFields...)

				// 终止当前请求并返回内部服务器错误状态码
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		// 调用下一个中间件或路由处理函数
		c.Next()
	}
}

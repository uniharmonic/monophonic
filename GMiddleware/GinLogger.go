package GMiddleware

/**
 * @File:   GinLogger.go
 * @Author: easternday <easterNday@foxmail.com>
 * @Date:   6/19/24 11:16 PM
 * @Create       easternday 2024-06-19 11:16 PM
 * @Update       easternday 2024-06-19 11:16 PM
 */

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xInitialization/xBackstage/internal/pkg/xLogger/GLogger"
	"github.com/xInitialization/xBackstage/internal/pkg/xLogger/GResponse"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TagDefault 定义了日志记录中的默认接收标签，用于标记接收到请求的记录。
const TagDefault = "[Receive]"

// maxMemory 定义了处理请求体时允许的最大内存大小，单位为字节。
const maxMemory = 32 << 20 // 32MB

// GinLogger 返回一个Gin中间件处理器，用于记录请求的详细日志信息。
// 每当请求到达时，此中间件会先记录请求的初步信息。
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		GLogger.Default.Info(TagDefault+c.FullPath(), GetFields(c)...)
	}
}

// GetFields 根据Gin的上下文信息构建日志字段切片。
// 这些字段包括请求处理耗时、响应状态码、请求方法、路径、查询参数、客户端IP、User-Agent、错误信息等。
// 若上下文中存在"result"键且其值不为空，则还会添加追踪ID字段。
func GetFields(c *gin.Context) []zapcore.Field {
	start := time.Now()
	c.Next() // 继续执行后续的处理函数
	cost := time.Since(start).Milliseconds()

	fields := []zapcore.Field{}

	// 尝试从上下文中提取并添加追踪ID
	if res, ok := c.Get("result"); ok && res != nil {
		traceID := res.(*GResponse.Response).TraceID
		fields = append(fields, zap.String("traceId", traceID))
	}

	// 添加其他标准日志字段
	return append(fields,
		zap.Int("status", c.Writer.Status()),                                 // HTTP响应状态码
		zap.String("method", c.Request.Method),                               // 请求方法
		zap.String("path", c.Request.URL.Path),                               // 请求路径
		zap.String("query", getParams(c)),                                    // 请求查询参数或POST数据
		zap.String("ip", c.ClientIP()),                                       // 客户端IP地址
		zap.String("user-agent", c.Request.UserAgent()),                      // 用户代理信息
		zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()), // 私有错误信息
		zap.Int64("cost", cost),                                              // 请求处理耗时（毫秒）
	)
}

// getParams 根据不同的请求类型解析并返回请求参数。
// 支持URL查询字符串、表单数据（包括x-www-form-urlencoded和multipart/form-data）以及直接读取请求体。
func getParams(c *gin.Context) string {
	// 解析请求类型，获取请求参数结构体
	if c.Request.Method == "POST" {
		contentType := strings.Split(c.Request.Header.Get("Content-Type"), ";")[0]
		switch contentType {
		case "application/x-www-form-urlencoded":
			if err := c.Request.ParseForm(); err == nil {
				values := c.Request.PostForm
				jsonByte, _ := json.Marshal(values)
				return string(jsonByte)
			}
		case "application/form-data":
			if err := c.Request.ParseMultipartForm(maxMemory); err == nil {
				values := c.Request.PostForm
				jsonByte, _ := json.Marshal(values)
				return string(jsonByte)
			}
		case "multipart/form-data":
			if err := c.Request.ParseMultipartForm(maxMemory); err == nil {
				values := c.Request.PostForm
				jsonByte, _ := json.Marshal(values)
				return string(jsonByte)
			}
		default:
			if requestBody, err := io.ReadAll(c.Request.Body); err == nil {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
				return string(requestBody)
			}
		}
	}
	return c.Request.URL.RawQuery
}

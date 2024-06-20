package logger

import "go.uber.org/zap/zapcore"

// GetLogLevel 根据输入的日志级别字符串，转换并返回对应的 zapcore.Level 枚举值。
// 此函数支持将常见的日志级别字符串（如 "debug", "info", "warn", "error", "fatal", "panic"）
// 转换为 zapcore 库使用的日志级别标识。
// 若提供的字符串无法识别，则默认返回 zapcore.InfoLevel。
//
// @param level string: 需要转换的日志级别字符串，大小写不敏感。
// @return zapcore.Level: 相应的日志级别枚举值。若转换失败，默认为InfoLevel。
//
// 示例：
//   - 输入 "debug" 或 "DEBUG"，转换结果为 zapcore.DebugLevel。
//   - 输入 "info" 或 "INFO"，转换结果为 zapcore.InfoLevel。
//   - 若输入如 "invalid" 等无法识别的字符串，则返回 zapcore.InfoLevel。
func GetLogLevel(level string) zapcore.Level {
	// 尝试将字符串转换为 zapcore.Level
	lv, err := zapcore.ParseLevel(level)
	if err != nil {
		// 转换失败时，返回默认的日志级别InfoLevel
		return zapcore.InfoLevel
	}
	return lv
}

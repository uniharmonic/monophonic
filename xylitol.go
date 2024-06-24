package xylitol

import (
	"github.com/xenochrony/xylitol/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Default 是一个默认初始化的 GLogger 实例，方便全局访问。
var Default = New("debug", "tmp/run.log")

// New 初始化并返回一个新的 Ginebra 日志实例。
// 此函数根据配置设置日志级别、路径以及输出目的地（控制台和/或文件）。
// 适合在应用程序启动时调用，以配置整个应用的日志行为。
//
// @Description:
//
//	初始化日志模块，配置日志级别、输出格式及存储位置。
//
// @Return *GLogger: 返回配置好的 GLogger 实例，可用于日志记录。
// TODO: 考虑后期日志输出级别从环境变量中获取，以及动态配置日志级别
func New(level string, logfile string) *logger.GLogger {
	// 从配置中获取日志级别和路径信息
	logLevel := logger.GetLogLevel(level)
	logPath := logfile

	// 配置日志编码器，用于格式化输出到控制台的日志
	encoder := logger.GetEncoder()

	// 准备文件写入器，用于将日志记录到指定文件
	fileWriteSyncer := logger.GetFileLogWriter(logPath)

	// 设置日志核心，允许同时输出到控制台和文件，根据环境调整此逻辑
	core := zapcore.NewTee(
		// 注意：生产环境中应考虑移除或调整控制台输出
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), logLevel),
		zapcore.NewCore(encoder, fileWriteSyncer, logLevel),
	)

	// 创建并返回 GLogger 实例，其中包含日志级别信息及 zap.Logger 的封装
	// 添加 zap.AddCaller 和 zap.AddCallerSkip 以便在日志中记录调用者信息
	return &logger.GLogger{
		ZapLogger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2)),
		LogLevel:  level,
		LogPath:   logPath,
	}
}

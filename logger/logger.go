package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

/*
LogInterface 接口定义了一套日志操作方法，
旨在为应用程序提供统一且灵活的日志记录能力。
包括生成追踪ID及不同日志级别（调试、信息、警告、错误、致命）的记录功能。

方法列表：
  - GenerateTraceId：生成一个唯一标识，便于日志追踪与问题排查。
  - Debug：记录调试信息，适用于开发阶段详细日志输出。
  - Info：记录一般信息日志，用于追踪应用运行流程。
  - Warn：记录警告信息，指出可能的问题但不影响当前执行流程。
  - Error：记录错误信息，指示发生了应当被关注并处理的错误情况。
  - Fatal：记录致命错误，并在执行该方法后终止程序运行。

参数说明：
  - msg：所有记录日志方法中的字符串参数，表示要记录的日志消息内容。
  - fields：可变参数，zapcore.Field类型的切片，用于携带额外的键值对信息，
    丰富日志内容，提高日志的可读性和分析便利性。
*/
type LogInterface interface {
	GenerateTraceId() string                   // 生成一个用于日志追踪的唯一ID。
	Debug(msg string, fields ...zapcore.Field) // 记录调试日志。
	Info(msg string, fields ...zapcore.Field)  // 记录信息日志。
	Warn(msg string, fields ...zapcore.Field)  // 记录警告日志。
	Error(msg string, fields ...zapcore.Field) // 记录错误日志。
	Fatal(msg string, fields ...zapcore.Field) // 记录致命错误日志后终止程序。
	SetLogLevel(level string)                  // 设置日志级别。
}

/*
GLogger 结构体封装了日志功能，集成 zap.Logger 提供高性能日志记录能力。
通过 ZapLogger 成员直接利用 zap 库的功能，并通过 LogLevel 成员控制日志输出级别。

属性说明：
  - ZapLogger：zap.Logger 的指针，作为日志记录的核心实现对象。
    提供了丰富的日志处理能力，如格式化、过滤和输出目标配置等。
  - LogLevel：日志级别枚举，来自 logger.LogLevel，用于设定日志输出的最低级别。
    允许动态调整以适应不同的运行环境（如生产、开发）对日志详略的需求。
*/
type GLogger struct {
	ZapLogger *zap.Logger // zap 日志库的实例，负责实际的日志处理工作。
	LogLevel  string      // 当前日志记录的最低级别门槛。
	LogPath   string      // 日志路径
}

// GetEncoder 创建并返回一个zapcore.Encoder，用于格式化日志输出至控制台。
// 该函数配置了日志的显示样式，包括时间格式、日志级别颜色高亮以及完整的调用者信息。
//
// @Description:
//
//	初始化控制台日志编码器，定制日志输出格式，包括时间、级别颜色及调用者信息。
//
// @Return zapcore.Encoder: 返回配置好的控制台日志编码器实例。
func GetEncoder() zapcore.Encoder {
	// 使用zap的生产环境默认配置作为基础
	encoderConfig := zap.NewProductionEncoderConfig()

	// 自定义配置：
	// 1. 时间戳格式设为ISO8601标准格式，便于国际标准化解析
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 2. 日志级别使用大写字母表示并带有颜色区分，增强可读性
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// 3. 调用者信息记录全路径，便于追踪日志来源
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder

	// 根据上述配置创建控制台编码器实例
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// GetFileLogWriter 根据给定的文件路径创建并返回一个实现了 zapcore.WriteSyncer 接口的对象，
// 用于日志文件的写入与同步。使用 lumberjack 库来支持日志文件的滚动、压缩和清理。
//
// @param logPath string: 日志文件的保存路径。
// @return zapcore.WriteSyncer: 返回配置好的日志文件写入器。
func GetFileLogWriter(logPath string) zapcore.WriteSyncer {
	// lumberjack.Logger 配置：
	// - Filename: 日志文件名
	// - MaxSize: 单个日志文件最大大小，默认单位为MB
	// - MaxBackups: 保留的旧日志文件的最大数量
	// - MaxAge: 旧日志文件保留的最长时间，单位天
	// - Compress: 是否启用日志文件压缩，默认不压缩
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100,  // 修改为单个文件最大100M，原说明有误
		MaxBackups: 60,   // 保留最多60个备份文件
		MaxAge:     30,   // 修改为最多保留30天的日志文件，原说明有误
		Compress:   true, // 压缩旧日志文件，提高存储效率
	}

	// zapcore.AddSync 将 lumberjack.Logger 包装成 zapcore.WriteSyncer
	return zapcore.AddSync(lumberJackLogger)
}

// Info 记录信息级别的日志。
// @param msg string: 日志消息。
// @param fields ...zapcore.Field: 额外的结构化日志字段。
func (log *GLogger) Info(msg string, fields ...zapcore.Field) {
	log.ZapLogger.Info(msg, fields...)
}

// Debug 记录调试级别的日志。
// @param msg string: 日志消息。
// @param fields ...zapcore.Field: 额外的结构化日志字段。
func (log *GLogger) Debug(msg string, fields ...zapcore.Field) {
	log.ZapLogger.Debug(msg, fields...)
}

// Warn 记录警告级别的日志。
// @param msg string: 日志消息。
// @param fields ...zapcore.Field: 额外的结构化日志字段。
func (log *GLogger) Warn(msg string, fields ...zapcore.Field) {
	log.ZapLogger.Warn(msg, fields...)
}

// Error 记录错误级别的日志。
// @param msg string: 日志消息。
// @param fields ...zapcore.Field: 额外的结构化日志字段。
func (log *GLogger) Error(msg string, fields ...zapcore.Field) {
	log.ZapLogger.Error(msg, fields...)
}

// Fatal 记录致命错误级别的日志，并在记录后调用 os.Exit(1) 终止程序。
// @param msg string: 日志消息。
// @param fields ...zapcore.Field: 额外的结构化日志字段。
func (log *GLogger) Fatal(msg string, fields ...zapcore.Field) {
	log.ZapLogger.Fatal(msg, fields...)
}

func (log *GLogger) SetLogLevel(level string) {
	log.LogLevel = level

	// 从配置中获取日志级别和路径信息
	logLevel := GetLogLevel(log.LogLevel)
	logPath := log.LogPath

	// 配置日志编码器，用于格式化输出到控制台的日志
	encoder := GetEncoder()

	// 准备文件写入器，用于将日志记录到指定文件
	fileWriteSyncer := GetFileLogWriter(logPath)

	// 设置日志核心，允许同时输出到控制台和文件，根据环境调整此逻辑
	core := zapcore.NewTee(
		// 注意：生产环境中应考虑移除或调整控制台输出
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), logLevel),
		zapcore.NewCore(encoder, fileWriteSyncer, logLevel),
	)

	log.ZapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
}

package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/uniharmonic/monophonic"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const TAG = "[GORM]"

//type GormLoggerInterface interface {
//	LogMode(level logger.LogLevel) GormLoggerInterface
//	Info(context.Context, string, ...interface{})
//	Warn(context.Context, string, ...interface{})
//	Error(context.Context, string, ...interface{})
//	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
//}

type GormLogger struct {
	SlowThreshold time.Duration
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	switch level {
	case logger.Silent:
		monophonic.Default.SetLogLevel("fatal")
	case logger.Info:
		monophonic.Default.SetLogLevel("debug")
	case logger.Warn:
		monophonic.Default.SetLogLevel("warn")
	case logger.Error:
		monophonic.Default.SetLogLevel("error")
	default:
		monophonic.Default.SetLogLevel("debug")
	}
	return l
}

func (l *GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	msg := fmt.Sprintf("%s Info: %s", TAG, fmt.Sprintf(str, args...))
	monophonic.Default.Info(msg)
}

func (l *GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	msg := fmt.Sprintf("%s Warn: %s", TAG, fmt.Sprintf(str, args...))
	monophonic.Default.Warn(msg)
}

func (l *GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	msg := fmt.Sprintf("%s Error: %s", TAG, fmt.Sprintf(str, args...))
	monophonic.Default.Error(msg)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// 获取运行时间
	elapsed := time.Since(begin)
	// 获取 SQL 请求和返回条数
	sql, rows := fc()
	// 通用字段
	logFields := []zap.Field{
		zap.String("sql", sql),
		zap.Float64("time", elapsed.Seconds()),
		zap.Int64("rows", rows),
	}
	// Gorm 错误
	if err != nil {
		// 记录未找到的错误使用 warning 等级
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := fmt.Sprintf("%s %s", TAG, "ErrRecordNotFound")
			monophonic.Default.Warn(msg, logFields...)
		} else {
			msg := fmt.Sprintf("%s %s", TAG, "Error")
			// 其他错误使用 error 等级
			logFields = append(logFields, zap.Error(err))
			monophonic.Default.Error(msg, logFields...)
		}
	} else if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		msg := fmt.Sprintf("%s %s", TAG, "Slow Log")
		monophonic.Default.Warn(msg, logFields...)
	} else {
		msg := fmt.Sprintf("%s %s", TAG, "Query")
		monophonic.Default.Debug(msg, logFields...)
	}
}

func GetGormConfig(level string) *gorm.Config {
	gormLogger := &GormLogger{}
	monophonic.Default.SetLogLevel(level)
	return &gorm.Config{
		Logger: gormLogger,
	}
}

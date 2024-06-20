package logger

import "github.com/google/uuid"

// GenerateTraceId 为 GLogger 类型实例提供生成全局唯一追踪ID的功能。
// 这个方法利用 UUID 生成一个字符串，确保了每个调用生成的ID都是唯一的，
// 有助于在分布式系统中跨服务追踪请求和日志。
//
// @receiver log *GLogger: GLogger 结构体的指针，尽管在此方法中未直接使用，
//
//	但作为接收者表明此方法属于 GLogger 的实例方法。
//
// @return string: 返回一个全局唯一标识符（UUID）的字符串表示形式，用作追踪ID。
func (log *GLogger) GenerateTraceId() string {
	// 使用uuid包生成一个新的UUID
	return uuid.New().String()
}

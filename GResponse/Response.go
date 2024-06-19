package GResponse

/**
 * @File:   Response.go
 * @Author: easternday <easterNday@foxmail.com>
 * @Date:   6/19/24 11:04 PM
 * @Create       easternday 2024-06-19 11:04 PM
 * @Update       easternday 2024-06-19 11:04 PM
 */

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ResponseInterface 定义了一个接口，用于构建和操作HTTP响应对象。
// 实现此接口的类型需提供设置响应状态码、追踪ID、消息、附加信息、响应数据的方法，
// 以及标记响应成功与否、获取日志记录所需字段和克隆自身的能力。
type ResponseInterface interface {
	// SetCode 用于设置HTTP响应的状态码。
	SetCode(int32)
	// SetTraceID 设置响应的追踪ID，便于日志追踪和问题排查。
	SetTraceID(string)
	// SetMsg 设置响应的详细消息内容。
	SetMsg(string)
	// SetInfo 设置附加的响应信息，提供更多上下文。
	SetInfo(string)
	// SetData 用于装载响应携带的具体数据。
	SetData(interface{})
	// Success 依据参数决定响应是否标记为成功，影响Status字段。
	Success(bool)
	// GetFields 返回zapcore.Field的切片，便于将响应信息记录到日志中。
	GetFields() []zapcore.Field
	// Clone 返回当前响应对象的副本，实现深拷贝或浅拷贝逻辑。
	Clone() ResponseInterface
}

// Response 是一个实现了ResponseInterface接口的结构体，专为构建和管理HTTP响应而设计。
// 它封装了请求的追踪ID、响应代码、附加信息、消息描述、响应状态以及响应数据等关键信息，
// 支持JSON序列化，并通过接口方法提供灵活的操作能力。
type Response struct {
	// TraceID 是请求的唯一追踪标识，用于日志追踪。
	TraceID string `json:"requestId,omitempty"`

	// Code 表示响应的HTTP状态码，指示请求处理结果。
	Code int32 `json:"code,omitempty"`

	// Info 为响应提供额外的上下文信息。
	Info string `json:"info,omitempty"`

	// Msg 包含响应的详细消息或描述。
	Msg string `json:"msg,omitempty"`

	// Status 概括响应处理的总体状态，如"success"或"error"。
	Status string `json:"status,omitempty"`

	// Data 字段携带响应的主要数据内容，类型为any，支持任意数据类型。
	Data any `json:"data"`
}

// SetTraceID 修改响应实例的追踪ID，该ID用于日志记录和请求的唯一标识。
// @param id string: 新的追踪ID值。
func (res *Response) SetTraceID(id string) {
	res.TraceID = id
}

// SetCode 用于设定响应的代码状态，表示请求处理的结论。
// @param code int32: 响应的状态码。
func (res *Response) SetCode(code int32) {
	res.Code = code
}

// SetMsg 更新响应中的消息内容，提供详细的响应描述或信息。
// @param msg string: 响应消息文本。
func (res *Response) SetMsg(msg string) {
	res.Msg = msg
}

// SetInfo 设置额外的响应信息字段，可以包含更多上下文详情。
// @param info string: 附加的响应信息。
func (res *Response) SetInfo(info string) {
	res.Info = info
}

// SetData 用于装载响应携带的数据部分，数据类型可以是任意类型。
// @param data any: 响应携带的实际数据内容。
func (res *Response) SetData(data any) {
	res.Data = data
}

// Success 根据输入参数更新响应的状态标记，若非成功状态则设置状态为"error"。
// @param isSuccess bool: 响应是否成功的标志。
func (res *Response) Success(isSuccess bool) {
	if !isSuccess {
		res.Status = "error"
	}
}

// GetFields 将响应对象的各属性转化为zap日志字段数组。
// 这些字段通常用于记录响应的详细信息到日志中，便于后续分析与调试。
// @return []zapcore.Field: 包含响应各属性的日志字段数组。
func (res *Response) GetFields() []zapcore.Field {
	return []zapcore.Field{
		zap.String("requestId", res.TraceID),
		zap.Int32("code", res.Code),
		zap.String("info", res.Info),
		zap.String("msg", res.Msg),
		zap.String("status", res.Status),
		zap.Any("data", res.Data),
	}
}

// Clone 生成并返回当前响应对象的副本。
// 该方法执行的是浅拷贝，即响应数据的指针会被复制，而非其底层数据的深拷贝。
// @return ResponseInterface: 返回响应对象的副本，实现接口类型。
func (res *Response) Clone() ResponseInterface {
	// 通过结构体复制创建一个新的Response实例
	clonedRes := *res
	// 返回新实例的指针，符合ResponseInterface接口
	return &clonedRes
}

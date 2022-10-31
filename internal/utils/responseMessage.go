package utils

var (
	// OK
	OK  = response(200, "OK")    // 通用成功
	Err = response(500, "Error") // 通用错误

	// 服务级错误码
	ErrParam = response(10001, "Params Error")

	// 模块级错误码 - 用户模块
	ErrUserService = response(20100, "Service Exception")

	// ......
)

type ResponseRes struct {
	Code int         `json:"code"`    // 错误码
	Msg  string      `json:"message"` // 错误描述
	Data interface{} `json:"data"`    // 返回数据的数量
}

type ResponsePageRes struct {
	Code  int         `json:"code"`    // 错误码
	Msg   string      `json:"message"` // 错误描述
	Total int         `json:"total"`   // 返回数据的数量
	Data  interface{} `json:"data"`    // 返回数据的数量
}

// 自定义响应信息
func (res *ResponseRes) WithMsg(message string) ResponseRes {
	return ResponseRes{
		Code: res.Code,
		Msg:  message,
		Data: res.Data,
	}
}

// 追加响应数据
func (res *ResponseRes) WithData(data interface{}) ResponseRes {
	return ResponseRes{
		Code: res.Code,
		Msg:  res.Msg,
		Data: data,
	}
}

// 追加响应数据
func (res *ResponseRes) WithDataAndTotal(data interface{}, total int) ResponsePageRes {
	return ResponsePageRes{
		Code:  res.Code, // 错误码
		Msg:   res.Msg,  // 错误描述
		Total: total,    // 返回数据的数量
		Data:  data,     // 返回数据的数量
	}
}

// 构造函数
func response(code int, msg string) *ResponseRes {
	return &ResponseRes{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

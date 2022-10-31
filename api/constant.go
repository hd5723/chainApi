package api

type Result struct {
	Code    int         `json:"code"           `    // 编码
	Message string      `json:"message"           ` // 信息
	Data    interface{} `json:"data"           `    // 数据
}

type Results struct {
	Code    int         `json:"code"           `    // 编码
	Message string      `json:"message"           ` // 信息
	Total   int         `json:"total"           `   // 条数
	Data    interface{} `json:"data"           `    // 数据
}

package response

// AjaxResult 统一响应结构
type AjaxResult struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// 状态码常量
const (
	SUCCESS = 0
	ERROR   = 500
	WARN    = 601
)

// 业务错误码
const (
	ErrCodeUnitNameExists  = 2001
	ErrCodeGoodsNameExists = 2101
	ErrCodeDataNotFound    = 4004
)

// 成功返回
func Success() *AjaxResult {
	return &AjaxResult{
		Code:    SUCCESS,
		Message: "操作成功",
		Data:    nil,
	}
}

func SuccessWithData(data interface{}) *AjaxResult {
	return &AjaxResult{
		Code:    SUCCESS,
		Message: "操作成功",
		Data:    data,
	}
}

func SuccessWithMsg(msg string) *AjaxResult {
	return &AjaxResult{
		Code:    SUCCESS,
		Message: msg,
		Data:    nil,
	}
}

func SuccessWithMsgData(msg string, data interface{}) *AjaxResult {
	return &AjaxResult{
		Code:    SUCCESS,
		Message: msg,
		Data:    data,
	}
}

// 错误返回
func Error() *AjaxResult {
	return &AjaxResult{
		Code:    ERROR,
		Message: "操作失败",
		Data:    nil,
	}
}

func ErrorWithMsg(msg string) *AjaxResult {
	return &AjaxResult{
		Code:    ERROR,
		Message: msg,
		Data:    nil,
	}
}

func ErrorWithCode(code int, msg string) *AjaxResult {
	return &AjaxResult{
		Code:    code,
		Message: msg,
		Data:    nil,
	}
}

// 警告返回
func Warn(msg string) *AjaxResult {
	return &AjaxResult{
		Code:    WARN,
		Message: msg,
		Data:    nil,
	}
}

// 自定义返回
func Result(code int, msg string, data interface{}) *AjaxResult {
	return &AjaxResult{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

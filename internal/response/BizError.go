package response

// BizError 业务错误
type BizError struct {
	Code    int
	Message string
}

func (e *BizError) Error() string {
	return e.Message
}

func NewBizError(code int, msg string) *BizError {
	return &BizError{
		Code:    code,
		Message: msg,
	}
}

// 预定义业务错误
var (
	ErrUnitNameExists  = NewBizError(ErrCodeUnitNameExists, "单位名称已存在")
	ErrGoodsNameExists = NewBizError(ErrCodeGoodsNameExists, "商品名称已存在")
	ErrDataNotFound    = NewBizError(ErrCodeDataNotFound, "数据不存在")
)

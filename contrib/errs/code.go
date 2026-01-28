package errs

// 公用错误码
const (
	ErrOK           = 0
	ErrInternal     = 100001 // 系统内部错误
	ErrInvalidArgs  = 100002 // 参数错误
	ErrUnauthorized = 100003
)

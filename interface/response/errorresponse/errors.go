package errorresponse

type ErrorCode string

// 4xx errors
const (
	ErrInvalidRequestBody ErrorCode = "INVALID_REQUEST_BODY"
	ErrMethodNotAllowed   ErrorCode = "METHOD_NOT_ALLOWED"
	ErrNotFound           ErrorCode = "NOT_FOUND"
)

// 5xx errors
const (
	ErrInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
)

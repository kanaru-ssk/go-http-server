package errorresponse

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

type ErrorResponse struct {
	ErrorCode ErrorCode `json:"errorCode"`
}

func (m *Mapper) MapErrorResponse(errorCode ErrorCode) *ErrorResponse {
	return &ErrorResponse{errorCode}
}

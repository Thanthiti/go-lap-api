package response

type ErrorResponse struct{
	Error ErrorDetaile `json:"error"`
}

type ErrorDetaile struct{
	Code string `json:"code"` 
	Message string `json:"message"`
}
func NewErrorResponse(code ,message string) ErrorResponse{
	return ErrorResponse{
		Error: ErrorDetaile{
			Code: code,
			Message:message,
		},
	}
}
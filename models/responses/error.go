package responses

type ApiErrorResponse struct {
	Message       string   `json:"message"`
	MessageCode   string   `json:"messageCode,omitempty"`
	MessageDetail []string `json:"messageDetail,omitempty"`
}

func MakeApiErrorResponse(errs ...string) ApiErrorResponse {
	if len(errs) < 1 {
		panic("api error response must have at least (1) error")
	}

	res := ApiErrorResponse{}
	for k, e := range errs {
		if k < 1 {
			res.Message = e
		} else {
			res.MessageDetail = append(res.MessageDetail, e)
		}
	}

	return res
}

func MakeApiErrorResponseFromError(err error) ApiErrorResponse {
	return ApiErrorResponse{
		Message: "API error",
		MessageDetail: []string{
			err.Error(),
		},
	}
}

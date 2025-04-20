package views

type Baseresponse struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Error   []string    `json:"error"`
	Message string      `json:"message"`
}

func SuccessResponse(data interface{}) Baseresponse {
	return Baseresponse{
		Status: 200,
		Data:   data,
		Error:  nil,
	}
}

func ErrorResponse(message string) Baseresponse {
	return Baseresponse{
		Status:  400,
		Message: message,
	}
}

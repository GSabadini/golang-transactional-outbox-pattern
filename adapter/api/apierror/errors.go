package apierror

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var (
	ErrInvalidPayload = Error{
		Code:    "0010",
		Message: "Invalid payload",
	}

	ErrInvalidOperationType = Error{
		Code:    "0020",
		Message: "Invalid operation type",
	}

	ErrInvalidCurrency = Error{
		Code:    "0030",
		Message: "Invalid currency",
	}

	ErrUseCaseProcessing = Error{
		Code:    "0040",
		Message: "Error when processing the use case",
	}
)

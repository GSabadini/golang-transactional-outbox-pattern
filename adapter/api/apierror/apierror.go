package apierror

type (
	APIError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
)

var (
	ErrInvalidPayload = APIError{
		Code:    "0010",
		Message: "Invalid payload",
	}

	ErrInvalidOperationType = APIError{
		Code:    "0020",
		Message: "Invalid operation type",
	}

	ErrInvalidCurrency = APIError{
		Code:    "0030",
		Message: "Invalid currency",
	}

	ErrUseCaseProcessing = APIError{
		Code:    "0040",
		Message: "APIError when processing the use case",
	}
)

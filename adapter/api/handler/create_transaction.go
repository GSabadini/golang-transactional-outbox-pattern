package handler

import (
	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/api/apierror"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/logger"
	"log/slog"
	"net/http"

	"github.com/GSabadini/golang-transactional-outbox-pattern/usecase"

	"github.com/labstack/echo"
)

type CreateTransactionHandler struct {
	usecase usecase.CreateTransactionUseCase
}

func NewCreateTransactionHandler(usecase usecase.CreateTransactionUseCase) CreateTransactionHandler {
	return CreateTransactionHandler{usecase: usecase}
}

func (cth CreateTransactionHandler) Handle(c echo.Context) error {
	input := usecase.CreateTransactionInput{}

	err := c.Bind(&input)
	if err != nil {
		logger.Slog.Error("Invalid payload", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, apierror.ErrInvalidPayload)
	}

	if !input.OperationType.IsValid() {
		logger.Slog.Error("Operation type invalid", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, apierror.ErrInvalidPayload)
	}

	output, err := cth.usecase.Execute(c.Request().Context(), input)
	if err != nil {
		logger.Slog.Error("Error when processing the use case", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, apierror.ErrUseCaseProcessing)
	}

	return c.JSON(http.StatusOK, output)
}

package handler

import (
	"log/slog"
	"net/http"

	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/api/apierror"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/logger"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/opentelemetry"
	"github.com/GSabadini/golang-transactional-outbox-pattern/usecase"

	"github.com/labstack/echo"
)

type (
	TransactionHandler struct {
		usecase usecase.TransactionUseCase
	}
)

func NewTransactionHandler(usecase usecase.TransactionUseCase) TransactionHandler {
	return TransactionHandler{usecase: usecase}
}

func (t TransactionHandler) Create(c echo.Context) error {
	ctx, span := opentelemetry.NewSpan(c.Request().Context(), "handler.transaction.create")
	defer span.End()

	input := usecase.CreateTransactionInput{}

	err := c.Bind(&input)
	if err != nil {
		opentelemetry.SetError(span, err)
		logger.Slog.Error("Invalid payload", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, apierror.ErrInvalidPayload)
	}

	if !input.OperationType.IsValid() {
		opentelemetry.SetError(span, err)
		logger.Slog.Error("Invalid operation type")
		return c.JSON(http.StatusBadRequest, apierror.ErrInvalidOperationType)
	}

	if !input.Currency.IsValid() {
		opentelemetry.SetError(span, err)
		logger.Slog.Error("Invalid currency")
		return c.JSON(http.StatusBadRequest, apierror.ErrInvalidCurrency)
	}

	output, err := t.usecase.Create(ctx, input)
	if err != nil {
		opentelemetry.SetError(span, err)
		logger.Slog.Error("Error when processing the use case", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, apierror.ErrUseCaseProcessing)
	}

	return c.JSON(http.StatusOK, output)
}

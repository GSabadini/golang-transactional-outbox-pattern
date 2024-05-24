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
	AccountHandler struct {
		usecase usecase.AccountUseCase
	}
)

func NewAccountHandler(usecase usecase.AccountUseCase) AccountHandler {
	return AccountHandler{usecase: usecase}
}

func (a AccountHandler) Create(c echo.Context) error {
	ctx, span := opentelemetry.NewSpan(c.Request().Context(), "handler.account.create")
	defer span.End()

	input := usecase.CreateAccountInput{}

	err := c.Bind(&input)
	if err != nil {
		opentelemetry.SetError(span, err)
		logger.Slog.Error("Invalid payload", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, apierror.ErrInvalidPayload)
	}

	output, err := a.usecase.Create(ctx, input)
	if err != nil {
		opentelemetry.SetError(span, err)
		logger.Slog.Error("Error when processing the use case", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, apierror.ErrUseCaseProcessing)
	}

	return c.JSON(http.StatusOK, output)
}

package handler

import (
	"github.com/labstack/echo"
	"log/slog"
	"net/http"

	"github.com/GSabadini/golang-transactional-outbox-pattern/adapter/api/apierror"
	"github.com/GSabadini/golang-transactional-outbox-pattern/infra/logger"
	"github.com/GSabadini/golang-transactional-outbox-pattern/usecase"
)

type CreateAccountHandler struct {
	usecase usecase.CreateAccountUseCase
}

func NewCreateAccountHandler(usecase usecase.CreateAccountUseCase) CreateAccountHandler {
	return CreateAccountHandler{usecase: usecase}
}

func (cah CreateAccountHandler) Handle(c echo.Context) error {
	input := usecase.CreateAccountInput{}

	err := c.Bind(&input)
	if err != nil {
		logger.Slog.Error("Invalid payload", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, apierror.ErrInvalidPayload)
	}

	output, err := cah.usecase.Execute(c.Request().Context(), input)
	if err != nil {
		logger.Slog.Error("Error when processing the use case", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, apierror.ErrUseCaseProcessing)
	}

	return c.JSON(http.StatusOK, output)
}

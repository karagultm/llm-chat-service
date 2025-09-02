package chat

import (
	"myapp/pkg/logger"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type Handler interface {
	Send(c echo.Context) error
	ShowHistory(c echo.Context) error
}
type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Send(c echo.Context) error {
	logger.Log.Info("received send request")
	input := new(Chat)
	if err := c.Bind(input); err != nil {
		logger.Log.Warn("failed to bind request", zap.Error(err))
		return c.String(http.StatusBadRequest, "bad request")
	}
	if input.SessionID == "" {
		input.SessionID = uuid.New().String()
	}
	if _, err := uuid.Parse(input.SessionID); err != nil {
		logger.Log.Warn("UUID is not correct format", zap.Error(err))
		return c.String(http.StatusBadRequest, "uuid is not correct format")
	}
	if len(input.Message) < 3 || len(input.Message) > 2048 {
		logger.Log.Warn("Message is not correct format")
		return c.String(http.StatusBadRequest, "message length should be between 3 and 2048")
	}
	response, err := h.service.SendMessage(input.SessionID, input.Message)
	if err != nil {
		logger.Log.Error("service error occured", zap.Error(err))
		return c.String(http.StatusInternalServerError, "service error occured")
	}

	logger.Log.Info("request sent successfully",
		zap.String("sessionID", input.SessionID),
		zap.String("message", input.Message))

	return c.JSON(http.StatusOK, response)
}

func (h *handler) ShowHistory(c echo.Context) error {
	logger.Log.Info("received show history request")
	session_id := c.Param("sessionId")
	if _, err := uuid.Parse(session_id); err != nil {
		logger.Log.Warn("UUID is not correct format", zap.Error(err))
		return c.String(http.StatusBadRequest, "uuid is not correct format")
	}

	history, err := h.service.FindHistory(session_id)
	if err != nil {
		logger.Log.Error("service error occured", zap.Error(err))
		return c.String(http.StatusInternalServerError, "session id bulunamadı db de")
	}
	logger.Log.Info("request sent successfully",
		zap.String("sessionID", session_id))
	return c.JSON(http.StatusOK, history)

}

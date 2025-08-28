package chat

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
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

	input := new(Chat)
	if err := c.Bind(input); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if input.SessionID == "" {
		input.SessionID = uuid.New().String()
	}
	if _, err := uuid.Parse(input.SessionID); err != nil {
		return c.String(http.StatusBadRequest, "uuid is not correct format")
	}
	if len(input.Message) < 3 || len(input.Message) > 2048 {
		return c.String(http.StatusBadRequest, "message length should be between 3 and 2048")
	}
	response, err := h.service.SendMessage(input.SessionID, input.Message)
	if err != nil {
		return c.String(http.StatusInternalServerError, "service error occured")
	}

	return c.JSON(http.StatusOK, response)
}

func (h *handler) ShowHistory(c echo.Context) error {
	session_id := c.Param("sessionId")
	if _, err := uuid.Parse(session_id); err != nil {
		return c.String(http.StatusBadRequest, "uuid is not correct format")
	}

	history, err := h.service.FindHistory(session_id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "session id bulunamadı db de")
	}
	return c.JSON(http.StatusOK, history)

}

package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestSend_Success(t *testing.T) {
	// Setup
	e := echo.New()

	msg := "merhaba canım"
	id := "811360d0-462f-4fbf-b90b-ccba665986f1"
	chatJSON := `{"Message":"merhaba canım" ,"SessionID":"811360d0-462f-4fbf-b90b-ccba665986f1"}`

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := NewMockService(ctrl)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(chatJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := NewHandler(serviceMock)

	expectedChat := Chat{
		SessionID: id,
		Message:   "sana nasıl yardımcı olabilirim",
	}
	expectedJSON, err := json.Marshal(expectedChat)
	if err != nil {
		t.Fatal(err)
	}

	serviceMock.EXPECT().
		SendMessage(id, msg).
		Return(expectedChat, nil).
		Times(1)

	// Act
	err = handler.Send(c)
	//  Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, string(expectedJSON), rec.Body.String()) //handler da c.json döndüğümüzden burda bunu kullanıyoruz.Diğer error durumlarında c.string döndüğümüzden error nil döner ve metni bu equal ile karşılaştırabiliriz.

}

func TestSend_BindError(t *testing.T) {
	// Setup
	e := echo.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := NewMockService(ctrl)
	handler := NewHandler(serviceMock)

	chatJSON := `{"Message":"merhaba canım" "SessionID":"811360d0-462f-4fbf-b90b-ccba665986f1"}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(chatJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Act
	handler.Send(c)

	//  Assert

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "bad request", rec.Body.String())

}

func TestSend_EmptySessionID_GeneratesUUID(t *testing.T) {
	//gerek var mı buna bilemedim?
}
func TestSend_InvalidSessionID(t *testing.T) {
	// Setup
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := NewMockService(ctrl)
	handler := NewHandler(serviceMock)

	chatJSON := `{"Message":"merhaba canım" ,"SessionID":"bozukid"}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(chatJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Act
	handler.Send(c)

	//  Assert

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "uuid is not correct format", rec.Body.String())
}

func TestSend_MessageTooShort(t *testing.T) {
	// Setup
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := NewMockService(ctrl)
	handler := NewHandler(serviceMock)

	chatJSON := `{"Message":"Sa" ,"SessionID":"811360d0-462f-4fbf-b90b-ccba665986f1"}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(chatJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Act
	handler.Send(c)

	//  Assert

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "message length should be between 3 and 2048", rec.Body.String())
}

func TestSend_MessageTooLong(t *testing.T) {
	// Setup
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := NewMockService(ctrl)
	handler := NewHandler(serviceMock)

	longMessage := strings.Repeat("a", 3000) // 3000 karakterlik "aaaaa..."
	chatJSON := fmt.Sprintf(`{"Message":"%s","SessionID":"811360d0-462f-4fbf-b90b-ccba665986f1"}`, longMessage)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(chatJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Act
	handler.Send(c)

	//  Assert
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "message length should be between 3 and 2048", rec.Body.String())
}

func TestSend_ServiceError(t *testing.T) {
	// Setup
	e := echo.New()

	msg := "merhaba canım"
	id := "811360d0-462f-4fbf-b90b-ccba665986f1"
	chatJSON := `{"Message":"merhaba canım" ,"SessionID":"811360d0-462f-4fbf-b90b-ccba665986f1"}`

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := NewMockService(ctrl)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(chatJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := NewHandler(serviceMock)

	serviceMock.EXPECT().
		SendMessage(id, msg).
		Return(Chat{}, errors.New("service error")). //bu error ü neden hiçbir yerde çeklemiyoruz onu anlamadım?
		Times(1)

		// Act
	handler.Send(c)
	//  Assert
	// assert.Error(t, err)  bu kontrol neden yok yani. echo nil döndüğünden err yok. ozaman neden return e bir şey yazıyom ?
	// assert.EqualError(t, err, "database save error")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "service error occured", rec.Body.String())

}

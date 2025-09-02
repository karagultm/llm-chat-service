package chat

import (
	"errors"
	"myapp/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestSendMessage_NoHistory_Success(t *testing.T) {
	//arrange
	logger.Log = zap.NewNop()
	ctrl := gomock.NewController(t)
	repoMock := NewMockRepository(ctrl)
	clientMock := NewMockClient(ctrl)
	service := NewService(repoMock, clientMock)

	message := "merhaba"
	openaiMsg := "merhaba, size nasıl yardımcı olabilirim?"
	sessionId := "sess123"

	response := Chat{
		Message:   openaiMsg,
		SessionID: sessionId,
	}
	history := []ChatMessage{}

	gomock.InOrder(
		repoMock.EXPECT().Save(gomock.Any()).Do(func(msg *ChatMessage) {
			assert.Equal(t, message, msg.Message)
			assert.Equal(t, sessionId, msg.SessionID)
		}).Return(nil).Times(1),
		repoMock.EXPECT().Find(sessionId).Return(history, nil).Times(1),
		clientMock.EXPECT().GetCompletion(message, history).Return(openaiMsg, nil).Times(1),
		repoMock.EXPECT().Save(gomock.Any()).Do(func(msg *ChatMessage) {
			assert.Equal(t, openaiMsg, msg.Message)
			assert.Equal(t, sessionId, msg.SessionID)
		}).Return(nil).Times(1),
	)

	//act
	result, err := service.SendMessage(sessionId, message)
	//assert
	assert.Equal(t, response, result)
	assert.Nil(t, err)
}
func TestSendMessage_WithHistory_Success(t *testing.T) {
	//gerek var mı buna anlamadım ya
	//arrange
	logger.Log = zap.NewNop()
	ctrl := gomock.NewController(t)
	repoMock := NewMockRepository(ctrl)
	clientMock := NewMockClient(ctrl)
	service := NewService(repoMock, clientMock)

	message := "merhaba"
	openaiMsg := "merhaba, size nasıl yardımcı olabilirim?"
	sessionId := "sess123"

	response := Chat{
		Message:   openaiMsg,
		SessionID: sessionId,
	}
	history := []ChatMessage{
		{
			ID:        69,
			Kind:      "USER_PROMPT",
			Message:   "selam naber? ben talha",
			Timestamp: 1756212819,
			SessionID: "sess123",
		},
		{
			ID:        70,
			Kind:      "LLM_OUTPUT",
			Message:   "Selam Talha! Ben iyiyim, teşekkür ederim. Sen nasılsın? Nasıl yardımcı olabilirim?",
			Timestamp: 1756212821,
			SessionID: "sess123",
		},
	}

	gomock.InOrder(
		repoMock.EXPECT().Save(gomock.Any()).Do(func(msg *ChatMessage) {
			assert.Equal(t, message, msg.Message)
			assert.Equal(t, sessionId, msg.SessionID)
		}).Return(nil).Times(1),
		repoMock.EXPECT().Find(gomock.Any()).Do(func(id string) {
			assert.Equal(t, sessionId, id)
		}).Return(history, nil).Times(1),
		clientMock.EXPECT().GetCompletion(message, history).Return(openaiMsg, nil).Times(1),
		repoMock.EXPECT().Save(gomock.Any()).Do(func(msg *ChatMessage) {
			assert.Equal(t, openaiMsg, msg.Message)
			assert.Equal(t, sessionId, msg.SessionID)
		}).Return(nil).Times(1),
	)

	//act
	result, err := service.SendMessage(sessionId, message)
	//assert
	assert.Equal(t, response, result)
	assert.Nil(t, err)
}

func TestSendMessage_SaveUserMessageFails(t *testing.T) {
	//arrange
	logger.Log = zap.NewNop()
	ctrl := gomock.NewController(t)
	repoMock := NewMockRepository(ctrl)
	clientMock := NewMockClient(ctrl)
	service := NewService(repoMock, clientMock)

	message := "merhaba"
	sessionId := "sess123"
	response := Chat{}

	repoMock.EXPECT().Save(gomock.Any()).Do(func(msg *ChatMessage) {
		assert.Equal(t, message, msg.Message)
		assert.Equal(t, sessionId, msg.SessionID)
	}).Return(errors.New("database save error")).Times(1)

	//act
	result, err := service.SendMessage(sessionId, message)
	//assert
	assert.Equal(t, response, result)
	assert.Error(t, err)
	assert.EqualError(t, err, "database save error")
}

func TestSendMessage_FindHistoryFails(t *testing.T) {
	//arrange
	logger.Log = zap.NewNop()
	ctrl := gomock.NewController(t)
	repoMock := NewMockRepository(ctrl)
	clientMock := NewMockClient(ctrl)
	service := NewService(repoMock, clientMock)

	message := "merhaba"
	sessionId := "sess123"
	response := Chat{}

	gomock.InOrder(
		repoMock.EXPECT().Save(gomock.Any()).Do(func(msg *ChatMessage) {
			assert.Equal(t, message, msg.Message)
			assert.Equal(t, sessionId, msg.SessionID)
		}).Return(nil).Times(1),
		repoMock.EXPECT().Find(sessionId).Return(nil, gorm.ErrRecordNotFound).Times(1),
	)

	//act
	result, err := service.SendMessage(sessionId, message) //nasıl oldu kafam gitti
	//assert
	assert.Equal(t, response, result)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestSendMessage_GetCompletionFails(t *testing.T) {
	//arrange
	logger.Log = zap.NewNop()
	ctrl := gomock.NewController(t)
	repoMock := NewMockRepository(ctrl)
	clientMock := NewMockClient(ctrl)
	service := NewService(repoMock, clientMock)
	message := "merhaba"
	sessionId := "sess123"
	response := Chat{}
	history := []ChatMessage{}

	gomock.InOrder(
		repoMock.EXPECT().Save(gomock.Any()).Do(func(msg *ChatMessage) {
			assert.Equal(t, message, msg.Message)
			assert.Equal(t, sessionId, msg.SessionID)
		}).Return(nil).Times(1),
		repoMock.EXPECT().Find(sessionId).Return(history, nil).Times(1),
		clientMock.EXPECT().GetCompletion(message, history).Return("", errors.New("llm error")).Times(1),
	)

	//act
	result, err := service.SendMessage(sessionId, message) //nasıl oldu kafam gitti
	//assert
	assert.Equal(t, response, result)
	assert.Error(t, err)
	assert.EqualError(t, err, "llm error")
}

func TestSendMessage_SaveLLMMessageFails(t *testing.T) {
	//arrange
	logger.Log = zap.NewNop()
	ctrl := gomock.NewController(t)
	repoMock := NewMockRepository(ctrl)
	clientMock := NewMockClient(ctrl)
	service := NewService(repoMock, clientMock)

	message := "merhaba"
	openaiMsg := "merhaba, size nasıl yardımcı olabilirim?"
	sessionId := "sess123"

	response := Chat{}
	history := []ChatMessage{}

	gomock.InOrder(
		repoMock.EXPECT().Save(gomock.Any()).Do(func(msg *ChatMessage) {
			assert.Equal(t, message, msg.Message)
			assert.Equal(t, sessionId, msg.SessionID)
		}).Return(nil).Times(1),
		repoMock.EXPECT().Find(sessionId).Return(history, nil).Times(1),
		clientMock.EXPECT().GetCompletion(message, history).Return(openaiMsg, nil).Times(1),
		repoMock.EXPECT().Save(gomock.Any()).Do(func(msg *ChatMessage) {
			assert.Equal(t, openaiMsg, msg.Message)
			assert.Equal(t, sessionId, msg.SessionID)
		}).Return(errors.New("db save response error")).Times(1),
	)

	//act
	result, err := service.SendMessage(sessionId, message) //nasıl oldu kafam gitti
	//assert
	assert.Equal(t, response, result)
	assert.Error(t, err)
	assert.EqualError(t, err, "db save response error")
}
func TestFindHistory_Success(t *testing.T) {
	//arrange
	logger.Log = zap.NewNop()
	ctrl := gomock.NewController(t)
	repoMock := NewMockRepository(ctrl)
	service := NewService(repoMock, nil)

	history := []ChatMessage{
		{ID: 1, Kind: UserPrompt, Message: "merhaba", Timestamp: 111, SessionID: "session1"},
	}

	repoMock.EXPECT().Find("sess1").Return(history, nil).Times(1)
	//act
	result, err := service.FindHistory("sess1")
	//assert
	assert.Equal(t, len(history), len(result))
	assert.Nil(t, err)
	assert.Equal(t, history, result)

}
func TestFindHistory_NotFound(t *testing.T) {
	logger.Log = zap.NewNop()
	ctrl := gomock.NewController(t)
	repoMock := NewMockRepository(ctrl)
	service := NewService(repoMock, nil)

	repoMock.EXPECT().Find("sess1").Return(nil, gorm.ErrRecordNotFound)

	result, err := service.FindHistory("sess1")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

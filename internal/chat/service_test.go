package chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

// func TestSendMessage(t *testing.T) {

// }

func TestFindHistory_Success(t *testing.T) {
	//arrange
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
	assert.Equal(t, result, history)

}
func TestFindHistory_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	repoMock := NewMockRepository(ctrl)
	service := NewService(repoMock, nil)

	repoMock.EXPECT().Find("sess1").Return(nil, gorm.ErrRecordNotFound)

	result, err := service.FindHistory("sess1")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

package chat

import (
	"myapp/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	Save(message *ChatMessage) error
	Find(sessionID string) ([]ChatMessage, error)
}
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Save(message *ChatMessage) error {

	return r.db.Create(message).Error
}

func (r *repository) Find(sessionID string) ([]ChatMessage, error) {
	var messages []ChatMessage
	result := r.db.Where("session_id = ?", sessionID).Find(&messages).Order("timestamp desc")

	if result.Error != nil {
		logger.Log.Error("database find error", zap.Error(result.Error))
		return []ChatMessage{}, result.Error
	} else {
		return messages, nil
	}

}

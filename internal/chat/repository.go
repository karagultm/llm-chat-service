package chat

import (
	"myapp/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	Save(message *models.ChatMessage) error
	Find(sessionID string) ([]models.ChatMessage, error)
}
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Save(message *models.ChatMessage) error {

	return r.db.Create(message).Error
}

func (r *repository) Find(sessionID string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	result := r.db.Where("session_id = ?", sessionID).Find(&messages).Order("timestamp desc")

	if result.Error != nil {
		return []models.ChatMessage{}, result.Error
	} else {
		return messages, nil
	}

}

package repository

import (
	"context"
	"log/slog"

	"github.com/kartikx04/chat/internal/domain"
	"github.com/kartikx04/chat/internal/models"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewSessionRepository(db *gorm.DB, logger *slog.Logger) domain.SessionRepository {
	return &sessionRepository{
		db:     db,
		logger: logger,
	}
}

func (sr *sessionRepository) Create(ctx context.Context, session *models.Sessions) error {
	if err := sr.db.WithContext(ctx).Create(session).Error; err != nil {
		sr.logger.ErrorContext(ctx,
			"failed to create session",
			"id", session.ID,
			"error", err,
		)
		return err
	}

	sr.logger.InfoContext(ctx,
		"session created",
		"id", session.ID,
		"user_id", session.UserID,
	)

	return nil
}

func (sr *sessionRepository) GetByID(ctx context.Context, id string) (models.Sessions, error) {
	var session models.Sessions

	err := sr.db.WithContext(ctx).
		Where("session_id = ?", id).
		First(session).Error

	if err != nil {
		return models.Sessions{}, err
	}

	return session, nil
}

func (sr *sessionRepository) DeleteByID(ctx context.Context, id string) error {
	result := sr.db.WithContext(ctx).
		Where("session_id = ?", id).
		Delete(&models.Sessions{})

	return result.Error
}

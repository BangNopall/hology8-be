package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
)

type PresenceRepository interface {
	InsertPresence(ctx context.Context, presence *entity.Presence) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.Presence, error)
	FetchAll(ctx context.Context, offset, limit int) ([]entity.Presence, int64, error)
}

type PresenceService interface {
	CreatePresence(ctx context.Context, userID uuid.UUID) (dto.PresenceCreateResponse, error)
	CheckPresence(ctx context.Context, userID uuid.UUID) (dto.PresenceCheckResponse, error)
	FetchPresences(ctx context.Context, req *dto.PaginationRequest) ([]dto.PresenceListItem, dto.PaginationResponse, error)
}

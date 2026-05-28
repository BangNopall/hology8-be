package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
)

type PresenceService struct {
	presenceRepo contracts.PresenceRepository
}

func NewPresenceService(presenceRepo contracts.PresenceRepository) contracts.PresenceService {
	return &PresenceService{presenceRepo}
}

func (presenceSvc *PresenceService) CreatePresence(ctx context.Context, userID uuid.UUID) (dto.PresenceCreateResponse, error) {
	payload := entity.NewPresence(userID)

	if err := presenceSvc.presenceRepo.InsertPresence(ctx, payload); err != nil {
		return dto.PresenceCreateResponse{}, err
	}

	// Re-fetch with relations to obtain user fullname and team name
	pres, err := presenceSvc.presenceRepo.FindByUserID(ctx, userID)
	if err != nil {
		return dto.PresenceCreateResponse{}, err
	}
	var fullname string
	var teamName string
	var createdAt = payload.CreatedAt

	if pres != nil {
		fullname = pres.User.Fullname
		createdAt = pres.CreatedAt
		// pick the first team membership's team name if any
		if len(pres.User.Teams) > 0 && pres.User.Teams[0].Team.Name != "" {
			teamName = pres.User.Teams[0].Team.Name
		}
	}

	return dto.PresenceCreateResponse{
		UserID:    userID,
		Fullname:  fullname,
		TeamName:  teamName,
		CreatedAt: createdAt,
	}, nil
}

func (presenceSvc *PresenceService) CheckPresence(ctx context.Context, userID uuid.UUID) (dto.PresenceCheckResponse, error) {
	pres, err := presenceSvc.presenceRepo.FindByUserID(ctx, userID)
	if err != nil {
		return dto.PresenceCheckResponse{}, err
	}
	if pres == nil {
		return dto.PresenceCheckResponse{
			UserID:    userID,
			Exists:    false,
			CreatedAt: nil,
		}, nil
	}
	t := pres.CreatedAt
	return dto.PresenceCheckResponse{
		UserID:    userID,
		Exists:    true,
		CreatedAt: &t,
	}, nil
}

func (presenceSvc *PresenceService) FetchPresences(ctx context.Context, req *dto.PaginationRequest) ([]dto.PresenceListItem, dto.PaginationResponse, error) {
	presences, total, err := presenceSvc.presenceRepo.FetchAll(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, dto.PaginationResponse{}, err
	}

	items := make([]dto.PresenceListItem, 0, len(presences))
	for _, p := range presences {
		teamName := ""
		if len(p.User.Teams) > 0 {
			teamName = p.User.Teams[0].Team.Name
		}
		items = append(items, dto.PresenceListItem{
			UserID:    p.UserID,
			Fullname:  p.User.Fullname,
			TeamName:  teamName,
			CreatedAt: p.CreatedAt,
		})
	}

	totalPages := 0
	if dto.DEFAULT_SIZE > 0 {
		totalPages = int((total + int64(dto.DEFAULT_SIZE) - 1) / int64(dto.DEFAULT_SIZE))
	}

	return items, dto.PaginationResponse{
		TotalPages: totalPages,
		Page:       req.Page,
	}, nil
}

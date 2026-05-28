package service

import (
	"context"
	"time"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
)

type announcementService struct {
	announcementRepo contracts.AnnouncementRepository
	teamRepo         contracts.TeamRepository
	competitionRepo  contracts.CompetitionRepository
	timeout          time.Duration
}

func NewAnnouncementService(announcementRepo contracts.AnnouncementRepository, teamRepo contracts.TeamRepository, competitionRepo contracts.CompetitionRepository, timeout time.Duration) contracts.AnnouncementService {
	return &announcementService{announcementRepo, teamRepo, competitionRepo, timeout}
}

func (announcementSvc *announcementService) FetchAnnouncementByTo(ctx context.Context, teamID string, competitionID int) ([]dto.AnnouncementResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, announcementSvc.timeout)
	defer cancel()

	if teamID != "" && competitionID != 0 {
		return []dto.AnnouncementResponse{}, domain.ErrIllegalEntry
	}

	announcements, err := announcementSvc.announcementRepo.FetchAnnouncementByTo(ctx, teamID, competitionID)

	res := dto.AnnouncementSliceEntityToResponse(announcements)

	select {
	case <-ctx.Done():
		return []dto.AnnouncementResponse{}, domain.ErrTimeout
	default:
		return res, err
	}
}

func (announcementSvc *announcementService) CreateAnnouncement(ctx context.Context, announcement *dto.AnnouncementRequest) error {
	ctx, cancel := context.WithTimeout(ctx, announcementSvc.timeout)
	defer cancel()

	if announcement.TeamID != "" && announcement.CompetitionID != 0 {
		return domain.ErrIllegalEntry
	}

	payload := entity.NewAnnouncement(announcement.TeamID, announcement.CompetitionID, announcement.Description, announcement.AdminID)

	err := announcementSvc.announcementRepo.InsertAnnouncement(ctx, payload)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (announcementSvc *announcementService) UpdateAnnouncement(ctx context.Context, announcement *dto.AnnouncementRequest) error {
	ctx, cancel := context.WithTimeout(ctx, announcementSvc.timeout)
	defer cancel()

	if announcement.TeamID != "" && announcement.CompetitionID != 0 {
		return domain.ErrIllegalEntry
	}

	payload := entity.NewAnnouncement(announcement.TeamID, announcement.CompetitionID, announcement.Description, announcement.AdminID)
	payload.ID = announcement.ID

	err := announcementSvc.announcementRepo.UpdateAnnouncement(ctx, payload)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (announcementSvc *announcementService) DeleteAnnouncement(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, announcementSvc.timeout)
	defer cancel()

	err := announcementSvc.announcementRepo.DeleteAnnouncement(ctx, id)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

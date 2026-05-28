package contracts

import (
	"context"

	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
)

type AnnouncementRepository interface {
	FetchAnnouncementByTo(ctx context.Context, teamID string, competitionID int) ([]entity.Announcement, error)
	InsertAnnouncement(ctx context.Context, announcement *entity.Announcement) error
	UpdateAnnouncement(ctx context.Context, announcement *entity.Announcement) error
	DeleteAnnouncement(ctx context.Context, id int) error
}

type AnnouncementService interface {
	FetchAnnouncementByTo(ctx context.Context, teamID string, competitionID int) ([]dto.AnnouncementResponse, error)
	CreateAnnouncement(ctx context.Context, announcement *dto.AnnouncementRequest) error
	UpdateAnnouncement(ctx context.Context, announcement *dto.AnnouncementRequest) error
	DeleteAnnouncement(ctx context.Context, id int) error
}

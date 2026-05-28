package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/domain/entity"
	"github.com/BangNopall/hology8-be/pkg/helpers"
	"github.com/BangNopall/hology8-be/pkg/log"
)

type announcementRepository struct {
	conn *gorm.DB
}

func NewAnnouncementRepository(conn *gorm.DB) contracts.AnnouncementRepository {
	return &announcementRepository{conn}
}

func (announcementRepo *announcementRepository) FetchAnnouncementByTo(ctx context.Context, teamID string, competitionID int) ([]entity.Announcement, error) {
	var announcements []entity.Announcement
	var err error

	if teamID != "" && competitionID == 0 {
		err = announcementRepo.conn.Where("team_id = ? AND competition_id IS NULL", teamID).Find(&announcements).Error

	} else if competitionID != 0 && teamID == "" {
		err = announcementRepo.conn.Where("team_id IS NULL AND competition_id = ?", competitionID).Find(&announcements).Error

	} else {
		err = announcementRepo.conn.Where("team_id IS NULL AND competition_id IS NULL").Find(&announcements).Error

	}

	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[ANNOUNCEMENT REPOSITORY][FetchAnnouncementByTo] failed to fetch announcements")

		return []entity.Announcement{}, err
	}

	return announcements, nil
}

func (announcementRepo *announcementRepository) InsertAnnouncement(ctx context.Context, announcement *entity.Announcement) error {
	res := announcementRepo.conn.Create(announcement)

	if res.Error != nil {
		if res.Error == gorm.ErrForeignKeyViolated {
			return domain.ErrInvalidCompeTeamID
		}

		log.Error(log.LogInfo{
			"error": res.Error.Error(),
		}, "[ANNOUNCEMENT REPOSITORY][InsertAnnouncement] failed to create announcement")

		return res.Error
	}

	return nil
}

func (announcementRepo *announcementRepository) UpdateAnnouncement(ctx context.Context, announcement *entity.Announcement) error {
	res := announcementRepo.conn.Model(announcement).Updates(announcement)

	if err := helpers.CheckRowsAffected(res.RowsAffected); err != nil {
		return err
	}

	if err := res.Error; err != nil {

		if res.Error == gorm.ErrForeignKeyViolated {
			return domain.ErrInvalidCompeTeamID
		}

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[ANNOUNCEMENT REPOSITORY][UpdateAnnouncement] failed to update announcement")

		return err
	}

	return nil
}

func (announcementRepo *announcementRepository) DeleteAnnouncement(ctx context.Context, id int) error {
	res := announcementRepo.conn.Delete(&entity.Announcement{ID: id})

	if err := helpers.CheckRowsAffected(res.RowsAffected); err != nil {
		return err
	}

	if err := res.Error; err != nil {

		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[ANNOUNCEMENT REPOSITORY][DeleteAnnouncement] failed to delete announcement")

		return err
	}

	return nil
}

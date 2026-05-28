package service

import (
	"context"
	"time"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/pkg/helpers"
)

type competitionService struct {
	competitionRepo contracts.CompetitionRepository
	timeout         time.Duration
}

func NewCompetitionService(competitionRepo contracts.CompetitionRepository, timeout time.Duration) contracts.CompetitionService {
	return &competitionService{competitionRepo, timeout}
}

func (s *competitionService) FetchAll(ctx context.Context, relations string) ([]dto.CompetitionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	rels := helpers.GetRelations(relations)

	competitions, err := s.competitionRepo.FetchAllByConditionAndRelation(ctx, "", nil, rels...)

	res := make([]dto.CompetitionResponse, 0)

	for _, c := range competitions {
		compeResp := dto.CompetitionEntityToResponse(&c)

		if c.Teams != nil {
			compeResp.Teams = dto.TeamSliceEntityToResponse(c.Teams)
		}

		if c.Announcements != nil {
			compeResp.Announcements = dto.AnnouncementSliceEntityToResponse(c.Announcements)
		}

		res = append(res, compeResp)
	}

	select {
	case <-ctx.Done():
		return nil, domain.ErrTimeout
	default:
		return res, err
	}
}

func (s *competitionService) FetchOne(ctx context.Context, id int, relations string) (dto.CompetitionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	rels := helpers.GetRelations(relations)

	competition, err := s.competitionRepo.FetchOneWithRelations(ctx, id, rels...)

	if err != nil {
		return dto.CompetitionResponse{}, err
	}

	res := dto.CompetitionEntityToResponse(&competition)

	if competition.Announcements != nil {
		res.Announcements = dto.AnnouncementSliceEntityToResponse(competition.Announcements)
	}

	if competition.Teams != nil {
		res.Teams = dto.TeamSliceEntityToResponse(competition.Teams)
	}

	select {
	case <-ctx.Done():
		return dto.CompetitionResponse{}, domain.ErrTimeout
	default:
		return res, err
	}
}

func (s *competitionService) InsertCompe(ctx context.Context, compe *dto.CompetitionRequest) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	payload := entity.NewCompetition(compe.Name, compe.Desc)

	err := s.competitionRepo.InsertCompe(ctx, payload)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *competitionService) UpdateCompe(ctx context.Context, compe *dto.CompetitionRequest) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	payload := entity.NewCompetition(compe.Name, compe.Desc)
	payload.LinkWA = compe.LinkWA
	payload.ID = compe.ID

	err := s.competitionRepo.UpdateCompe(ctx, payload)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *competitionService) DeleteCompe(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.competitionRepo.DeleteCompe(ctx, id)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

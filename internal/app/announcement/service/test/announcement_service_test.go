package test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	announcementMocks "github.com/hology8/hology-be/internal/app/announcement/repository/mock"
	"github.com/hology8/hology-be/internal/app/announcement/service"
	competitionMocks "github.com/hology8/hology-be/internal/app/competition/repository/mock"
	teamMocks "github.com/hology8/hology-be/internal/app/team/repository/mock"
	"github.com/stretchr/testify/assert"
)

func TestFetchAnnouncementByTo(t *testing.T) {
	type args struct {
		ctx           context.Context
		teamID        string
		competitionID int
	}

	dummyTeamID := "1"
	dummyCompetitionID := 1

	tests := []struct {
		name       string
		args       args
		want       interface{}
		wantErr    bool
		wantErrMsg string
		beforeTest func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository)
	}{
		{
			name: "When fetching announcements to all, it should not return error",
			args: args{context.TODO(), "", 0},
			want: []dto.AnnouncementResponse{
				{
					ID:            1,
					Description:   "Announcement 1",
					TeamID:        nil,
					CompetitionID: nil,
				},
				{
					ID:            2,
					Description:   "Announcement 2",
					TeamID:        nil,
					CompetitionID: nil,
				},
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {
				mockServiceRepo.EXPECT().
					FetchAnnouncementByTo(gomock.Any(), "", 0).
					Return([]entity.Announcement{
						{
							ID:            1,
							Description:   "Announcement 1",
							TeamID:        nil,
							CompetitionID: nil,
						},
						{
							ID:            2,
							Description:   "Announcement 2",
							TeamID:        nil,
							CompetitionID: nil,
						},
					}, nil)
			},
		},
		{
			name: "When fetching announcements to a team, it should not return error",
			args: args{context.TODO(), "1", 0},
			want: []dto.AnnouncementResponse{
				{
					ID:            1,
					Description:   "Announcement 1",
					TeamID:        &dummyTeamID,
					CompetitionID: nil,
				},
				{
					ID:            2,
					Description:   "Announcement 2",
					TeamID:        &dummyTeamID,
					CompetitionID: nil,
				},
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {

				mockServiceRepo.EXPECT().
					FetchAnnouncementByTo(gomock.Any(), "1", 0).
					Return([]entity.Announcement{
						{
							ID:            1,
							Description:   "Announcement 1",
							TeamID:        &dummyTeamID,
							CompetitionID: nil,
						},
						{
							ID:            2,
							Description:   "Announcement 2",
							TeamID:        &dummyTeamID,
							CompetitionID: nil,
						},
					}, nil)
			},
		},
		{
			name: "When fetching announcements to a competition, it should not return error",
			args: args{context.TODO(), "", 1},
			want: []dto.AnnouncementResponse{
				{
					ID:            1,
					Description:   "Announcement 1",
					TeamID:        nil,
					CompetitionID: &dummyCompetitionID,
				},
				{
					ID:            2,
					Description:   "Announcement 2",
					TeamID:        nil,
					CompetitionID: &dummyCompetitionID,
				},
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {

				mockServiceRepo.EXPECT().
					FetchAnnouncementByTo(gomock.Any(), "", 1).
					Return([]entity.Announcement{
						{
							ID:            1,
							Description:   "Announcement 1",
							TeamID:        nil,
							CompetitionID: &dummyCompetitionID,
						},
						{
							ID:            2,
							Description:   "Announcement 2",
							TeamID:        nil,
							CompetitionID: &dummyCompetitionID,
						},
					}, nil)
			},
		},
		{
			name:       "When fetching announcements to a team and competition, it should return error",
			args:       args{context.TODO(), "1", 1},
			want:       []dto.AnnouncementResponse{},
			wantErr:    true,
			wantErrMsg: domain.ErrIllegalEntry.Error(),
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {
			},
		},
		{
			name:       "When fetching announcements to a non existing team, it should return error",
			args:       args{context.TODO(), "1", 0},
			want:       []dto.AnnouncementResponse{},
			wantErr:    true,
			wantErrMsg: domain.ErrTeamNotFound.Error(),
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {

				mockServiceRepo.EXPECT().
					FetchAnnouncementByTo(gomock.Any(), "1", 0).
					Return([]entity.Announcement{}, domain.ErrTeamNotFound)
			},
		},
		{
			name:       "When fetching announcements to a non existing competition, it should return error",
			args:       args{context.TODO(), "", 1},
			want:       []dto.AnnouncementResponse{},
			wantErr:    true,
			wantErrMsg: domain.ErrCompetitionNotFound.Error(),
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {

				mockServiceRepo.EXPECT().
					FetchAnnouncementByTo(gomock.Any(), "", 1).
					Return([]entity.Announcement{}, domain.ErrCompetitionNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockAnnouncementRepo := announcementMocks.NewMockAnnouncementRepository(mockCtr)
			mockTeamRepo := teamMocks.NewMockTeamRepository(mockCtr)
			mockCompetitionRepo := competitionMocks.NewMockCompetitionRepository(mockCtr)

			test.beforeTest(mockAnnouncementRepo, mockTeamRepo, mockCompetitionRepo)

			announcementService := service.NewAnnouncementService(mockAnnouncementRepo, mockTeamRepo, mockCompetitionRepo, 500*time.Millisecond)

			announcements, err := announcementService.FetchAnnouncementByTo(test.args.ctx, test.args.teamID, test.args.competitionID)

			assert.Equal(t, test.want, announcements, "Expected announcements result to be match")
			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting error message to be same")
			} else {
				assert.Nil(t, err, "Expecting error to be nil")
			}
		})
	}
}

func TestCreateAnnouncement(t *testing.T) {
	type args struct {
		ctx          context.Context
		announcement dto.AnnouncementRequest
	}

	dummyTeamID := "1"
	dummyCompetitionID := 1

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		beforeTest func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository)
	}{
		{
			name: "When creating announcement to all, it should not return error",
			args: args{
				ctx: context.TODO(),
				announcement: dto.AnnouncementRequest{
					Description:   "Announcement 1",
					TeamID:        "",
					CompetitionID: 0,
				},
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {
				mockServiceRepo.EXPECT().
					InsertAnnouncement(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			name: "When creating announcement to a team, it should not return error",
			args: args{
				ctx: context.TODO(),
				announcement: dto.AnnouncementRequest{
					Description:   "Announcement 1",
					TeamID:        dummyTeamID,
					CompetitionID: 0,
				},
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {

				mockServiceRepo.EXPECT().
					InsertAnnouncement(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			name: "When creating announcement to a competition, it should not return error",
			args: args{
				ctx: context.TODO(),
				announcement: dto.AnnouncementRequest{
					Description:   "Announcement 1",
					TeamID:        "",
					CompetitionID: dummyCompetitionID,
				},
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {

				mockServiceRepo.EXPECT().
					InsertAnnouncement(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			name: "When creating announcement to a team and competition, it should return error",
			args: args{
				ctx: context.TODO(),
				announcement: dto.AnnouncementRequest{
					Description:   "Announcement 1",
					TeamID:        dummyTeamID,
					CompetitionID: dummyCompetitionID,
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrIllegalEntry.Error(),
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {
			},
		},
		{
			name: "When creating announcement to a non existing team, it should return error",
			args: args{
				ctx: context.TODO(),
				announcement: dto.AnnouncementRequest{
					Description:   "Announcement 1",
					TeamID:        dummyTeamID,
					CompetitionID: 0,
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrTeamNotFound.Error(),
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {

				mockServiceRepo.EXPECT().
					InsertAnnouncement(gomock.Any(), gomock.Any()).
					Return(domain.ErrTeamNotFound)
			},
		},
		{
			name: "When creating announcement to a non existing competition, it should return error",
			args: args{
				ctx: context.TODO(),
				announcement: dto.AnnouncementRequest{
					Description:   "Announcement 1",
					TeamID:        "",
					CompetitionID: dummyCompetitionID,
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrCompetitionNotFound.Error(),
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {

				mockServiceRepo.EXPECT().
					InsertAnnouncement(gomock.Any(), gomock.Any()).
					Return(domain.ErrCompetitionNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockAnnouncementRepo := announcementMocks.NewMockAnnouncementRepository(mockCtr)
			mockTeamRepo := teamMocks.NewMockTeamRepository(mockCtr)
			mockCompetitionRepo := competitionMocks.NewMockCompetitionRepository(mockCtr)

			test.beforeTest(mockAnnouncementRepo, mockTeamRepo, mockCompetitionRepo)

			announcementService := service.NewAnnouncementService(mockAnnouncementRepo, mockTeamRepo, mockCompetitionRepo, 500*time.Millisecond)

			err := announcementService.CreateAnnouncement(test.args.ctx, &test.args.announcement)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting error message to be same")
			} else {
				assert.Nil(t, err, "Expecting error to be nil")
			}
		})
	}
}

func TestUpdateAnnouncement(t *testing.T) {
	type args struct {
		ctx          context.Context
		announcement dto.AnnouncementRequest
	}

	dummyTeamID := "1"
	dummyCompetitionID := 1

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		beforeTest func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository)
	}{
		{
			name: "When updating announcement to all, it should not return error",
			args: args{
				ctx: context.TODO(),
				announcement: dto.AnnouncementRequest{
					ID:            1,
					Description:   "Announcement 1",
					TeamID:        "",
					CompetitionID: 0,
				},
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {
				mockServiceRepo.EXPECT().
					UpdateAnnouncement(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			name: "When updating announcement to a team, it should not return error",
			args: args{
				ctx: context.TODO(),
				announcement: dto.AnnouncementRequest{
					ID:            1,
					Description:   "Announcement 1",
					TeamID:        dummyTeamID,
					CompetitionID: 0,
				},
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {

				mockServiceRepo.EXPECT().
					UpdateAnnouncement(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			name: "When updating announcement to a competition, it should not return error",
			args: args{
				ctx: context.TODO(),
				announcement: dto.AnnouncementRequest{
					ID:            1,
					Description:   "Announcement 1",
					TeamID:        "",
					CompetitionID: dummyCompetitionID,
				},
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {

				mockServiceRepo.EXPECT().
					UpdateAnnouncement(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			name: "When updating announcement to a team and competition, it should return error",
			args: args{
				ctx: context.TODO(),
				announcement: dto.AnnouncementRequest{
					ID:            1,
					Description:   "Announcement 1",
					TeamID:        dummyTeamID,
					CompetitionID: dummyCompetitionID,
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrIllegalEntry.Error(),
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {
			},
		},
		{
			name: "When updating announcement to a non existing team, it should return error",
			args: args{
				ctx: context.TODO(),
				announcement: dto.AnnouncementRequest{
					ID:            1,
					Description:   "Announcement 1",
					TeamID:        dummyTeamID,
					CompetitionID: 0,
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrTeamNotFound.Error(),
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {

				mockServiceRepo.EXPECT().
					UpdateAnnouncement(gomock.Any(), gomock.Any()).
					Return(domain.ErrTeamNotFound)
			},
		},
		{
			name: "When updating announcement to a non existing competition, it should return error",
			args: args{
				ctx: context.TODO(),
				announcement: dto.AnnouncementRequest{
					ID:            1,
					Description:   "Announcement 1",
					TeamID:        "",
					CompetitionID: dummyCompetitionID,
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrCompetitionNotFound.Error(),
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository, mockTeamRepo *teamMocks.MockTeamRepository, mockCompetitionRepo *competitionMocks.MockCompetitionRepository) {

				mockServiceRepo.EXPECT().
					UpdateAnnouncement(gomock.Any(), gomock.Any()).
					Return(domain.ErrCompetitionNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockAnnouncementRepo := announcementMocks.NewMockAnnouncementRepository(mockCtr)
			mockTeamRepo := teamMocks.NewMockTeamRepository(mockCtr)
			mockCompetitionRepo := competitionMocks.NewMockCompetitionRepository(mockCtr)

			test.beforeTest(mockAnnouncementRepo, mockTeamRepo, mockCompetitionRepo)

			announcementService := service.NewAnnouncementService(mockAnnouncementRepo, mockTeamRepo, mockCompetitionRepo, 500*time.Millisecond)

			err := announcementService.UpdateAnnouncement(test.args.ctx, &test.args.announcement)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting error message to be same")
			} else {
				assert.Nil(t, err, "Expecting error to be nil")
			}
		})
	}
}

func TestDeleteAnnouncement(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		beforeTest func(mockServiceRepo *announcementMocks.MockAnnouncementRepository)
	}{
		{
			name:    "When deleting existing announcement, it should not return error",
			args:    args{context.TODO(), 1},
			wantErr: false,
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository) {
				mockServiceRepo.EXPECT().
					DeleteAnnouncement(gomock.Any(), 1).
					Return(nil)
			},
		},
		{
			name:       "When deleting non-existing announcement, it should return error",
			args:       args{context.TODO(), 1},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(mockServiceRepo *announcementMocks.MockAnnouncementRepository) {
				mockServiceRepo.EXPECT().
					DeleteAnnouncement(gomock.Any(), 1).
					Return(domain.ErrNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockAnnouncementRepo := announcementMocks.NewMockAnnouncementRepository(mockCtr)

			test.beforeTest(mockAnnouncementRepo)

			announcementService := service.NewAnnouncementService(mockAnnouncementRepo, nil, nil, 500*time.Millisecond)

			err := announcementService.DeleteAnnouncement(test.args.ctx, test.args.id)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting error message to be same")
			} else {
				assert.Nil(t, err, "Expecting error to be nil")
			}
		})
	}
}

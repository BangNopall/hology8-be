package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
	mocks "github.com/BangNopall/hology8-be/internal/app/competition/repository/mock"
	"github.com/BangNopall/hology8-be/internal/app/competition/service"
)

func TestFetchAll(t *testing.T) {
	tests := []struct {
		name       string
		want       interface{}
		wantErr    bool
		beforeTest func(mockCompeRepo *mocks.MockCompetitionRepository)
	}{
		{
			name: "When fetching all competitions, it should not return error",
			want: []dto.CompetitionResponse{
				{
					ID:   1,
					Name: "Competition 1",
					Desc: "This is cybersec Competition",
				},
				{
					ID:   2,
					Name: "Competition 2",
					Desc: "This is AI Competition",
				},
			},
			wantErr: false,
			beforeTest: func(mockCompeRepo *mocks.MockCompetitionRepository) {
				mockCompeRepo.EXPECT().
					FetchAllByConditionAndRelation(gomock.Any(), "", nil, gomock.Any()).
					Return([]entity.Competition{
						{
							ID:   1,
							Name: "Competition 1",
							Desc: "This is cybersec Competition",
						},
						{
							ID:   2,
							Name: "Competition 2",
							Desc: "This is AI Competition",
						},
					}, nil)
			},
		},
		{
			name:    "When fetching competition but operation time exceeded, it shoild return error timeout",
			want:    []dto.CompetitionResponse(nil),
			wantErr: true,
			beforeTest: func(mockCompeRepo *mocks.MockCompetitionRepository) {
				mockCompeRepo.EXPECT().
					FetchAllByConditionAndRelation(gomock.Any(), "", nil, gomock.Any()).
					DoAndReturn(func(ctx context.Context, condition string, args []interface{}, preload ...string) ([]entity.Competition, error) {
						time.Sleep(1 * time.Second)
						return []entity.Competition(nil), nil
					})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockCompeRepo := mocks.NewMockCompetitionRepository(mockCtr)

			test.beforeTest(mockCompeRepo)

			compeSvc := service.NewCompetitionService(mockCompeRepo, 500*time.Millisecond)

			competitions, err := compeSvc.FetchAll(context.TODO(), "")

			assert.Equal(t, test.want, competitions, "Expected competitions result to be match")
			if test.wantErr {
				assert.NotNil(t, err, "Error expected to be thrown")
			} else {
				assert.Nil(t, err, "Expected error to be nil")
			}
		})
	}
}
func TestFetchOne(t *testing.T) {
	type args struct {
		ctx            context.Context
		id             int
		relations      string
		relationSpread []interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		want       interface{}
		beforeTest func(args args, mockCompeRepo *mocks.MockCompetitionRepository)
	}{
		{
			name: "When fetching competition by id, it should not return eror",
			args: args{
				context.TODO(),
				1,
				"",
				[]interface{}{""},
			},
			want: dto.CompetitionResponse{
				ID:   1,
				Name: "CTF",
				Desc: "Ini Lomba CTF",
			},
			wantErr: false,
			beforeTest: func(args args, mockCompeRepo *mocks.MockCompetitionRepository) {
				mockCompeRepo.EXPECT().
					FetchOneWithRelations(gomock.Any(), gomock.Any(), args.relationSpread...).
					Return(entity.Competition{
						ID:   1,
						Name: "CTF",
						Desc: "Ini Lomba CTF",
					}, nil)
			},
		},
		{
			name: "When fetching non existing competition by id, it should return eror",
			args: args{
				context.TODO(),
				1,
				"",
				[]interface{}{""},
			},
			wantErr:    true,
			want:       dto.CompetitionResponse{},
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(args args, mockCompeRepo *mocks.MockCompetitionRepository) {
				mockCompeRepo.EXPECT().
					FetchOneWithRelations(gomock.Any(), gomock.Any(), args.relationSpread...).
					Return(entity.Competition{}, domain.ErrNotFound)
			},
		},
		{
			name: "When fetching competition by id with team relation, it should return compe with team relation preloaded",
			args: args{
				context.TODO(),
				1,
				"team",
				[]interface{}{"Teams"},
			},
			wantErr: false,
			want: dto.CompetitionResponse{
				ID:   1,
				Name: "CTF",
				Desc: "Ini Lomba CTF",
				Teams: []dto.TeamResponse{
					{
						ID:      "1",
						Name:    "AcRtf",
						Members: []dto.UserResponse{},
						University: dto.UniversityResponse{
							ID:   1,
							Name: "ASASAS",
						},
					},
				},
			},
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(args args, mockCompeRepo *mocks.MockCompetitionRepository) {
				mockCompeRepo.EXPECT().
					FetchOneWithRelations(gomock.Any(), gomock.Any(), args.relationSpread...).
					Return(entity.Competition{
						ID:   1,
						Name: "CTF",
						Desc: "Ini Lomba CTF",
						Teams: []entity.Team{
							{
								ID:   "1",
								Name: "AcRtf",
								University: entity.University{
									ID:   1,
									Name: "ASASAS",
								},
							},
						},
					}, nil)
			},
		},
		{
			name: "When fetching competition by id with announcement relation, it should return compe with announcement relation preloaded",
			args: args{
				context.TODO(),
				1,
				"announcement",
				[]interface{}{"Announcements"},
			},
			wantErr: false,
			want: dto.CompetitionResponse{
				ID:   1,
				Name: "CTF",
				Desc: "Ini Lomba CTF",
				Announcements: []dto.AnnouncementResponse{
					{
						ID:          1,
						Description: "woi",
					},
				},
			},
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(args args, mockCompeRepo *mocks.MockCompetitionRepository) {
				mockCompeRepo.EXPECT().
					FetchOneWithRelations(gomock.Any(), gomock.Any(), args.relationSpread...).
					Return(entity.Competition{
						ID:   1,
						Name: "CTF",
						Desc: "Ini Lomba CTF",
						Announcements: []entity.Announcement{
							{
								ID:          1,
								Description: "woi",
							},
						},
					}, nil)
			},
		},
		{
			name: "When fetching competition by id with team and announcement relation, it should return compe with team and announcement relation preloaded",
			args: args{
				context.TODO(),
				1,
				"team,announcement",
				[]interface{}{"Teams", "Announcements"},
			},
			wantErr: false,
			want: dto.CompetitionResponse{
				ID:   1,
				Name: "CTF",
				Desc: "Ini Lomba CTF",
				Announcements: []dto.AnnouncementResponse{
					{
						ID:          1,
						Description: "woi",
					},
				},
				Teams: []dto.TeamResponse{
					{
						ID:      "1",
						Name:    "AcRtf",
						Members: []dto.UserResponse{},
						University: dto.UniversityResponse{
							ID:   1,
							Name: "ASASAS",
						},
					},
				},
			},
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(args args, mockCompeRepo *mocks.MockCompetitionRepository) {
				mockCompeRepo.EXPECT().
					FetchOneWithRelations(gomock.Any(), gomock.Any(), gomock.All()).
					Return(entity.Competition{
						ID:   1,
						Name: "CTF",
						Desc: "Ini Lomba CTF",
						Announcements: []entity.Announcement{
							{
								ID:          1,
								Description: "woi",
							},
						},
						Teams: []entity.Team{
							{
								ID:   "1",
								Name: "AcRtf",
								University: entity.University{
									ID:   1,
									Name: "ASASAS",
								},
							},
						},
					}, nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockCompeRepo := mocks.NewMockCompetitionRepository(mockCtr)

			test.beforeTest(test.args, mockCompeRepo)

			compSvc := service.NewCompetitionService(mockCompeRepo, time.Second*1)

			compe, err := compSvc.FetchOne(test.args.ctx, test.args.id, test.args.relations)

			if !test.wantErr {
				assert.Nil(t, err, "Error should not be expected")
			} else {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error message")
			}

			assert.Equal(t, test.want, compe, "Expecting same compe result")
		})
	}
}
func TestInsertCompe(t *testing.T) {
	type args struct {
		ctx   context.Context
		compe *dto.CompetitionRequest
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeTest func(mockCompeRepo *mocks.MockCompetitionRepository)
	}{
		{
			name: "When creating compe, it should not return eror",
			args: args{
				context.TODO(),
				&dto.CompetitionRequest{
					Name: "Competition 1",
					Desc: "This is cybersec compe",
				},
			},
			wantErr: false,
			beforeTest: func(mockCompeRepo *mocks.MockCompetitionRepository) {
				mockCompeRepo.EXPECT().
					InsertCompe(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockCompeRepo := mocks.NewMockCompetitionRepository(mockCtr)

			test.beforeTest(mockCompeRepo)

			compeSvc := service.NewCompetitionService(mockCompeRepo, time.Second*2)

			err := compeSvc.InsertCompe(test.args.ctx, test.args.compe)

			if !test.wantErr {
				assert.Nil(t, err, "Error should not be expected")
			} else {
				assert.NotNil(t, err, "Expecting error to be thrown")
			}
		})
	}
}
func TestUpdateCompe(t *testing.T) {
	type args struct {
		ctx   context.Context
		compe *dto.CompetitionRequest
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		beforeTest func(mockCompeRepo *mocks.MockCompetitionRepository)
	}{
		{
			name: "When updating existing compe, it should not return error",
			args: args{
				context.TODO(),
				&dto.CompetitionRequest{
					ID:   1,
					Name: "Competition 1",
					Desc: "This is cybersec compe",
				},
			},
			wantErr: false,
			beforeTest: func(mockCompeRepo *mocks.MockCompetitionRepository) {
				mockCompeRepo.EXPECT().
					UpdateCompe(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			name: "When updating non-existing compe, it should return error item not found",
			args: args{
				context.TODO(),
				&dto.CompetitionRequest{
					ID:   1,
					Name: "Competition 1",
					Desc: "This is cybersec compe",
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(mockCompeRepo *mocks.MockCompetitionRepository) {
				mockCompeRepo.EXPECT().
					UpdateCompe(gomock.Any(), gomock.Any()).
					Return(domain.ErrNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockCompeRepo := mocks.NewMockCompetitionRepository(mockCtr)

			test.beforeTest(mockCompeRepo)

			compeSvc := service.NewCompetitionService(mockCompeRepo, time.Second*2)

			err := compeSvc.UpdateCompe(test.args.ctx, test.args.compe)

			if !test.wantErr {
				assert.Nil(t, err, "Error should not be expected")
			} else {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error message")
			}
		})
	}
}
func TestDeleteCompe(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		beforeTest func(mockCompeRepo *mocks.MockCompetitionRepository)
	}{
		{
			name: "When deleting existing competition, it should not return error",
			args: args{
				context.TODO(),
				1,
			},
			wantErr: false,
			beforeTest: func(mockCompeRepo *mocks.MockCompetitionRepository) {
				mockCompeRepo.EXPECT().
					DeleteCompe(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			name: "When deleting non-existing competition, it should return error item not found",
			args: args{
				context.TODO(),
				1,
			},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(mockCompeRepo *mocks.MockCompetitionRepository) {
				mockCompeRepo.EXPECT().
					DeleteCompe(gomock.Any(), gomock.Any()).
					Return(domain.ErrNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockCompeRepo := mocks.NewMockCompetitionRepository(mockCtr)

			test.beforeTest(mockCompeRepo)

			compeSvc := service.NewCompetitionService(mockCompeRepo, time.Second*2)

			err := compeSvc.DeleteCompe(test.args.ctx, test.args.id)

			if !test.wantErr {
				assert.Nil(t, err, "Error should not be expected")
			} else {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error message")
			}
		})
	}
}

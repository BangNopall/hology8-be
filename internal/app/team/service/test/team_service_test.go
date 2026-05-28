package service_test

import (
	"context"
	"mime/multipart"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
	"github.com/BangNopall/hology8-be/domain/enums"
	mockCompetitionRepo "github.com/BangNopall/hology8-be/internal/app/competition/repository/mock"
	mockTeamRepo "github.com/BangNopall/hology8-be/internal/app/team/repository/mock"
	"github.com/BangNopall/hology8-be/internal/app/team/service"
	mockUserRepo "github.com/BangNopall/hology8-be/internal/app/user/repository/mock"
	mockAws "github.com/BangNopall/hology8-be/pkg/aws/mock"
)

type mockObjects struct {
	teamRepo        *mockTeamRepo.MockTeamRepository
	competitionRepo *mockCompetitionRepo.MockCompetitionRepository
	userRepo        *mockUserRepo.MockUserRepository
}

func TestFetchTeamData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	teamRepo := mockTeamRepo.NewMockTeamRepository(ctrl)

	mockObj := mockObjects{
		teamRepo: teamRepo,
	}

	type args struct {
		teamId string
	}

	tests := []struct {
		name        string
		args        args
		beforeTest  func(mockObj mockObjects)
		want        interface{}
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when fetching existing team data, it should return the data without error",
			args: args{
				teamId: "hehe",
			},
			beforeTest: func(mockObj mockObjects) {
				mockObj.teamRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), "id = ?", []interface{}{"hehe"}, gomock.Any(), nil, gomock.All()).Return([]entity.Team{
					{
						ID:            "hehe",
						LeaderID:      uuid.MustParse("26474a36-e0e1-4ac0-aec0-23d7e760295a"),
						Name:          "ini nama teamnya",
						CompetitionID: 1,
						Leader: entity.User{
							Fullname: "Devan",
						},
						University: entity.University{
							ID:   1,
							Name: "ASASAS",
						},
						Members: []entity.DetailTeams{},
					},
				}, dto.PaginationResponse{}, nil)
			},
			want: dto.TeamResponse{
				ID:            "hehe",
				LeaderID:      uuid.MustParse("26474a36-e0e1-4ac0-aec0-23d7e760295a"),
				Name:          "ini nama teamnya",
				CompetitionID: 1,
				UniversityID:  0,
				Leader: dto.UserResponse{
					Fullname: "Devan",
				},
				Competition: dto.CompetitionResponse{
					Announcements: []dto.AnnouncementResponse{},
				},
				University: dto.UniversityResponse{
					ID:   1,
					Name: "ASASAS",
				},
				Announcements: []dto.AnnouncementResponse{},
				Members:       []dto.UserResponse{},
			},
		},
		{
			name: "when fetching non-existing team data, it should return error not found",
			args: args{
				teamId: "hehe",
			},
			beforeTest: func(mockObj mockObjects) {
				mockObj.teamRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), "id = ?", []interface{}{"hehe"}, gomock.Any(), nil, gomock.All()).Return([]entity.Team{}, dto.PaginationResponse{}, nil)
			},
			wantErr:     true,
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.beforeTest(mockObj)

			s := service.NewTeamService(mockObj.teamRepo, nil, nil, time.Millisecond*500, nil)

			res, err := s.FetchTeamData(context.Background(), test.args.teamId)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.expectedErr, err, "Expecting error to be %v", test.expectedErr)
			} else {
				assert.Nil(t, err, "Expecting no error")
				assert.Equal(t, test.want, res, "Expecting result to be equal")
			}
		})
	}
}

func TestFetchUserTeams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	teamRepo := mockTeamRepo.NewMockTeamRepository(ctrl)

	mockObj := mockObjects{
		teamRepo: teamRepo,
	}

	type args struct {
		userId uuid.UUID
	}

	tests := []struct {
		name        string
		args        args
		beforeTest  func(mockObj mockObjects)
		want        interface{}
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when fetching user teams, it should return all user's team has joined",
			args: args{
				userId: uuid.MustParse("26474a36-e0e1-4ac0-aec0-23d7e760295a"),
			},
			beforeTest: func(mockObj mockObjects) {
				mockObj.teamRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), "leader_id = ?",
					[]interface{}{uuid.MustParse("26474a36-e0e1-4ac0-aec0-23d7e760295a")}, gomock.Any(), nil, "Leader", "Competition", "Announcements").
					Return([]entity.Team{
						{
							ID:            "hehe",
							LeaderID:      uuid.MustParse("26474a36-e0e1-4ac0-aec0-23d7e760295a"),
							Name:          "ini nama teamnya",
							CompetitionID: 2,
							Members:       []entity.DetailTeams{},
							University: entity.University{
								ID:   1,
								Name: "ASASAS",
							},
							Announcements: []entity.Announcement{
								{
									ID:          1,
									Description: "awawaw",
								},
							},
						},
					}, dto.PaginationResponse{}, nil)

				mockObj.teamRepo.EXPECT().FetchMemberTeams(gomock.Any(), uuid.MustParse("26474a36-e0e1-4ac0-aec0-23d7e760295a")).
					Return([]entity.DetailTeams{
						{
							Team: entity.Team{
								ID:            "hehe",
								LeaderID:      uuid.MustParse("26474a36-e0e1-4ac0-aec0-23d7e760295a"),
								Name:          "ini nama teamnya",
								CompetitionID: 1,
								University: entity.University{
									ID:   1,
									Name: "ASASAS",
								},
								Announcements: []entity.Announcement{
									{
										ID:          1,
										Description: "awawaw",
									},
								},
							},
						},
					}, nil)
			},
			want: dto.UserTeamsResponse{
				Teams: []dto.TeamResponse{
					{
						ID:            "hehe",
						LeaderID:      uuid.MustParse("26474a36-e0e1-4ac0-aec0-23d7e760295a"),
						Name:          "ini nama teamnya",
						CompetitionID: 2,
						Members:       []dto.UserResponse{},
						University: dto.UniversityResponse{
							ID:   1,
							Name: "ASASAS",
						},
						Announcements: []dto.AnnouncementResponse{
							{
								ID:          1,
								Description: "awawaw",
							},
						},
					},
					{
						ID:            "hehe",
						LeaderID:      uuid.MustParse("26474a36-e0e1-4ac0-aec0-23d7e760295a"),
						Name:          "ini nama teamnya",
						CompetitionID: 1,
						Members:       []dto.UserResponse{},
						University: dto.UniversityResponse{
							ID:   1,
							Name: "ASASAS",
						},
						Announcements: []dto.AnnouncementResponse{
							{
								ID:          1,
								Description: "awawaw",
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.beforeTest(mockObj)

			s := service.NewTeamService(mockObj.teamRepo, nil, nil, time.Millisecond*500, nil)

			res, err := s.FetchUserTeams(context.Background(), test.args.userId)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.expectedErr, err, "Expecting error to be %v", test.expectedErr)
			} else {
				assert.Nil(t, err, "Expecting no error")
				assert.Equal(t, test.want, res, "Expecting result to be equal")
			}
		})
	}
}

func TestCreateTeam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	teamRepo := mockTeamRepo.NewMockTeamRepository(ctrl)
	competitionRepo := mockCompetitionRepo.NewMockCompetitionRepository(ctrl)
	userRepo := mockUserRepo.NewMockUserRepository(ctrl)
	storage := mockAws.NewMockCloudStorage(ctrl)

	mockObject := mockObjects{
		teamRepo:        teamRepo,
		competitionRepo: competitionRepo,
		userRepo:        userRepo,
	}

	type args struct {
		leaderID uuid.UUID
		team     dto.TeamRegister
	}

	dummyUUID := uuid.New()
	dummyCompetition := entity.Competition{
		ID:   1,
		Name: "TestCompetition",
		Desc: "TestDesc",
	}

	mockArgs := args{
		leaderID: dummyUUID,
		team: dto.TeamRegister{
			CompetitionID: 1,
			Name:          "TestTeam",
		},
	}

	tests := []struct {
		name          string
		args          args
		beforeTest    func(mockObject mockObjects)
		wantErr       bool
		expectedError error
	}{
		{
			name: "when create team, it should success without error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.competitionRepo.EXPECT().FetchOneByID(gomock.Any(), gomock.Any()).Return(dummyCompetition, nil)
				mockObject.teamRepo.EXPECT().FetchAllTeams(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, domain.ErrNotFound)
				mockObject.userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(nil)
				mockObject.teamRepo.EXPECT().InsertTeam(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "when create team but it get error when fetch competition by id, it should return internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.competitionRepo.EXPECT().FetchOneByID(gomock.Any(), gomock.Any()).Return(entity.Competition{}, domain.ErrInternalServer)
			},
			wantErr:       true,
			expectedError: domain.ErrInternalServer,
		},
		{
			name: "when create team but competition not found, it should return not found error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.competitionRepo.EXPECT().FetchOneByID(gomock.Any(), gomock.Any()).Return(entity.Competition{}, domain.ErrNotFound)
			},
			wantErr:       true,
			expectedError: domain.ErrNotFound,
		},
		{
			name: "when create team but it get error when fetch all teams, it should return internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.competitionRepo.EXPECT().FetchOneByID(gomock.Any(), gomock.Any()).Return(dummyCompetition, nil)
				mockObject.teamRepo.EXPECT().FetchAllTeams(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, domain.ErrInternalServer)
			},
			wantErr:       true,
			expectedError: domain.ErrInternalServer,
		},
		{
			name: "when create team but it get error because leader already register in the same competition, it should return error user already registered",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.competitionRepo.EXPECT().FetchOneByID(gomock.Any(), gomock.Any()).Return(dummyCompetition, nil)
				mockObject.teamRepo.EXPECT().FetchAllTeams(gomock.Any(), gomock.Any(), gomock.Any()).Return([]entity.Team{{Members: []entity.DetailTeams{{UserID: mockArgs.leaderID}}}}, nil)
			},
			wantErr:       true,
			expectedError: domain.ErrUserAlreadyRegistered,
		},
		{
			name: "when create team but it get error when insert team, it should return internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.competitionRepo.EXPECT().FetchOneByID(gomock.Any(), gomock.Any()).Return(dummyCompetition, nil)
				mockObject.teamRepo.EXPECT().FetchAllTeams(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, domain.ErrNotFound)
				mockObject.userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(nil)
				mockObject.teamRepo.EXPECT().InsertTeam(gomock.Any(), gomock.Any()).Return(domain.ErrInternalServer)
			},
			wantErr:       true,
			expectedError: domain.ErrInternalServer,
		},
		{
			name: "when create team but operation time exceeded, it should failed with error time out",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.competitionRepo.EXPECT().FetchOneByID(gomock.Any(), gomock.Any()).Return(dummyCompetition, nil)
				mockObject.teamRepo.EXPECT().FetchAllTeams(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, domain.ErrNotFound)
				mockObject.userRepo.EXPECT().FindUser(gomock.Any(), gomock.Any()).Return(nil)
				mockObject.teamRepo.EXPECT().InsertTeam(gomock.Any(), gomock.Any()).
					DoAndReturn(func(any any, any2 any) error {
						time.Sleep(time.Millisecond * 1000)
						return nil
					})
			},
			wantErr:       true,
			expectedError: domain.ErrTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeTest(mockObject)

			s := service.NewTeamService(mockObject.teamRepo, mockObject.competitionRepo, mockObject.userRepo, time.Millisecond*500, storage)

			err := s.CreateTeam(context.Background(), tt.args.leaderID, tt.args.team)

			if tt.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, tt.expectedError, err, "Expecting error to be %v", tt.expectedError)
			} else {
				assert.Nil(t, err, "Expecting no error")
			}
		})
	}
}

func TestUploadPaymentProof(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     string
		userId string
		file   *multipart.FileHeader
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		beforeTest func(
			args args,
			mockTeamRepo *mockTeamRepo.MockTeamRepository,
			mockAws *mockAws.MockCloudStorage,
		)
	}{
		{
			name: "When uploading a payment proof, it should return no error",
			args: args{
				context.TODO(),
				"TeamID",
				"26474a36-e0e1-4ac0-aec0-23d7e760295a",
				&multipart.FileHeader{
					Filename: "ini file.png",
					Size:     999,
				},
			},
			wantErr: false,
			beforeTest: func(
				args args,
				mockTeamRepo *mockTeamRepo.MockTeamRepository,
				mockAws *mockAws.MockCloudStorage) {
				mockTeamRepo.EXPECT().
					FetchOneByParams(gomock.Any(), &dto.TeamParams{ID: "TeamID"}, gomock.All()).
					Return(entity.Team{
						ID:               "TeamID",
						LeaderID:         uuid.MustParse("26474a36-e0e1-4ac0-aec0-23d7e760295a"),
						TwibbonProofLink: "",
						PaymentProofLink: "",
					}, nil)

				mockAws.EXPECT().
					Upload(gomock.Any(), args.file).
					Return("link", nil)

				mockTeamRepo.EXPECT().
					UpdateTeam(gomock.Any(), args.id, &dto.TeamUpdate{PaymentProofLink: "link"}).
					Return(nil)
			},
		},
		{
			name: "When uploading a payment proof but the file exceeds 1MB, it should return error file too big",
			args: args{
				context.TODO(),
				"TeamID",
				"26474a36-e0e1-4ac0-aec0-23d7e760295a",
				&multipart.FileHeader{
					Filename: "ini file.png",
					Size:     1024*1024 + 500,
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrFileTooBig.Error(),
		},
		{
			name: "When uploading non existing team's payment proof, it should return error not found",
			args: args{
				context.TODO(),
				"TeamID",
				"26474a36-e0e1-4ac0-aec0-23d7e760295a",
				&multipart.FileHeader{
					Filename: "ini file.png",
					Size:     1024 * 1022,
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(
				args args,
				mockTeamRepo *mockTeamRepo.MockTeamRepository,
				mockAws *mockAws.MockCloudStorage) {
				mockTeamRepo.EXPECT().
					FetchOneByParams(gomock.Any(), &dto.TeamParams{ID: "TeamID"}, gomock.All()).
					Return(entity.Team{}, domain.ErrNotFound)
			},
		},
		{
			name: "When updating team's payment proof, it should return no error",
			args: args{
				context.TODO(),
				"TeamID",
				"26474a36-e0e1-4ac0-aec0-23d7e760295a",
				&multipart.FileHeader{
					Filename: "ini file.png",
					Size:     1024 * 1022,
				},
			},
			wantErr: false,
			beforeTest: func(
				args args,
				mockTeamRepo *mockTeamRepo.MockTeamRepository,
				mockAws *mockAws.MockCloudStorage) {
				mockTeamRepo.EXPECT().
					FetchOneByParams(gomock.Any(), &dto.TeamParams{ID: "TeamID"}, gomock.All()).
					Return(entity.Team{
						ID:               "TeamID",
						LeaderID:         uuid.MustParse("26474a36-e0e1-4ac0-aec0-23d7e760295a"),
						TwibbonProofLink: "",
						PaymentProofLink: "payment1.com",
					}, nil)

				mockAws.EXPECT().
					Update(gomock.Any(), args.file, "payment1.com").
					Return("newLink", nil)

				mockTeamRepo.EXPECT().
					UpdateTeam(gomock.Any(), args.id, &dto.TeamUpdate{PaymentProofLink: "newLink"}).
					Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)

			defer mockCtr.Finish()

			repoMock := mockTeamRepo.NewMockTeamRepository(mockCtr)
			awsMock := mockAws.NewMockCloudStorage(mockCtr)

			if test.beforeTest != nil {
				test.beforeTest(test.args, repoMock, awsMock)
			}

			svc := service.NewTeamService(repoMock, nil, nil, time.Second, awsMock)

			err := svc.UploadPaymentProof(test.args.ctx, test.args.id, test.args.userId, test.args.file)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error result")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}
		})
	}
}

func TestUploadTwibbonProof(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     string
		userId string
		file   *multipart.FileHeader
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		beforeTest func(
			args args,
			mockTeamRepo *mockTeamRepo.MockTeamRepository,
			mockAws *mockAws.MockCloudStorage,
		)
	}{
		{
			name: "When uploading a twibbon proof, it should return no error",
			args: args{
				context.TODO(),
				"TeamID",
				"26474a36-e0e1-4ac0-aec0-23d7e760295a",
				&multipart.FileHeader{
					Filename: "ini file.png",
					Size:     999,
				},
			},
			wantErr: false,
			beforeTest: func(
				args args,
				mockTeamRepo *mockTeamRepo.MockTeamRepository,
				mockAws *mockAws.MockCloudStorage) {
				mockTeamRepo.EXPECT().
					FetchOneByParams(gomock.Any(), &dto.TeamParams{ID: "TeamID"}, gomock.All()).
					Return(entity.Team{
						ID:               "TeamID",
						LeaderID:         uuid.MustParse("26474a36-e0e1-4ac0-aec0-23d7e760295a"),
						TwibbonProofLink: "",
						PaymentProofLink: "",
					}, nil)

				mockAws.EXPECT().
					Upload(gomock.Any(), args.file).
					Return("link", nil)

				mockTeamRepo.EXPECT().
					UpdateTeam(gomock.Any(), args.id, &dto.TeamUpdate{TwibbonProofLink: "link"}).
					Return(nil)
			},
		},
		{
			name: "When uploading a twibbon proof but the file exceeds 1MB, it should return error file too big",
			args: args{
				context.TODO(),
				"TeamID",
				"26474a36-e0e1-4ac0-aec0-23d7e760295a",
				&multipart.FileHeader{
					Filename: "ini file.png",
					Size:     1024*1024 + 500,
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrFileTooBig.Error(),
		},
		{
			name: "When uploading non existing team's payment proof, it should return error not found",
			args: args{
				context.TODO(),
				"TeamID",
				"26474a36-e0e1-4ac0-aec0-23d7e760295a",
				&multipart.FileHeader{
					Filename: "ini file.png",
					Size:     1024 * 1022,
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(
				args args,
				mockTeamRepo *mockTeamRepo.MockTeamRepository,
				mockAws *mockAws.MockCloudStorage) {
				mockTeamRepo.EXPECT().
					FetchOneByParams(gomock.Any(), &dto.TeamParams{ID: "TeamID"}, gomock.All()).
					Return(entity.Team{}, domain.ErrNotFound)
			},
		},
		{
			name: "When updating team's twibbon proof, it should return no error",
			args: args{
				context.TODO(),
				"TeamID",
				"26474a36-e0e1-4ac0-aec0-23d7e760295a",
				&multipart.FileHeader{
					Filename: "ini file.png",
					Size:     1024 * 1022,
				},
			},
			wantErr: false,
			beforeTest: func(
				args args,
				mockTeamRepo *mockTeamRepo.MockTeamRepository,
				mockAws *mockAws.MockCloudStorage) {
				mockTeamRepo.EXPECT().
					FetchOneByParams(gomock.Any(), &dto.TeamParams{ID: "TeamID"}, gomock.All()).
					Return(entity.Team{
						ID:               "TeamID",
						LeaderID:         uuid.MustParse("26474a36-e0e1-4ac0-aec0-23d7e760295a"),
						TwibbonProofLink: "payment1.com",
						PaymentProofLink: "",
					}, nil)

				mockAws.EXPECT().
					Update(gomock.Any(), args.file, "payment1.com").
					Return("newLink", nil)

				mockTeamRepo.EXPECT().
					UpdateTeam(gomock.Any(), args.id, &dto.TeamUpdate{TwibbonProofLink: "newLink"}).
					Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)

			defer mockCtr.Finish()

			repoMock := mockTeamRepo.NewMockTeamRepository(mockCtr)
			awsMock := mockAws.NewMockCloudStorage(mockCtr)

			if test.beforeTest != nil {
				test.beforeTest(test.args, repoMock, awsMock)
			}

			svc := service.NewTeamService(repoMock, nil, nil, time.Second, awsMock)

			err := svc.UploadTwibbonProof(test.args.ctx, test.args.id, test.args.userId, test.args.file)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error result")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}
		})
	}
}

func TestUpdateTeamData(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     string
		userId string
		team   *dto.TeamUpdate
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		beforeTest func(
			args args,
			mockTeamRepo *mockTeamRepo.MockTeamRepository,
		)
	}{
		{
			name: "When updating a team data, it should not return error",
			args: args{
				context.TODO(),
				"77FA",
				"ad172b18-f829-400f-8d75-cb16538c8503",
				&dto.TeamUpdate{
					UniversityID:      1,
					SenderPaymentName: "Devan",
					BankAccountNumber: "77",
					BankName:          "Mandiri",
					ProposalDocLink:   "link.com",
					VideoLink:         "link.com",
				},
			},
			wantErr: false,
			beforeTest: func(args args, mockTeamRepo *mockTeamRepo.MockTeamRepository) {
				mockTeamRepo.EXPECT().
					FetchOneByID(gomock.Any(), args.id).
					Return(entity.Team{
						LeaderID:          uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
						SenderPaymentName: "Devan",
						BankAccountNumber: "77",
						BankName:          "Mandiri",
						ProposalDocLink:   "link.com",
						VideoLink:         "link.com",
					}, nil)

				mockTeamRepo.EXPECT().
					UpdateTeam(gomock.Any(), args.id, args.team).
					Return(nil)
			},
		},
		{
			name: "When updating a attribute that is illegal to be updated in this service, it should return error forbidden update",
			args: args{
				context.TODO(),
				"77FA",
				"ad172b18-f829-400f-8d75-cb16538c8503",
				&dto.TeamUpdate{
					UniversityID:      1,
					SenderPaymentName: "Devan",
					BankAccountNumber: "77",
					BankName:          "Mandiri",
					ProposalDocLink:   "link.com",
					VideoLink:         "link.com",
					Status:            enums.Verified,
					Phase:             enums.Elimination,
				},
			},

			wantErr:    true,
			wantErrMsg: domain.ErrForbiddenUpdate.Error(),
			beforeTest: func(args args, mockTeamRepo *mockTeamRepo.MockTeamRepository) {
			},
		},
		{
			name: "When updating a attribute that is illegal to be updated in this service, it should return error forbidden update",
			args: args{
				context.TODO(),
				"77FA",
				"ad172b18-f829-400f-8d75-cb16538c8503",
				&dto.TeamUpdate{
					UniversityID:      1,
					SenderPaymentName: "Devan",
					BankAccountNumber: "77",
					BankName:          "Mandiri",
					ProposalDocLink:   "link.com",
					VideoLink:         "link.com",
					Status:            enums.Verified,
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrForbiddenUpdate.Error(),
			beforeTest: func(args args, mockTeamRepo *mockTeamRepo.MockTeamRepository) {
			},
		},
		{
			name: "When updating a attribute that is illegal to be updated in this service, it should return error forbidden update",
			args: args{
				context.TODO(),
				"77FA",
				"ad172b18-f829-400f-8d75-cb16538c8503",
				&dto.TeamUpdate{
					UniversityID:      1,
					SenderPaymentName: "Devan",
					BankAccountNumber: "77",
					BankName:          "Mandiri",
					ProposalDocLink:   "link.com",
					VideoLink:         "link.com",
					Phase:             enums.Elimination,
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrForbiddenUpdate.Error(),
			beforeTest: func(args args, mockTeamRepo *mockTeamRepo.MockTeamRepository) {
			},
		},
		{
			name: "When updating a attribute that is illegal to be updated in this service, it should return error forbidden update",
			args: args{
				context.TODO(),
				"77FA",
				"ad172b18-f829-400f-8d75-cb16538c8503",
				&dto.TeamUpdate{
					LeaderID: uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
				},
			},

			wantErr:    true,
			wantErrMsg: domain.ErrForbiddenUpdate.Error(),
			beforeTest: func(args args, mockTeamRepo *mockTeamRepo.MockTeamRepository) {
			},
		},
		{
			name: "When updating a non-existing team status, it should return error item not found",
			args: args{
				context.TODO(),
				"77FA",
				"ad172b18-f829-400f-8d75-cb16538c8503",
				&dto.TeamUpdate{
					UniversityID:      1,
					SenderPaymentName: "Devan",
					BankAccountNumber: "77",
					BankName:          "Mandiri",
					ProposalDocLink:   "link.com",
					VideoLink:         "link.com",
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(args args, mockTeamRepo *mockTeamRepo.MockTeamRepository) {
				mockTeamRepo.EXPECT().
					FetchOneByID(gomock.Any(), args.id).
					Return(entity.Team{
						LeaderID:          uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
						SenderPaymentName: "Devan",
						BankAccountNumber: "77",
						BankName:          "Mandiri",
						ProposalDocLink:   "link.com",
						VideoLink:         "link.com",
					}, nil)

				mockTeamRepo.EXPECT().
					UpdateTeam(gomock.Any(), args.id, args.team).
					Return(domain.ErrNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)

			defer mockCtr.Finish()

			repoMock := mockTeamRepo.NewMockTeamRepository(mockCtr)
			awsMock := mockAws.NewMockCloudStorage(mockCtr)

			if test.beforeTest != nil {
				test.beforeTest(test.args, repoMock)
			}

			svc := service.NewTeamService(repoMock, nil, nil, time.Second, awsMock)

			err := svc.UpdateTeamData(test.args.ctx, test.args.id, test.args.userId, test.args.team)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error message")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}
		})
	}
}

func TestUpdateTeamStatus(t *testing.T) {
	type args struct {
		ctx  context.Context
		id   string
		team *dto.TeamUpdate
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		beforeTest func(
			args args,
			mockTeamRepo *mockTeamRepo.MockTeamRepository,
		)
	}{
		{
			name: "When updating a team status, it should not return error",
			args: args{
				context.TODO(),
				"77FA",
				&dto.TeamUpdate{
					Status: enums.Verified,
					Phase:  enums.Elimination,
				},
			},
			wantErr: false,
			beforeTest: func(args args, mockTeamRepo *mockTeamRepo.MockTeamRepository) {
				mockTeamRepo.EXPECT().
					UpdateTeam(gomock.Any(), args.id, args.team).
					Return(nil)
			},
		},
		{
			name: "When updating a non-existing team status, it should return error item not found",
			args: args{
				context.TODO(),
				"77FA",
				&dto.TeamUpdate{
					Status: enums.Verified,
					Phase:  enums.Elimination,
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(args args, mockTeamRepo *mockTeamRepo.MockTeamRepository) {
				mockTeamRepo.EXPECT().
					UpdateTeam(gomock.Any(), args.id, args.team).
					Return(domain.ErrNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)

			defer mockCtr.Finish()

			repoMock := mockTeamRepo.NewMockTeamRepository(mockCtr)
			awsMock := mockAws.NewMockCloudStorage(mockCtr)

			if test.beforeTest != nil {
				test.beforeTest(test.args, repoMock)
			}

			svc := service.NewTeamService(repoMock, nil, nil, time.Second, awsMock)

			err := svc.UpdateTeamStatus(test.args.ctx, test.args.id, test.args.team)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error message")
			} else {
				assert.Nil(t, err, "Error should not be expected")

			}
		})
	}
}

func TestJoinTeam(t *testing.T) {
	type args struct {
		ctx       context.Context
		joinToken string
		userID    uuid.UUID
	}

	dummyUUID := uuid.New()
	teamXComp := entity.Team{
		ID:        "TeamXComp",
		JoinToken: "joinTokenTeamWithXComp",
		Competition: entity.Competition{
			ID:   1,
			Name: "XComp",
		},
	}
	teamNonXComp := entity.Team{
		ID:        "TeamNonXComp",
		JoinToken: "joinTokenTeamNonXComp",
		Competition: entity.Competition{
			ID:   2,
			Name: "NonXComp",
		},
	}
	teamUserLeader := entity.Team{
		ID:        "TeamUserLeader",
		LeaderID:  dummyUUID,
		JoinToken: "joinTokenTeamWhereUserLeader",
	}
	teamFull := entity.Team{
		ID:        "TeamFull",
		JoinToken: "joinTokenTeamFull",
		Members: []entity.DetailTeams{
			{UserID: uuid.New()},
			{UserID: uuid.New()},
		},
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		beforeTest func(
			args args,
			mockTeamRepo *mockTeamRepo.MockTeamRepository,
			mockUserRepo *mockUserRepo.MockUserRepository,
		)
	}{
		{
			name: "When joining first team, it should not return error",
			args: args{
				context.TODO(),
				teamXComp.JoinToken,
				dummyUUID,
			},
			wantErr: false,
			beforeTest: func(args args, mockTeamRepo *mockTeamRepo.MockTeamRepository, mockUserRepo *mockUserRepo.MockUserRepository) {
				mockUserRepo.EXPECT().FindUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).SetArg(0, entity.User{
					ID: dummyUUID,
					Teams: []entity.DetailTeams{
						{
							TeamID: "hehe",
							Team: entity.Team{
								CompetitionID: 24,
							},
						},
					},
				})
				mockTeamRepo.EXPECT().FetchOneByParams(gomock.Any(), &dto.TeamParams{JoinToken: teamXComp.JoinToken}, "Members").Return(teamXComp, nil)
				mockTeamRepo.EXPECT().FetchOneByParams(gomock.Any(), gomock.Any()).Return(entity.Team{}, domain.ErrNotFound)
				mockTeamRepo.EXPECT().InsertTeamMember(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "When joining a team but user already registered to a team with the different competition, it should not return error",
			args: args{
				context.TODO(),
				teamNonXComp.JoinToken,
				dummyUUID,
			},
			wantErr: false,
			beforeTest: func(args args, mockTeamRepo *mockTeamRepo.MockTeamRepository, mockUserRepo *mockUserRepo.MockUserRepository) {
				mockUserRepo.EXPECT().FindUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockTeamRepo.EXPECT().FetchOneByParams(gomock.Any(), &dto.TeamParams{JoinToken: teamNonXComp.JoinToken}, "Members").Return(teamXComp, nil)
				mockTeamRepo.EXPECT().FetchOneByParams(gomock.Any(), gomock.Any()).Return(entity.Team{}, domain.ErrNotFound)
				mockTeamRepo.EXPECT().InsertTeamMember(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "When joining a non-existing team, it should return error not found",
			args: args{
				context.TODO(),
				"nonExistingJoinToken",
				dummyUUID,
			},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(args args, mockTeamRepo *mockTeamRepo.MockTeamRepository, mockUserRepo *mockUserRepo.MockUserRepository) {
				mockUserRepo.EXPECT().FindUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockTeamRepo.EXPECT().FetchOneByParams(gomock.Any(), &dto.TeamParams{JoinToken: "nonExistingJoinToken"}, "Members").Return(entity.Team{}, domain.ErrNotFound)
			},
		},
		{
			name: "When joining a team but user already registered to a team with the same competition, it should return error user already registered",
			args: args{
				context.TODO(),
				teamXComp.JoinToken,
				dummyUUID,
			},
			wantErr:    true,
			wantErrMsg: domain.ErrUserAlreadyRegistered.Error(),
			beforeTest: func(args args, mockTeamRepo *mockTeamRepo.MockTeamRepository, mockUserRepo *mockUserRepo.MockUserRepository) {
				mockUserRepo.EXPECT().FindUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockTeamRepo.EXPECT().FetchOneByParams(gomock.Any(), &dto.TeamParams{JoinToken: teamXComp.JoinToken}, "Members").Return(teamXComp, domain.ErrUserAlreadyRegistered)
			},
		},
		{
			name: "When joining a team but user already registered to a team as leader, it should return error user already registered",
			args: args{
				context.TODO(),
				teamUserLeader.JoinToken,
				dummyUUID,
			},
			wantErr:    true,
			wantErrMsg: domain.ErrUserAlreadyRegistered.Error(),
			beforeTest: func(args args, mockTeamRepo *mockTeamRepo.MockTeamRepository, mockUserRepo *mockUserRepo.MockUserRepository) {
				mockUserRepo.EXPECT().FindUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockTeamRepo.EXPECT().FetchOneByParams(gomock.Any(), &dto.TeamParams{JoinToken: teamUserLeader.JoinToken}, "Members").Return(teamUserLeader, domain.ErrUserAlreadyRegistered)
			},
		},
		{
			name: "When joining a team but team is full, it should return error team full",
			args: args{
				context.TODO(),
				teamFull.JoinToken,
				dummyUUID,
			},
			wantErr:    true,
			wantErrMsg: domain.ErrTeamFull.Error(),
			beforeTest: func(args args, mockTeamRepo *mockTeamRepo.MockTeamRepository, mockUserRepo *mockUserRepo.MockUserRepository) {
				mockUserRepo.EXPECT().FindUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockTeamRepo.EXPECT().FetchOneByParams(gomock.Any(), &dto.TeamParams{JoinToken: teamFull.JoinToken}, "Members").Return(teamFull, nil)
				mockTeamRepo.EXPECT().FetchOneByParams(gomock.Any(), gomock.Any()).Return(entity.Team{}, domain.ErrNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)

			defer mockCtr.Finish()

			repoMock := mockTeamRepo.NewMockTeamRepository(mockCtr)

			userRepoMock := mockUserRepo.NewMockUserRepository(mockCtr)

			if test.beforeTest != nil {
				test.beforeTest(test.args, repoMock, userRepoMock)
			}

			svc := service.NewTeamService(repoMock, nil, userRepoMock, time.Second, nil)

			err := svc.JoinTeam(test.args.ctx, test.args.joinToken, test.args.userID)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error message")
			} else {
				assert.Nil(t, err, "Error should not be expected")

			}
		})
	}
}

func TestUpdateLeader(t *testing.T) {
	type args struct {
		ctx      context.Context
		teamId   string
		leaderId uuid.UUID
	}

	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
		beforeTest  func(args args, repoMock *mockTeamRepo.MockTeamRepository)
	}{
		{
			name: "When update a leader, it should return no error",
			args: args{
				ctx:      context.TODO(),
				teamId:   "88XX",
				leaderId: uuid.MustParse("5fca6a81-d366-4791-babd-dcc9bfd0805c"),
			},
			wantErr: false,
			beforeTest: func(args args, repoMock *mockTeamRepo.MockTeamRepository) {
				repoMock.EXPECT().
					FetchOneByParams(gomock.Any(), &dto.TeamParams{ID: args.teamId}).
					Return(entity.Team{
						ID:            args.teamId,
						LeaderID:      uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
						CompetitionID: 1,
					}, nil)

				repoMock.EXPECT().
					FetchAllTeams(gomock.Any(), &dto.TeamParams{CompetitionID: 1}, "Members").
					Return([]entity.Team{
						{
							ID:            args.teamId,
							LeaderID:      uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
							CompetitionID: 1,
							Members: []entity.DetailTeams{
								{
									UserID: uuid.MustParse("83e2f35a-7b81-4243-a533-754ebf1727f0"),
									TeamID: args.teamId,
								},
								{
									UserID: uuid.MustParse("bb2ce198-eda8-4d8e-bfd5-817ecfd3f0e6"),
									TeamID: args.teamId,
								},
							},
						},
						{
							ID:            "77ff",
							LeaderID:      uuid.MustParse("e814a712-cc83-4ca0-8b5c-ceeb7e900405"),
							CompetitionID: 2,
							Members: []entity.DetailTeams{
								{
									UserID: uuid.MustParse("83e2f35a-7b81-4243-a533-754ebf1727f0"),
									TeamID: "77ff",
								},
								{
									UserID: uuid.MustParse("bb2ce198-eda8-4d8e-bfd5-817ecfd3f0e6"),
									TeamID: "77ff",
								},
							},
						},
					}, nil)

				repoMock.EXPECT().Begin()

				repoMock.EXPECT().
					UpdateTeam(gomock.Any(), args.teamId, &dto.TeamUpdate{LeaderID: args.leaderId}).
					Return(nil).
					AnyTimes()

				repoMock.EXPECT().
					DeleteTeamMember(gomock.Any(), &entity.DetailTeams{UserID: args.leaderId, TeamID: args.teamId}).
					Return(nil).
					AnyTimes()

				repoMock.EXPECT().Commit()
			},
		},
		{
			name: "When update a leader with a user that already join the same competition, it should return error already registered",
			args: args{
				ctx:      context.TODO(),
				teamId:   "88XX",
				leaderId: uuid.MustParse("5fca6a81-d366-4791-babd-dcc9bfd0805c"),
			},
			wantErr:     true,
			expectedErr: domain.ErrUserAlreadyRegistered,
			beforeTest: func(args args, repoMock *mockTeamRepo.MockTeamRepository) {
				repoMock.EXPECT().
					FetchOneByParams(gomock.Any(), &dto.TeamParams{ID: args.teamId}).
					Return(entity.Team{
						ID:            args.teamId,
						LeaderID:      uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
						CompetitionID: 1,
					}, nil)

				repoMock.EXPECT().
					FetchAllTeams(gomock.Any(), &dto.TeamParams{CompetitionID: 1}, "Members").
					Return([]entity.Team{
						{
							ID:            args.teamId,
							LeaderID:      uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
							CompetitionID: 1,
							Members: []entity.DetailTeams{
								{
									UserID: uuid.MustParse("83e2f35a-7b81-4243-a533-754ebf1727f0"),
									TeamID: args.teamId,
								},
								{
									UserID: uuid.MustParse("bb2ce198-eda8-4d8e-bfd5-817ecfd3f0e6"),
									TeamID: args.teamId,
								},
							},
						},
						{
							ID:            "77ff",
							LeaderID:      uuid.MustParse("e814a712-cc83-4ca0-8b5c-ceeb7e900405"),
							CompetitionID: 1,
							Members: []entity.DetailTeams{
								{
									UserID: uuid.MustParse("5fca6a81-d366-4791-babd-dcc9bfd0805c"),
									TeamID: "77ff",
								},
								{
									UserID: uuid.MustParse("bb2ce198-eda8-4d8e-bfd5-817ecfd3f0e6"),
									TeamID: "77ff",
								},
							},
						},
					}, nil)
			},
		},
		{
			name: "When update a leader with a user that is already a leader from a team, it should return error already registered",
			args: args{
				ctx:      context.TODO(),
				teamId:   "88XX",
				leaderId: uuid.MustParse("5fca6a81-d366-4791-babd-dcc9bfd0805c"),
			},
			wantErr:     true,
			expectedErr: domain.ErrDuplicateEntry,
			beforeTest: func(args args, repoMock *mockTeamRepo.MockTeamRepository) {
				repoMock.EXPECT().
					FetchOneByParams(gomock.Any(), &dto.TeamParams{ID: args.teamId}).
					Return(entity.Team{
						ID:            args.teamId,
						LeaderID:      uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
						CompetitionID: 1,
					}, nil)

				repoMock.EXPECT().
					FetchAllTeams(gomock.Any(), &dto.TeamParams{CompetitionID: 1}, "Members").
					Return([]entity.Team{
						{
							ID:            args.teamId,
							LeaderID:      uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
							CompetitionID: 1,
							Members: []entity.DetailTeams{
								{
									UserID: uuid.MustParse("83e2f35a-7b81-4243-a533-754ebf1727f0"),
									TeamID: args.teamId,
								},
								{
									UserID: uuid.MustParse("bb2ce198-eda8-4d8e-bfd5-817ecfd3f0e6"),
									TeamID: args.teamId,
								},
							},
						},
						{
							ID:            "77ff",
							LeaderID:      uuid.MustParse("e814a712-cc83-4ca0-8b5c-ceeb7e900405"),
							CompetitionID: 2,
							Members: []entity.DetailTeams{
								{
									UserID: uuid.MustParse("5fca6a81-d366-4791-babd-dcc9bfd0807d"),
									TeamID: "77ff",
								},
								{
									UserID: uuid.MustParse("bb2ce198-eda8-4d8e-bfd5-817ecfd3f0e6"),
									TeamID: "77ff",
								},
							},
						},
					}, nil)

				repoMock.EXPECT().Begin()

				repoMock.EXPECT().
					UpdateTeam(gomock.Any(), args.teamId, &dto.TeamUpdate{LeaderID: args.leaderId}).
					Return(domain.ErrDuplicateEntry).
					AnyTimes()

				repoMock.EXPECT().
					DeleteTeamMember(gomock.Any(), &entity.DetailTeams{UserID: args.leaderId, TeamID: args.teamId}).
					Return(nil).
					AnyTimes()

				repoMock.EXPECT().Rollback()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)

			defer mockCtr.Finish()

			repoMock := mockTeamRepo.NewMockTeamRepository(mockCtr)

			if test.beforeTest != nil {
				test.beforeTest(test.args, repoMock)
			}

			svc := service.NewTeamService(repoMock, nil, nil, 1*time.Second, nil)

			err := svc.UpdateLeader(test.args.ctx, test.args.teamId, test.args.leaderId)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.expectedErr, err, "Expecting same error")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}

		})
	}
}

func TestRemoveMember(t *testing.T) {
	type args struct {
		ctx    context.Context
		member *entity.DetailTeams
	}

	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
		beforeTest  func(args args, repoMock *mockTeamRepo.MockTeamRepository)
	}{
		{
			name:    "When deleting a member from a team, it should return no error",
			wantErr: false,
			args: args{
				context.TODO(),
				&entity.DetailTeams{
					TeamID: "77FF",
					UserID: uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
				},
			},
			beforeTest: func(args args, repoMock *mockTeamRepo.MockTeamRepository) {
				repoMock.EXPECT().DeleteTeamMember(gomock.Any(), args.member).Return(nil)
			},
		},
		{
			name:        "When deleting non existing member from a team, it should return error not found",
			wantErr:     true,
			expectedErr: domain.ErrNotFound,
			args: args{
				context.TODO(),
				&entity.DetailTeams{
					TeamID: "77FF",
					UserID: uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
				},
			},
			beforeTest: func(args args, repoMock *mockTeamRepo.MockTeamRepository) {
				repoMock.EXPECT().DeleteTeamMember(gomock.Any(), args.member).Return(domain.ErrNotFound)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)

			defer mockCtr.Finish()

			repoMock := mockTeamRepo.NewMockTeamRepository(mockCtr)

			if test.beforeTest != nil {
				test.beforeTest(test.args, repoMock)
			}

			svc := service.NewTeamService(repoMock, nil, nil, 1*time.Second, nil)

			err := svc.RemoveMember(test.args.ctx, test.args.member)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.expectedErr, err, "Expecting same error")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}
		})
	}
}

func TestCountTeamNUniv(t *testing.T) {
	tests := []struct {
		name        string
		beforeTest  func(repoMock *mockTeamRepo.MockTeamRepository)
		wantErr     bool
		want        dto.TeamNUnivCounter
		expectedErr error
	}{
		{
			name:    "When counting team and university, it should return no error",
			wantErr: false,
			want: dto.TeamNUnivCounter{
				TeamCounter: 10,
				UnivCounter: 5,
			},
			beforeTest: func(repoMock *mockTeamRepo.MockTeamRepository) {
				repoMock.EXPECT().Count(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(10), nil)
				repoMock.EXPECT().Count(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(5), nil)
			},
			expectedErr: nil,
		},
		{
			name:    "When counting team and university but there's an error when count team, it should return error",
			wantErr: true,
			want:    dto.TeamNUnivCounter{},
			beforeTest: func(repoMock *mockTeamRepo.MockTeamRepository) {
				repoMock.EXPECT().Count(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(10), domain.ErrInternalServer)
			},
			expectedErr: domain.ErrInternalServer,
		},
		{
			name:    "When counting team and university but there's an error when count univ, it should return error",
			wantErr: true,
			want:    dto.TeamNUnivCounter{},
			beforeTest: func(repoMock *mockTeamRepo.MockTeamRepository) {
				repoMock.EXPECT().Count(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(10), nil)
				repoMock.EXPECT().Count(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(5), domain.ErrInternalServer)
			},
			expectedErr: domain.ErrInternalServer,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)

			defer mockCtr.Finish()

			repoMock := mockTeamRepo.NewMockTeamRepository(mockCtr)

			if test.beforeTest != nil {
				test.beforeTest(repoMock)
			}

			svc := service.NewTeamService(repoMock, nil, nil, 1*time.Second, nil)

			result, err := svc.CountTeamNUniv(context.Background())

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.expectedErr, err, "Expecting same error")
			} else {
				assert.Nil(t, err, "Error should not be expected")
				assert.Equal(t, test.want, result, "Expecting same result")
			}
		})
	}
}

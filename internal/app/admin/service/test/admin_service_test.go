package test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	mockAdminRepo "github.com/hology8/hology-be/internal/app/admin/repository/mock"
	"github.com/hology8/hology-be/internal/app/admin/service"
	mockCompeRepo "github.com/hology8/hology-be/internal/app/competition/repository/mock"
	mockTeamRepo "github.com/hology8/hology-be/internal/app/team/repository/mock"
	mockUserRepo "github.com/hology8/hology-be/internal/app/user/repository/mock"
	mockBcrypt "github.com/hology8/hology-be/pkg/bcrypt/mock"
	mockGomail "github.com/hology8/hology-be/pkg/gomail/mock"
	mockJwt "github.com/hology8/hology-be/pkg/jwt/mock"
	"github.com/stretchr/testify/assert"
)

type mockObjects struct {
	adminRepo       *mockAdminRepo.MockAdminRepository
	userRepo        *mockUserRepo.MockUserRepository
	teamRepo        *mockTeamRepo.MockTeamRepository
	competitionRepo *mockCompeRepo.MockCompetitionRepository
	jwt             *mockJwt.MockJwtInterface
	bcrypt          *mockBcrypt.MockBcryptInterface
	gomail          *mockGomail.MockGoMailInterface
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mockAdminRepo.NewMockAdminRepository(ctrl)
	jwt := mockJwt.NewMockJwtInterface(ctrl)
	bcrypt := mockBcrypt.NewMockBcryptInterface(ctrl)
	gomail := mockGomail.NewMockGoMailInterface(ctrl)

	mockObj := mockObjects{
		adminRepo: r,
		jwt:       jwt,
		bcrypt:    bcrypt,
		gomail:    gomail,
	}

	type args struct {
		adminLogin    dto.AdminLogin
		admin         entity.Admin
		loginResponse dto.AdminLoginResponse
	}

	mockArgs := args{
		adminLogin: dto.AdminLogin{
			Username: "admin",
			Password: "randomstring",
		},
		admin: entity.Admin{
			Username: "admin",
			Password: "randomstring",
		},
		loginResponse: dto.AdminLoginResponse{
			Token: "token",
		},
	}

	test := []struct {
		name        string
		args        args
		beforeTest  func(mockObject mockObjects)
		want        dto.AdminLoginResponse
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when login success, it should return token",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.adminRepo.EXPECT().FindAdmin(&entity.Admin{}, &dto.AdminParam{Username: mockArgs.adminLogin.Username}).SetArg(0, mockArgs.admin).Return(nil)
				mockObject.bcrypt.EXPECT().Compare(mockArgs.adminLogin.Password, mockArgs.admin.Password).Return(true)
				mockObject.jwt.EXPECT().GenerateToken(mockArgs.admin.ID, "", mockArgs.admin.RoleID).Return("token", nil)
			},
			want:    mockArgs.loginResponse,
			wantErr: false,
		},
		{
			name: "when login but failed to find admin, it should return internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.adminRepo.EXPECT().FindAdmin(&entity.Admin{}, &dto.AdminParam{Username: mockArgs.adminLogin.Username}).SetArg(0, mockArgs.admin).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when login but failed failed to find admin(wrong user name), it should return username or password wrong",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.adminRepo.EXPECT().FindAdmin(&entity.Admin{}, &dto.AdminParam{Username: mockArgs.adminLogin.Username}).SetArg(0, entity.Admin{}).Return(domain.ErrNotFound)
			},
			wantErr:     true,
			expectedErr: domain.ErrWrongEmailOrPassword,
		},
		{
			name: "when login but failed to compare password, it should return username or password wrong",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.adminRepo.EXPECT().FindAdmin(&entity.Admin{}, &dto.AdminParam{Username: mockArgs.adminLogin.Username}).SetArg(0, mockArgs.admin).Return(nil)
				mockObject.bcrypt.EXPECT().Compare(mockArgs.adminLogin.Password, mockArgs.admin.Password).Return(false)
			},
			wantErr:     true,
			expectedErr: domain.ErrWrongEmailOrPassword,
		},
		{
			name: "when login but failed to generate token, it should return internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.adminRepo.EXPECT().FindAdmin(&entity.Admin{}, &dto.AdminParam{Username: mockArgs.adminLogin.Username}).SetArg(0, mockArgs.admin).Return(nil)
				mockObject.bcrypt.EXPECT().Compare(mockArgs.adminLogin.Password, mockArgs.admin.Password).Return(true)
				mockObject.jwt.EXPECT().GenerateToken(mockArgs.admin.ID, "", mockArgs.admin.RoleID).Return("token", domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when login but operation time exceeded, it should return timeout error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.adminRepo.EXPECT().FindAdmin(&entity.Admin{}, &dto.AdminParam{Username: mockArgs.adminLogin.Username}).SetArg(0, mockArgs.admin).Return(nil)
				mockObject.bcrypt.EXPECT().Compare(mockArgs.adminLogin.Password, mockArgs.admin.Password).Return(true)
				mockObject.jwt.EXPECT().GenerateToken(mockArgs.admin.ID, "", mockArgs.admin.RoleID).
					DoAndReturn(func(any1 any, any2 any, any3 any) (string, error) {
						time.Sleep(time.Millisecond * 1000)
						return "token", nil
					})
			},
			wantErr:     true,
			expectedErr: domain.ErrTimeout,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeTest(mockObj)

			s := service.NewAdminService(mockObj.adminRepo, nil, nil, nil, mockObj.bcrypt, mockObj.jwt, mockObj.gomail, time.Millisecond*500)

			got, err := s.Login(context.Background(), tt.args.adminLogin)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSendEmails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mockAdminRepo.NewMockAdminRepository(ctrl)
	userRepo := mockUserRepo.NewMockUserRepository(ctrl)
	teamRepo := mockTeamRepo.NewMockTeamRepository(ctrl)
	competitionRepo := mockCompeRepo.NewMockCompetitionRepository(ctrl)
	gomail := mockGomail.NewMockGoMailInterface(ctrl)

	mockObj := mockObjects{
		adminRepo:       r,
		userRepo:        userRepo,
		teamRepo:        teamRepo,
		competitionRepo: competitionRepo,
		gomail:          gomail,
	}

	type args struct {
		to           string
		emailMessage dto.EmailMessage
	}

	tests := []struct {
		name        string
		args        args
		beforeTest  func(mockObject mockObjects)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when send email to all users, it should return no error",
			args: args{
				to: "all",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), []interface{}{true}, nil, nil).
					Return([]entity.User{
						{
							Email: "email",
						},
					}, dto.PaginationResponse{}, nil)

				mockObject.gomail.EXPECT().SendEmails(gomock.Any(), gomock.Any(), []string{"email"}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "when send email but the param is empty, it should return error missing attribute",
			args: args{
				to: "",
			},
			wantErr:     true,
			expectedErr: domain.ErrMissingAttribute,
		},
		{
			name: "when send email to all users but failed to find users, it should return not found error",
			args: args{
				to: "all",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), []interface{}{true}, nil, nil).
					Return([]entity.User{}, dto.PaginationResponse{}, domain.ErrNotFound)
			},
			wantErr:     true,
			expectedErr: domain.ErrNotFound,
		},
		{
			name: "when send email to all users but failed to find users, it should return internal server error",
			args: args{
				to: "all",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), []interface{}{true}, nil, nil).
					Return([]entity.User{}, dto.PaginationResponse{}, domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when send email to all users but operation time exceeded, it should return timeout error",
			args: args{
				to: "all",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.All()).
					DoAndReturn(
						func(any1 any, any2 any, any3 any, any4 any, preload ...interface{}) ([]entity.User, dto.PaginationResponse, error) {
							time.Sleep(time.Millisecond * 1000)
							return nil, dto.PaginationResponse{}, nil
						},
					)
			},
			wantErr:     true,
			expectedErr: domain.ErrTimeout,
		},
		{
			name: "when send email to a specific competition, it should return no error",
			args: args{
				to: "competition",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
					Name:    []string{"kompetisi1"},
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.competitionRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.All()).
					Return([]entity.Competition{
						{
							Teams: []entity.Team{
								{
									Members: []entity.DetailTeams{
										{
											User: entity.User{
												Email:           "email",
												EmailIsVerified: true,
											},
										},
									},
								},
							},
						},
					}, nil)

				mockObject.userRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), nil, nil).
					Return([]entity.User{
						{
							Email:           "email2",
							EmailIsVerified: true,
						},
					}, dto.PaginationResponse{}, nil)
				mockObject.gomail.EXPECT().SendEmails(gomock.Any(), gomock.Any(), []string{"email", "email2"}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "when send email to a specific competition but the name struct is empty, it should return status bad request",
			args: args{
				to: "competition",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
				},
			},
			wantErr:     true,
			expectedErr: domain.ErrMissingAttribute,
		},
		{
			name: "when send email to a specific competition but failed to find any competition, it should return not found error",
			args: args{
				to: "competition",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
					Name:    []string{"kompetisi1"},
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.competitionRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]entity.Competition{}, domain.ErrNotFound)
			},
			wantErr:     true,
			expectedErr: domain.ErrNotFound,
		},
		{
			name: "when send email to a specific competition but failed to find any competition, it should return internal server error",
			args: args{
				to: "competition",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
					Name:    []string{"kompetisi1"},
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.competitionRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]entity.Competition{}, domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when send email to a specific competition but failed to find leader(user), it should return not found error",
			args: args{
				to: "competition",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
					Name:    []string{"kompetisi1"},
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.competitionRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]entity.Competition{
						{
							Teams: []entity.Team{
								{
									Members: []entity.DetailTeams{
										{
											User: entity.User{
												Email:           "email",
												EmailIsVerified: true,
											},
										},
									},
								},
							},
						},
					}, nil)

				mockObject.userRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), nil, nil).
					Return([]entity.User{}, dto.PaginationResponse{}, domain.ErrNotFound)
			},
			wantErr:     true,
			expectedErr: domain.ErrNotFound,
		},
		{
			name: "when send email to a specific competition but failed to find leader(user), it should return internal server error",
			args: args{
				to: "competition",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
					Name:    []string{"kompetisi1"},
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.competitionRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]entity.Competition{
						{
							Teams: []entity.Team{
								{
									Members: []entity.DetailTeams{
										{
											User: entity.User{
												Email:           "email",
												EmailIsVerified: true,
											},
										},
									},
								},
							},
						},
					}, nil)

				mockObject.userRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), nil, nil).
					Return([]entity.User{}, dto.PaginationResponse{}, domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when send email to a specific competition but operation time exceeded, it should return timeout error",
			args: args{
				to: "competition",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
					Name:    []string{"kompetisi1"},
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.competitionRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(any1 any, any2 any, any3 any, any4 ...any) ([]entity.Competition, error) {
						time.Sleep(time.Millisecond * 1000)
						return nil, nil
					})

				mockObject.userRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), nil, nil).
					DoAndReturn(func(any1 any, any2 any, any3 any, any4 any, preload ...any) ([]entity.User, dto.PaginationResponse, error) {
						time.Sleep(time.Millisecond * 1000)
						return nil, dto.PaginationResponse{}, nil
					})
			},
			wantErr:     true,
			expectedErr: domain.ErrTimeout,
		},
		{
			name: "when send email to a specific team, it should return no error",
			args: args{
				to: "team",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
					Name:    []string{"team1"},
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.teamRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]entity.Team{
						{
							LeaderID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						},
					}, dto.PaginationResponse{}, nil)

				mockObject.userRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), nil, nil).
					Return([]entity.User{
						{
							Email:           "email1",
							EmailIsVerified: true,
						},
					}, dto.PaginationResponse{}, nil)

				mockObject.gomail.EXPECT().SendEmails(gomock.Any(), gomock.Any(), []string{"email1"}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "when send email to a specific team but the name struct is empty, it should return status bad request",
			args: args{
				to: "team",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
				},
			},
			wantErr:     true,
			expectedErr: domain.ErrMissingAttribute,
		},
		{
			name: "when send email to a specific team but failed to find any team, it should return not found error",
			args: args{
				to: "team",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
					Name:    []string{"team1"},
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.teamRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]entity.Team{}, dto.PaginationResponse{}, domain.ErrNotFound)
			},
			wantErr:     true,
			expectedErr: domain.ErrNotFound,
		},
		{
			name: "when send email to a specific team but failed to find any team, it should return internal server error",
			args: args{
				to: "team",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
					Name:    []string{"team1"},
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.teamRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]entity.Team{}, dto.PaginationResponse{}, domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when send email to a specific team but failed to find leader(user), it should return not found error",
			args: args{
				to: "team",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
					Name:    []string{"team1"},
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.teamRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]entity.Team{
						{
							LeaderID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						},
					}, dto.PaginationResponse{}, nil)

				mockObject.userRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), nil, nil).
					Return([]entity.User{}, dto.PaginationResponse{}, domain.ErrNotFound)
			},
			wantErr:     true,
			expectedErr: domain.ErrNotFound,
		},
		{
			name: "when send email to a specific team but failed to find leader(user), it should return internal server error",
			args: args{
				to: "team",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
					Name:    []string{"team1"},
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.teamRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]entity.Team{
						{
							LeaderID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						},
					}, dto.PaginationResponse{}, nil)

				mockObject.userRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), nil, nil).
					Return([]entity.User{}, dto.PaginationResponse{}, domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when send email to a specific team but operation time exceeded, it should return timeout error",
			args: args{
				to: "team",
				emailMessage: dto.EmailMessage{
					Subject: "subject",
					Content: "body",
					Name:    []string{"team1"},
				},
			},
			beforeTest: func(mockObject mockObjects) {
				mockObject.teamRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(any1 any, any2 any, any3 any, any4 any, any5 any, preload ...any) ([]entity.Team, dto.PaginationResponse, error) {
						time.Sleep(time.Millisecond * 1000)
						return nil, dto.PaginationResponse{}, nil
					})

				mockObject.userRepo.EXPECT().FetchAllByConditionAndRelation(gomock.Any(), gomock.Any(), nil, nil).
					DoAndReturn(func(any1 any, any2 any, any3 any, any4 any, preload ...any) ([]entity.User, dto.PaginationResponse, error) {
						time.Sleep(time.Millisecond * 1000)
						return nil, dto.PaginationResponse{}, nil
					})
			},
			wantErr:     true,
			expectedErr: domain.ErrTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.beforeTest != nil {
				tt.beforeTest(mockObj)
			}

			s := service.NewAdminService(mockObj.adminRepo, mockObj.userRepo, mockObj.competitionRepo, mockObj.teamRepo, nil, nil, mockObj.gomail, time.Millisecond*500)

			err := s.SendEmail(context.Background(), tt.args.to, tt.args.emailMessage)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

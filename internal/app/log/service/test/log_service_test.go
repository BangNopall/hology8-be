package test

import (
	"context"
	"testing"
	"time"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
	logMocks "github.com/BangNopall/hology8-be/internal/app/log/repository/mock"
	"github.com/BangNopall/hology8-be/internal/app/log/service"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInsertLog(t *testing.T) {
	type args struct {
		ctx context.Context
		log *dto.LogRequest
	}

	dummyUUID := uuid.New().String()

	tests := []struct {
		name        string
		args        args
		wantErr     bool
		beforeTests func(mockServiceRepo *logMocks.MockLogRepository)
	}{
		{
			name: "When storing a log data, it should not return error",
			args: args{
				ctx: context.TODO(),
				log: &dto.LogRequest{
					Action:  "test action",
					AdminID: dummyUUID,
				},
			},
			wantErr: false,
			beforeTests: func(mockServiceRepo *logMocks.MockLogRepository) {
				mockServiceRepo.EXPECT().
					InsertLog(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockServiceRepo := logMocks.NewMockLogRepository(mockCtr)

			test.beforeTests(mockServiceRepo)

			logService := service.NewLogService(mockServiceRepo, time.Second*2)

			err := logService.InsertLog(test.args.ctx, test.args.log)

			if !test.wantErr {
				assert.Nil(t, err, "Error should not be expected")
			} else {
				assert.NotNil(t, err, "Expecting error to be thrown")
			}
		})
	}
}

func TestFetchAllLogs(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	mockTime := time.Now()
	dummyUUID := uuid.New()

	tests := []struct {
		name        string
		args        args
		wantErr     bool
		want        []dto.LogResponse
		expectedErr error
		beforeTests func(mockServiceRepo *logMocks.MockLogRepository)
	}{
		{
			name: "When fetching all logs, it should not return error",
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
			want: []dto.LogResponse{
				{
					ID:        1,
					Action:    "test action",
					Fullname:  "Test Admin",
					CreatedAt: mockTime,
				},
				{
					ID:        2,
					Action:    "test action2",
					Fullname:  "Test Admin2",
					CreatedAt: mockTime,
				},
				{
					ID:        3,
					Action:    "test action3",
					Fullname:  "Test Admin3",
					CreatedAt: mockTime,
				},
			},
			beforeTests: func(mockServiceRepo *logMocks.MockLogRepository) {
				mockServiceRepo.EXPECT().
					FetchAll(gomock.Any(), "Admin").
					Return([]entity.Log{
						{
							ID:      1,
							Action:  "test action",
							AdminID: dummyUUID,
							Admin: entity.Admin{
								ID:       dummyUUID,
								Fullname: "Test Admin",
								Username: "testadmin",
								Password: "testpassword",
							},
							CreatedAt: mockTime,
						},
						{
							ID:      2,
							Action:  "test action2",
							AdminID: dummyUUID,
							Admin: entity.Admin{
								ID:       dummyUUID,
								Fullname: "Test Admin2",
								Username: "testadmin2",
								Password: "testpassword2",
							},
							CreatedAt: mockTime,
						},
						{
							ID:      3,
							Action:  "test action3",
							AdminID: dummyUUID,
							Admin: entity.Admin{
								ID:       dummyUUID,
								Fullname: "Test Admin3",
								Username: "testadmin3",
								Password: "testpassword3",
							},
							CreatedAt: mockTime,
						},
					}, nil)
			},
		},
		{
			name: "When fetching all logs, it should return error",
			args: args{
				ctx: context.TODO(),
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
			beforeTests: func(mockServiceRepo *logMocks.MockLogRepository) {
				mockServiceRepo.EXPECT().
					FetchAll(gomock.Any(), "Admin").
					Return(nil, domain.ErrInternalServer)
			},
		},
		{
			name: "When fetching all logs but operation time exceeded, it should return time out",
			args: args{
				ctx: context.TODO(),
			},
			wantErr:     true,
			expectedErr: domain.ErrTimeout,
			beforeTests: func(mockServiceRepo *logMocks.MockLogRepository) {
				mockServiceRepo.EXPECT().
					FetchAll(gomock.Any(), "Admin").
					DoAndReturn(func(any any, any2 any) ([]entity.Log, error) {
						time.Sleep(time.Millisecond * 1000)
						return nil, nil
					})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockServiceRepo := logMocks.NewMockLogRepository(mockCtr)

			test.beforeTests(mockServiceRepo)

			logService := service.NewLogService(mockServiceRepo, time.Millisecond*500)

			logs, err := logService.FetchAllLogs(test.args.ctx)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.expectedErr, err, "Expecting error to be %v", test.expectedErr)
			} else {
				assert.Nil(t, err, "Error should not be expected")
				assert.Equal(t, test.want, logs, "Expecting logs to be equal, want %v, got %v", test.want, logs)
			}
		})
	}
}

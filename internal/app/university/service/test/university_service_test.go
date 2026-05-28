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
	mocks "github.com/BangNopall/hology8-be/internal/app/university/repository/mock"
	"github.com/BangNopall/hology8-be/internal/app/university/service"
)

func TestFetchAll(t *testing.T) {
	tests := []struct {
		name       string
		want       interface{}
		wantErr    bool
		beforeTest func(mockServiceRepo *mocks.MockUniversityRepository)
	}{
		{
			name: "When fetching all universities, it should not return error",
			want: []dto.UniversityResponse{
				{
					ID:   1,
					Name: "Universitas Brawijaya",
				},
				{
					ID:   2,
					Name: "Universitas Indonesia",
				},
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *mocks.MockUniversityRepository) {
				mockServiceRepo.EXPECT().
					FetchAll(gomock.Any(), nil).
					Return([]entity.University{
						{
							ID:   1,
							Name: "Universitas Brawijaya",
						},
						{
							ID:   2,
							Name: "Universitas Indonesia",
						},
					}, nil)
			},
		},
		{
			name:    "When fetching universities but operation time exceeded, it shoild return error timeout",
			want:    []dto.UniversityResponse(nil),
			wantErr: true,
			beforeTest: func(mockServiceRepo *mocks.MockUniversityRepository) {
				mockServiceRepo.EXPECT().
					FetchAll(gomock.Any(), nil).
					DoAndReturn(func(ctx context.Context, params *dto.UniversityParam) ([]entity.University, error) {
						time.Sleep(1 * time.Second)
						return []entity.University(nil), nil
					})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockServiceRepo := mocks.NewMockUniversityRepository(mockCtr)

			test.beforeTest(mockServiceRepo)

			uniSvc := service.NewUniversityService(mockServiceRepo, 500*time.Millisecond)

			universities, err := uniSvc.FetchAll(context.TODO(), nil)

			assert.Equal(t, test.want, universities, "Expeced universities result to be match")
			if test.wantErr {
				assert.NotNil(t, err, "Error expected to be thrown")
			} else {
				assert.Nil(t, err, "Expected error to be nil")
			}
		})
	}
}

func TestFetchByID(t *testing.T) {
	type args struct {
		id  int
		ctx context.Context
	}

	tests := []struct {
		name       string
		args       args
		want       interface{}
		wantErr    bool
		wantErrMsg string
		beforeTest func(mockServiceRepo *mocks.MockUniversityRepository)
	}{
		{
			name: "When fetching province with existing id, it should not return error",
			args: args{
				id:  1,
				ctx: context.TODO(),
			},
			want: dto.UniversityResponse{
				ID:   1,
				Name: "University Brawijaya",
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *mocks.MockUniversityRepository) {
				mockServiceRepo.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).
					Return(entity.University{
						ID:   1,
						Name: "University Brawijaya",
					}, nil)
			},
		},
		{
			name: "When fetching province with non-existing id, it should return error not found",
			args: args{
				id:  2,
				ctx: context.TODO(),
			},
			want:       dto.UniversityResponse{},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(mockServiceRepo *mocks.MockUniversityRepository) {
				mockServiceRepo.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).
					Return(entity.University{}, domain.ErrNotFound)
			},
		},
		{
			name: "When fetching province but operation timed out, it should return error time out",
			args: args{
				id:  2,
				ctx: context.TODO(),
			},
			want:       dto.UniversityResponse{},
			wantErr:    true,
			wantErrMsg: domain.ErrTimeout.Error(),
			beforeTest: func(mockServiceRepo *mocks.MockUniversityRepository) {
				mockServiceRepo.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int) (entity.University, error) {
						time.Sleep(1 * time.Second)
						return entity.University{}, nil
					})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockServiceRepo := mocks.NewMockUniversityRepository(mockCtr)

			test.beforeTest(mockServiceRepo)

			uniSvc := service.NewUniversityService(mockServiceRepo, 500*time.Millisecond)

			university, err := uniSvc.FetchByID(context.TODO(), test.args.id)

			assert.Equal(t, test.want, university, "Expeced university result to be match")
			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting error message to be same")
			} else {
				assert.Nil(t, err, "Expecting error to be nil")
			}
		})
	}
}

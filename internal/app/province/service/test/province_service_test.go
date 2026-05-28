package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	mocks "github.com/hology8/hology-be/internal/app/province/repository/mock"
	"github.com/hology8/hology-be/internal/app/province/service"
)

func TestFetchAll(t *testing.T) {
	tests := []struct {
		name       string
		want       interface{}
		wantErr    bool
		beforeTest func(mockServiceRepo *mocks.MockProvinceRepository)
	}{
		{
			name: "When fetching all provinces, it should not return error",
			want: []dto.ProvinceResponse{
				{
					ID:   1,
					Name: "Jawa Timur",
				},
				{
					ID:   2,
					Name: "Jawa Barat",
				},
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *mocks.MockProvinceRepository) {
				mockServiceRepo.EXPECT().
					FetchAll(gomock.Any()).
					Return([]entity.Province{
						{
							ID:   1,
							Name: "Jawa Timur",
						},
						{
							ID:   2,
							Name: "Jawa Barat",
						},
					}, nil)
			},
		},
		{
			name:    "When fetching provinces but operation time exceeded, it should return error timeout",
			want:    []dto.ProvinceResponse(nil),
			wantErr: true,
			beforeTest: func(mockServiceRepo *mocks.MockProvinceRepository) {
				mockServiceRepo.EXPECT().
					FetchAll(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]entity.Province, error) {
						time.Sleep(1 * time.Second)
						return []entity.Province(nil), nil
					})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockServiceRepo := mocks.NewMockProvinceRepository(mockCtr)

			test.beforeTest(mockServiceRepo)

			provinceService := service.NewProvinceService(mockServiceRepo, 500*time.Millisecond)

			provinces, err := provinceService.FetchAll(context.TODO())

			assert.Equal(t, test.want, provinces, "Expeced provinces result to be match")
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
		beforeTest func(mockServiceRepo *mocks.MockProvinceRepository)
	}{
		{
			name: "When fetching province with existing id, it should not return error",
			args: args{
				id:  1,
				ctx: context.TODO(),
			},
			want: dto.ProvinceResponse{
				ID:   1,
				Name: "Jawa Timur",
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *mocks.MockProvinceRepository) {
				mockServiceRepo.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).
					Return(entity.Province{
						ID:   1,
						Name: "Jawa Timur",
					}, nil)
			},
		},
		{
			name: "When fetching province with non-existing id, it should return error not found",
			args: args{
				id:  2,
				ctx: context.TODO(),
			},
			want:       dto.ProvinceResponse{},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(mockServiceRepo *mocks.MockProvinceRepository) {
				mockServiceRepo.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).
					Return(entity.Province{}, domain.ErrNotFound)
			},
		},
		{
			name: "When fetching province but operation timed out, it should return error time out",
			args: args{
				id:  2,
				ctx: context.TODO(),
			},
			want:       dto.ProvinceResponse{},
			wantErr:    true,
			wantErrMsg: domain.ErrTimeout.Error(),
			beforeTest: func(mockServiceRepo *mocks.MockProvinceRepository) {
				mockServiceRepo.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int) (entity.Province, error) {
						time.Sleep(1 * time.Second)
						return entity.Province{}, nil
					})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockServiceRepo := mocks.NewMockProvinceRepository(mockCtr)

			test.beforeTest(mockServiceRepo)

			provSvc := service.NewProvinceService(mockServiceRepo, 500*time.Millisecond)

			province, err := provSvc.FetchByID(context.TODO(), test.args.id)

			assert.Equal(t, test.want, province, "Expeced province result to be match")
			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting error message to be same")
			} else {
				assert.Nil(t, err, "Expecting error to be nil")
			}
		})
	}
}

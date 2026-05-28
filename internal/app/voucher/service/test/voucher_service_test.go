package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
	teamMocks "github.com/BangNopall/hology8-be/internal/app/team/repository/mock"
	voucherMocks "github.com/BangNopall/hology8-be/internal/app/voucher/repository/mock"
	"github.com/BangNopall/hology8-be/internal/app/voucher/service"
	pgsqlMock "github.com/BangNopall/hology8-be/internal/infra/database/mock"
)

func TestFetchAll(t *testing.T) {
	tests := []struct {
		name       string
		want       interface{}
		wantErr    bool
		beforeTest func(mockServiceRepo *voucherMocks.MockVoucherRepository)
	}{
		{
			name: "When fetching all vouchers, it should not return error",
			want: []dto.VoucherResponse{
				{
					ID:     "blabla-123",
					TeamID: "123",
				},
				{
					ID:     "holoholo-1",
					TeamID: "234",
				},
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *voucherMocks.MockVoucherRepository) {
				mockServiceRepo.EXPECT().
					FetchAll(gomock.Any()).
					Return([]entity.Voucher{
						{
							ID:     "blabla-123",
							TeamID: "123",
						},
						{
							ID:     "holoholo-1",
							TeamID: "234",
						},
					}, nil)
			},
		},
		{
			name:    "When fetching vouchers but operation time exceeded, it should return error timeout",
			want:    []dto.VoucherResponse(nil),
			wantErr: true,
			beforeTest: func(mockServiceRepo *voucherMocks.MockVoucherRepository) {
				mockServiceRepo.EXPECT().
					FetchAll(gomock.Any()).
					DoAndReturn(func(ctx context.Context) ([]entity.Voucher, error) {
						time.Sleep(1 * time.Second)
						return []entity.Voucher(nil), nil
					})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, _ := pgsqlMock.NewMockDB(t)

			defer close()

			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockVoucherRepo := voucherMocks.NewMockVoucherRepository(mockCtr)

			mockTeamRepo := teamMocks.NewMockTeamRepository(mockCtr)

			test.beforeTest(mockVoucherRepo)

			voucherService := service.NewVoucherService(mockVoucherRepo, mockTeamRepo, 500*time.Millisecond, db)

			vouchers, err := voucherService.FetchAll(context.TODO())

			assert.Equal(t, test.want, vouchers, "Expected vouchers result to be match")
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
		id  string
		ctx context.Context
	}

	tests := []struct {
		name       string
		args       args
		want       interface{}
		wantErr    bool
		wantErrMsg string
		beforeTest func(mockServiceRepo *voucherMocks.MockVoucherRepository)
	}{
		{
			name: "When fetching voucher with existing id, it should not return error",
			args: args{
				id:  "blabla-123",
				ctx: context.TODO(),
			},
			want: dto.VoucherResponse{
				ID:     "blabla-123",
				TeamID: "123",
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *voucherMocks.MockVoucherRepository) {
				mockServiceRepo.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).
					Return(entity.Voucher{
						ID:     "blabla-123",
						TeamID: "123",
					}, nil)
			},
		},
		{
			name: "When fetching voucher with non-existing id, it should return error not found",
			args: args{
				id:  "blabla-123",
				ctx: context.TODO(),
			},
			want:       dto.VoucherResponse{},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(mockServiceRepo *voucherMocks.MockVoucherRepository) {
				mockServiceRepo.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).
					Return(entity.Voucher{}, domain.ErrNotFound)
			},
		},
		{
			name: "When fetching voucher but operation timed out, it should return error time out",
			args: args{
				id:  "blabla-123",
				ctx: context.TODO(),
			},
			want:       dto.VoucherResponse{},
			wantErr:    true,
			wantErrMsg: domain.ErrTimeout.Error(),
			beforeTest: func(mockServiceRepo *voucherMocks.MockVoucherRepository) {
				mockServiceRepo.EXPECT().
					FetchByID(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id string) (entity.Voucher, error) {
						time.Sleep(1 * time.Second)
						return entity.Voucher{}, nil
					})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, _ := pgsqlMock.NewMockDB(t)

			defer close()

			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockVoucherRepo := voucherMocks.NewMockVoucherRepository(mockCtr)

			mockTeamRepo := teamMocks.NewMockTeamRepository(mockCtr)

			test.beforeTest(mockVoucherRepo)

			voucherService := service.NewVoucherService(mockVoucherRepo, mockTeamRepo, 500*time.Millisecond, db)

			voucher, err := voucherService.FetchByID(context.TODO(), test.args.id)

			assert.Equal(t, test.want, voucher, "Expected voucher result to be match")
			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting error message to be same")
			} else {
				assert.Nil(t, err, "Expecting error to be nil")
			}
		})
	}
}

func TestInsertVoucher(t *testing.T) {
	type args struct {
		voucher *dto.VoucherRequest
		ctx     context.Context
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		beforeTest func(mockServiceRepo *voucherMocks.MockVoucherRepository)
	}{
		{
			name: "When inserting compe, it should not return error",
			args: args{
				voucher: &dto.VoucherRequest{
					ID: "blabla-123",
				},
				ctx: context.TODO(),
			},
			wantErr: false,
			beforeTest: func(mockServiceRepo *voucherMocks.MockVoucherRepository) {
				mockServiceRepo.EXPECT().
					InsertVoucher(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, _ := pgsqlMock.NewMockDB(t)

			defer close()

			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockVoucherRepo := voucherMocks.NewMockVoucherRepository(mockCtr)

			mockTeamRepo := teamMocks.NewMockTeamRepository(mockCtr)

			test.beforeTest(mockVoucherRepo)

			voucherService := service.NewVoucherService(mockVoucherRepo, mockTeamRepo, 500*time.Millisecond, db)

			err := voucherService.InsertVoucher(context.TODO(), test.args.voucher)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting error message to be same")
			} else {
				assert.Nil(t, err, "Expecting error to be nil")
			}
		})
	}
}

func TestRedeemVoucher(t *testing.T) {
	type args struct {
		voucherRedeem *dto.VoucherRedeem
		ctx           context.Context
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		beforeTest func(mockVoucherRepo *voucherMocks.MockVoucherRepository, mockTeamRepo *teamMocks.MockTeamRepository, mock sqlmock.Sqlmock)
	}{
		{
			name: "When redeeming voucher, it should not return error",
			args: args{
				voucherRedeem: &dto.VoucherRedeem{
					ID:     "blabla-123",
					TeamID: "H8",
				},
				ctx: context.TODO(),
			},
			wantErr: false,
			beforeTest: func(mockVoucherRepo *voucherMocks.MockVoucherRepository, mockTeamRepo *teamMocks.MockTeamRepository, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mockVoucherRepo.EXPECT().FetchByID(gomock.Any(), gomock.Any())
				mockTeamRepo.EXPECT().FetchOneByID(gomock.Any(), gomock.Any())
				mockVoucherRepo.EXPECT().UpdateVoucher(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockTeamRepo.EXPECT().LinkVoucher(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mock.ExpectCommit()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockVoucherRepo := voucherMocks.NewMockVoucherRepository(mockCtr)

			mockTeamRepo := teamMocks.NewMockTeamRepository(mockCtr)

			test.beforeTest(mockVoucherRepo, mockTeamRepo, mock)

			voucherService := service.NewVoucherService(mockVoucherRepo, mockTeamRepo, 500*time.Millisecond, db)

			err := voucherService.RedeemVoucher(context.TODO(), test.args.voucherRedeem)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting error message to be same")
			} else {
				assert.Nil(t, err, "Expecting error to be nil")
			}
		})
	}
}

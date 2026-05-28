package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/entity"
	"github.com/BangNopall/hology8-be/internal/app/voucher/repository"
	pgsqlMock "github.com/BangNopall/hology8-be/internal/infra/database/mock"
)

func TestFetchAll(t *testing.T) {
	db, close, mock := pgsqlMock.NewMockDB(t)

	defer close()

	repo := repository.NewVoucherRepository(db)

	rows := sqlmock.NewRows([]string{"id", "team_id"}).
		AddRow("21-GACOR4-1212", "bla-bla-21").
		AddRow("21-GACOR4-1232", "2")

	mock.ExpectQuery(`SELECT \* FROM "vouchers"`).WillReturnRows(rows)

	vouchers, err := repo.FetchAll(context.TODO())

	if err != nil {
		t.Fatal(err)
	}

	expected := []entity.Voucher{
		{ID: "21-GACOR4-1212", TeamID: "bla-bla-21"},
		{ID: "21-GACOR4-1232", TeamID: "2"},
	}

	assert.Equal(t, expected, vouchers, "Fetched vouchers should match the expected data")
}

func TestFetchByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}

	tests := []struct {
		name       string
		args       args
		want       interface{}
		wantErr    bool
		wantErrMsg string
		beforeTest func(args args, mock sqlmock.Sqlmock)
	}{
		{
			name: "When fetching existing voucher, it should return the data",
			args: args{
				context.TODO(),
				"bla-bla-21",
			},
			want: entity.Voucher{
				ID:     "bla-bla-21",
				TeamID: "123",
			},
			wantErr: false,
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{"id", "team_id"}).
					AddRow("bla-bla-21", "123")

				mock.ExpectQuery(`SELECT \* FROM "vouchers" WHERE "vouchers"."id" = \$1 ORDER BY "vouchers"."id" LIMIT \$2`).
					WithArgs("bla-bla-21", 1).
					WillReturnRows(rows)
			},
		},
		{
			name: "When non-fetching existing voucher, it should return error not found",
			args: args{
				context.TODO(),
				"bla-bla-21",
			},
			want:       entity.Voucher{},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{"id", "team_id"})

				mock.ExpectQuery(`SELECT \* FROM "vouchers" WHERE "vouchers"."id" = \$1 ORDER BY "vouchers"."id" LIMIT \$2`).
					WithArgs("bla-bla-21", 1).
					WillReturnRows(rows)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			repo := repository.NewVoucherRepository(db)

			test.beforeTest(test.args, mock)

			voucher, err := repo.FetchByID(test.args.ctx, test.args.id)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error message")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}

			assert.Equal(t, test.want, voucher)
		})
	}
}

func TestInsertVoucher(t *testing.T) {
	type args struct {
		ctx     context.Context
		voucher *entity.Voucher
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		beforeTest func(args args, mock sqlmock.Sqlmock)
	}{
		{
			name: "When inserting voucher, it should not return error",
			args: args{
				context.TODO(),
				&entity.Voucher{
					ID: "something-123",
				},
			},
			wantErr: false,
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "vouchers" (.+) VALUES (.+)`).
					WithArgs(args.voucher.ID).
					WillReturnRows(sqlmock.NewRows([]string{"team_id"})).
					WillReturnError(nil)

				mock.ExpectCommit()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			repo := repository.NewVoucherRepository(db)

			test.beforeTest(test.args, mock)

			err := repo.InsertVoucher(test.args.ctx, test.args.voucher)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error message")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}
		})
	}
}

func TestUpdateVoucher(t *testing.T) {
	type args struct {
		ctx     context.Context
		voucher *entity.Voucher
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
		beforeTest func(args args, mock sqlmock.Sqlmock)
	}{
		{
			name: "When updating voucher, it should not return error",
			args: args{
				context.TODO(),
				&entity.Voucher{
					ID:     "bla-bla-21",
					TeamID: "123",
				},
			},
			wantErr: false,
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE "vouchers" SET .+`).
					WithArgs(args.voucher.TeamID, args.voucher.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "When updating non-existing voucher, it should return error",
			args: args{
				context.TODO(),
				&entity.Voucher{
					ID:     "bla-bla-21",
					TeamID: "123",
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE "vouchers" SET .+`).
					WithArgs(args.voucher.TeamID, args.voucher.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectCommit()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			repo := repository.NewVoucherRepository(db)

			test.beforeTest(test.args, mock)

			err := repo.UpdateVoucher(test.args.ctx, test.args.voucher, db)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error message")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}
		})
	}
}

package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/internal/app/province/repository"
	pgsqlMock "github.com/hology8/hology-be/internal/infra/database/mock"
)

func TestFetchAll(t *testing.T) {
	db, close, mock := pgsqlMock.NewMockDB(t)

	defer close()

	repo := repository.NewProvinceRepository(db)

	rows := sqlmock.NewRows([]string{"id", "province_name"}).
		AddRow(1, "Jawa Timur").
		AddRow(2, "Jawa Barat")

	mock.ExpectQuery(`SELECT \* FROM "provinces"`).WillReturnRows(rows)

	provinces, err := repo.FetchAll(context.TODO())

	if err != nil {
		t.Fatal(err)
	}

	expected := []entity.Province{
		{ID: 1, Name: "Jawa Timur"},
		{ID: 2, Name: "Jawa Barat"},
	}

	assert.Equal(t, expected, provinces, "Fetched provinces should match the expected data")
}

func TestFetchByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	tests := []struct {
		name       string
		args       args
		want       interface{}
		wantErr    bool
		wantErrMsg string
		beforeTest func(mock sqlmock.Sqlmock)
	}{
		{
			name: "When fetching existing province, it should return the data",
			args: args{
				context.TODO(),
				1,
			},
			want: entity.Province{
				ID:   1,
				Name: "Jawa Timur",
			},
			wantErr: false,
			beforeTest: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{"id", "province_name"}).
					AddRow(1, "Jawa Timur")

				mock.ExpectQuery(`SELECT \* FROM "provinces" WHERE "provinces"."id" = \$1 ORDER BY "provinces"."id" LIMIT \$2`).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
		},
		{
			name: "When fetching non-existing province, it should return error not found",
			args: args{
				context.TODO(),
				1,
			},
			want:       entity.Province{},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{"id", "province_name"})

				mock.ExpectQuery(`SELECT \* FROM "provinces" WHERE "provinces"."id" = \$1 ORDER BY "provinces"."id" LIMIT \$2`).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			repo := repository.NewProvinceRepository(db)

			test.beforeTest(mock)

			province, err := repo.FetchByID(test.args.ctx, test.args.id)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error message")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}

			assert.Equal(t, test.want, province)
		})
	}
}

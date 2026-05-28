package repository_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/internal/app/university/repository"
	pgsqlMock "github.com/hology8/hology-be/internal/infra/database/mock"
	"github.com/stretchr/testify/assert"
)

func TestFetchAll(t *testing.T) {
	db, close, mock := pgsqlMock.NewMockDB(t)

	defer close()

	repo := repository.NewUniversityRepository(db)

	rows := sqlmock.NewRows([]string{"id", "university_name"}).
		AddRow(1, "Universitas Brawijaya").
		AddRow(2, "Institut Teknologi Sepuluh November")

	mock.ExpectQuery(`SELECT \* FROM "universities"`).WillReturnRows(rows)

	universities, err := repo.FetchAll(context.TODO(), nil)

	if err != nil {
		t.Fatal(err)
	}

	expcted := []entity.University{
		{ID: 1, Name: "Universitas Brawijaya"},
		{ID: 2, Name: "Institut Teknologi Sepuluh November"},
	}

	assert.Equal(t, expcted, universities, "Fetched universities should match the expected data")
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
			name: "When fetching existing university, it should return the data",
			args: args{
				context.TODO(),
				1,
			},
			want: entity.University{
				ID:   1,
				Name: "Universitas Brawijaya",
			},
			wantErr: false,
			beforeTest: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{"id", "university_name"}).
					AddRow(1, "Universitas Brawijaya")

				mock.ExpectQuery(`SELECT \* FROM "universities" WHERE "universities"."id" = \$1 ORDER BY "universities"."id" LIMIT \$2`).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
		},
		{
			name: "When non-fetching existing province, it should return error not found",
			args: args{
				context.TODO(),
				1,
			},
			want:       entity.University{},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{"id", "university_name"})

				mock.ExpectQuery(`SELECT \* FROM "universities" WHERE "universities"."id" = \$1 ORDER BY "universities"."id" LIMIT \$2`).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			repo := repository.NewUniversityRepository(db)

			test.beforeTest(mock)

			university, err := repo.FetchByID(test.args.ctx, test.args.id)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error message")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}

			assert.Equal(t, test.want, university)
		})
	}
}

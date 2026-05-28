package test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/domain/entity"
	"github.com/BangNopall/hology8-be/internal/app/log/repository"
	pgsqlMock "github.com/BangNopall/hology8-be/internal/infra/database/mock"
)

func TestInsertLog(t *testing.T) {
	type args struct {
		ctx     context.Context
		log     *entity.Log
		fetchId int
	}

	type expectedFetch struct {
		log entity.Log
		err error
	}

	dummyUUID := uuid.New()

	tests := []struct {
		name          string
		args          args
		beforeTests   func(sqlmock.Sqlmock)
		wantErr       bool
		fetch         func(repo contracts.LogRepository, id int) (entity.Log, error)
		expectedFetch expectedFetch
	}{
		{
			name: "When storing a log data, it should return no error",
			args: args{
				ctx: context.TODO(),
				log: &entity.Log{
					ID:      1,
					Action:  "test action",
					AdminID: dummyUUID,
				},
				fetchId: 1,
			},
			beforeTests: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "logs" (.+) RETURNING`).
					WithArgs(dummyUUID, "test action", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()

				mock.ExpectQuery(`SELECT \* FROM "logs" WHERE "logs"."id" = \$1 ORDER BY "logs"."id" LIMIT \$2`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "action", "admin_id"}).AddRow(1, "test action", dummyUUID))
			},
			wantErr: false,
			fetch: func(repo contracts.LogRepository, id int) (entity.Log, error) {
				log, err := repo.FetchOneByID(context.TODO(), id)

				return log, err
			},
			expectedFetch: expectedFetch{
				log: entity.Log{
					ID:      1,
					Action:  "test action",
					AdminID: dummyUUID,
				},
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			if test.beforeTests != nil {
				test.beforeTests(mock)
			}

			repo := repository.NewLogRepository(db)

			err := repo.InsertLog(test.args.ctx, test.args.log)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}

			log, err := test.fetch(repo, test.args.fetchId)

			assert.Equal(t, test.expectedFetch.log, log, "Expecting log to be equal")
			assert.Equal(t, test.expectedFetch.err, err, "Expecting error to be equal")
		})
	}
}

func TestFetchAllLogs(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type expectedFetch struct {
		logs []entity.Log
		err  error
	}

	dummyUUID := uuid.New()

	tests := []struct {
		name          string
		args          args
		beforeTests   func(sqlmock.Sqlmock)
		wantErr       bool
		fetch         func(repo contracts.LogRepository) ([]entity.Log, error)
		expectedFetch expectedFetch
		expectedErr   error
	}{
		{
			name: "When fetching all logs, it should return no error",
			args: args{
				ctx: context.TODO(),
			},
			beforeTests: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "logs"`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "action", "admin_id"}).AddRow(1, "test action", dummyUUID))
			},
			wantErr: false,
			fetch: func(repo contracts.LogRepository) ([]entity.Log, error) {
				logs, err := repo.FetchAll(context.TODO())

				return logs, err
			},
			expectedFetch: expectedFetch{
				logs: []entity.Log{
					{
						ID:      1,
						Action:  "test action",
						AdminID: dummyUUID,
					},
				},
				err: nil,
			},
		},
		{
			name: "When fetching all logs but it get error, it should return internal server error",
			args: args{
				ctx: context.TODO(),
			},
			beforeTests: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "logs"`).
					WillReturnError(domain.ErrInternalServer)
			},
			fetch: func(repo contracts.LogRepository) ([]entity.Log, error) {
				logs, err := repo.FetchAll(context.TODO())

				return logs, err
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "When fetching all logs but it get error, it should return record not found",
			args: args{
				ctx: context.TODO(),
			},
			beforeTests: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "logs"`).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			fetch: func(repo contracts.LogRepository) ([]entity.Log, error) {
				logs, err := repo.FetchAll(context.TODO())

				return logs, err
			},
			wantErr:     true,
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			if test.beforeTests != nil {
				test.beforeTests(mock)
			}

			repo := repository.NewLogRepository(db)

			logs, err := test.fetch(repo)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.expectedErr, err, "Expecting error to be equal")
			} else {
				assert.Nil(t, err, "Error should not be expected")
				assert.Equal(t, test.expectedFetch.logs, logs, "Expecting logs to be equal")
			}
		})
	}
}

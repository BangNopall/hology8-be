package test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/internal/app/announcement/repository"
	pgsqlMock "github.com/hology8/hology-be/internal/infra/database/mock"
	_ "github.com/hology8/hology-be/pkg/log"
	"github.com/hology8/hology-be/pkg/uuid"
)

func TestFetchAnnouncementByTo(t *testing.T) {
	type args struct {
		ctx           context.Context
		teamID        string
		competitionID int
	}

	dummyCompetitionID := 1
	dummyTeamID := "1"

	tests := []struct {
		name        string
		args        args
		beforeTests func(sqlmock.Sqlmock)
		want        []entity.Announcement
	}{
		{
			name: "When fetching existing announcement data to all , it should return the announcement data to all",
			args: args{context.TODO(), "", 0},
			beforeTests: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{"id", "description"}).
					AddRow(1, "To All 1").
					AddRow(2, "To All 2")

				mock.NewRows([]string{"id", "description", "team_id", "competition_id"}).
					AddRow(3, "To Team 1", "1", nil).
					AddRow(4, "To Team 2", "2", nil).
					AddRow(5, "To Competition 1", nil, 1).
					AddRow(6, "To Competition 2", nil, 2)

				mock.ExpectQuery(`SELECT \* FROM "announcements"`).
					WillReturnRows(rows)
			},
			want: []entity.Announcement{{ID: 1, Description: "To All 1"}, {ID: 2, Description: "To All 2"}},
		},
		{
			name: "When fetching existing announcement data to competition, it should return the announcement data to competition",
			args: args{context.TODO(), "", 1},
			beforeTests: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{"id", "description", "competition_id"}).
					AddRow(5, "To Competition 1", 1)

				mock.NewRows([]string{"id", "description", "team_id", "competition_id"}).
					AddRow(1, "To All 1", nil, nil).
					AddRow(2, "To All 2", nil, nil).
					AddRow(3, "To Team 1", "1", nil).
					AddRow(4, "To Team 2", "2", nil).
					AddRow(6, "To Competition 2", 2, nil)

				mock.ExpectQuery(`SELECT \* FROM "announcements" WHERE team_id IS NULL AND competition_id = \$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: []entity.Announcement{{ID: 5, Description: "To Competition 1", CompetitionID: &dummyCompetitionID}},
		}, {
			name: "When fetching existing announcement data to team, it should return the announcement data to team",
			args: args{context.TODO(), "1", 0},
			beforeTests: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{"id", "description", "team_id"}).
					AddRow(3, "To Team 1", "1")

				mock.NewRows([]string{"id", "description", "team_id", "competition_id"}).
					AddRow(1, "To All 1", nil, nil).
					AddRow(2, "To All 2", nil, nil).
					AddRow(4, "To Team 2", "2", nil).
					AddRow(5, "To Competition 1", nil, 1).
					AddRow(6, "To Competition 2", nil, 2)

				mock.ExpectQuery(`SELECT \* FROM "announcements" WHERE team_id = \$1 AND competition_id IS NULL`).
					WithArgs("1").
					WillReturnRows(rows)
			},
			want: []entity.Announcement{{ID: 3, Description: "To Team 1", TeamID: &dummyTeamID}},
		},
		{
			name: "When fetching non-existing announcement data to team , it should return the empty announcement data",
			args: args{context.TODO(), "2", 0},
			beforeTests: func(mock sqlmock.Sqlmock) {
				mock.NewRows([]string{"id", "description", "team_id", "competition_id"}).
					AddRow(1, "To All 1", nil, nil).
					AddRow(2, "To All 2", nil, nil).
					AddRow(3, "To Team 1", "1", nil).
					AddRow(4, "To Team 2", "2", nil).
					AddRow(5, "To Competition 1", nil, 1).
					AddRow(6, "To Competition 2", nil, 2)

				emptyRows := mock.NewRows([]string{"id", "description", "team_id"})

				mock.ExpectQuery(`SELECT \* FROM "announcements" WHERE team_id = \$1 AND competition_id IS NULL`).
					WithArgs("2").
					WillReturnRows(emptyRows)
			},
			want: []entity.Announcement{},
		},
		{
			name: "When fetching non-existing announcement data to competition , it should return the empty announcement data",
			args: args{context.TODO(), "", 2},
			beforeTests: func(mock sqlmock.Sqlmock) {
				mock.NewRows([]string{"id", "description", "team_id", "competition_id"}).
					AddRow(1, "To All 1", nil, nil).
					AddRow(2, "To All 2", nil, nil).
					AddRow(3, "To Team 1", "1", nil).
					AddRow(4, "To Team 2", "2", nil).
					AddRow(5, "To Competition 1", nil, 1).
					AddRow(6, "To Competition 2", nil, 2)

				emptyRows := mock.NewRows([]string{"id", "description", "competition_id"})

				mock.ExpectQuery(`SELECT \* FROM "announcements" WHERE team_id IS NULL AND competition_id = \$1`).
					WithArgs(2).
					WillReturnRows(emptyRows)
			},
			want: []entity.Announcement{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			repo := repository.NewAnnouncementRepository(db)

			if test.beforeTests != nil {
				test.beforeTests(mock)
			}

			announcement, err := repo.FetchAnnouncementByTo(test.args.ctx, test.args.teamID, test.args.competitionID)

			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.want, announcement, "Expecting announcement data to be equal")
		})
	}
}

func TestInsertAnnouncement(t *testing.T) {
	type args struct {
		ctx          context.Context
		announcement *entity.Announcement
	}

	type expectedFetch struct {
		announcement entity.Announcement
		err          error
	}

	dummyUUID := uuid.New()

	tests := []struct {
		name          string
		args          args
		beforeTests   func(sqlmock.Sqlmock)
		wantErr       bool
		fetch         func(repo contracts.AnnouncementRepository) (entity.Announcement, error)
		expectedFetch expectedFetch
	}{
		{
			name: "When inserting announcement to all data, it should not return the error",
			args: args{context.TODO(), &entity.Announcement{Description: "Test Announcement", AdminID: dummyUUID}},
			beforeTests: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "announcements" (.+) RETURNING`).
					WithArgs(dummyUUID, "Test Announcement").
					WillReturnRows(sqlmock.NewRows([]string{"id", "description", "admin_id"}).AddRow(1, "Test Announcement", dummyUUID))

				mock.ExpectCommit()

				mock.ExpectQuery(`SELECT \* FROM "announcements" WHERE team_id IS NULL AND competition_id IS NULL`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "description", "admin_id"}).AddRow(1, "Test Announcement", dummyUUID))
			},
			wantErr: false,
			fetch: func(repo contracts.AnnouncementRepository) (entity.Announcement, error) {
				announcement, err := repo.FetchAnnouncementByTo(context.TODO(), "", 0)

				return announcement[0], err
			},
			expectedFetch: expectedFetch{
				announcement: entity.Announcement{ID: 1, Description: "Test Announcement", AdminID: dummyUUID},
				err:          nil,
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

			repo := repository.NewAnnouncementRepository(db)

			err := repo.InsertAnnouncement(test.args.ctx, test.args.announcement)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}

			announcement, err := test.fetch(repo)

			assert.Equal(t, test.expectedFetch.announcement, announcement, "Expecting announcement to be equal")
			assert.Equal(t, test.expectedFetch.err, err, "Expecting error to be equal")
		})
	}
}

func TestUpdateAnnouncement(t *testing.T) {
	type args struct {
		ctx          context.Context
		announcement *entity.Announcement
	}

	dummyUUID := uuid.New()

	dummyUpdatedAt, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", "2024-09-11 12:00:00 +0700 WIB")

	tests := []struct {
		name       string
		args       args
		beforeTest func(args args, mock sqlmock.Sqlmock)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "When updating existing announcement, it should not return an error",
			args: args{
				context.TODO(),
				&entity.Announcement{
					ID:          1,
					Description: "Test Announcement",
					AdminID:     dummyUUID,
				},
			},
			wantErr: false,
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE \"announcements\" SET .+`).
					WithArgs(
						args.announcement.AdminID,
						args.announcement.Description,
						sqlmock.AnyArg(),
						args.announcement.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "When updating non-existing announcement, it should return error item not found",
			args: args{
				context.TODO(),
				&entity.Announcement{
					ID:          1,
					Description: "Test Announcement",
					AdminID:     dummyUUID,
					UpdatedAt:   dummyUpdatedAt,
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE \"announcements\" SET .+`).
					WithArgs(
						args.announcement.AdminID,
						args.announcement.Description,
						sqlmock.AnyArg(),
						args.announcement.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectCommit()
			},
		},
		{
			name: "When updating data with rows affected more than 1, it should return weird behavior error",
			args: args{
				context.TODO(),
				&entity.Announcement{
					ID:          1,
					Description: "Test Announcement",
					UpdatedAt:   dummyUpdatedAt,
					AdminID:     dummyUUID,
				},
			},
			wantErr:    true,
			wantErrMsg: "weird behaviour. rows affected : 2",
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE \"announcements\" SET .+`).
					WithArgs(
						args.announcement.AdminID,
						args.announcement.Description,
						sqlmock.AnyArg(),
						args.announcement.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 2))

				mock.ExpectCommit()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			test.beforeTest(test.args, mock)

			repo := repository.NewAnnouncementRepository(db)

			err := repo.UpdateAnnouncement(test.args.ctx, test.args.announcement)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error to be thrown")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}
		})
	}
}

func TestDeleteAnnouncement(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(args args, mock sqlmock.Sqlmock)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "When dleting existing announcement, it should not return error",
			args: args{
				context.TODO(),
				1,
			},
			wantErr: false,
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`DELETE FROM "announcements" WHERE "announcements"."id" = \$1`).
					WithArgs(
						args.id,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "When deleting non-existing announcement, it should return error not found",
			args: args{
				context.TODO(),
				1,
			},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`DELETE FROM "announcements" WHERE "announcements"."id" = \$1`).
					WithArgs(
						args.id,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectCommit()
			},
		},
		{
			name: "When deleting with result rows affected more than 1, it should return error weird behaviour",
			args: args{
				context.TODO(),
				1,
			},
			wantErr:    true,
			wantErrMsg: "weird behaviour. rows affected : 2",
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`DELETE FROM "announcements" WHERE "announcements"."id" = \$1`).
					WithArgs(
						args.id,
					).
					WillReturnResult(sqlmock.NewResult(0, 2))

				mock.ExpectCommit()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			test.beforeTest(test.args, mock)

			repo := repository.NewAnnouncementRepository(db)

			err := repo.DeleteAnnouncement(test.args.ctx, test.args.id)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}
		})
	}
}

package test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/domain/enums"
	"github.com/hology8/hology-be/internal/app/competition/repository"
	pgsqlMock "github.com/hology8/hology-be/internal/infra/database/mock"
)

func TestFetchAll(t *testing.T) {
	db, close, mock := pgsqlMock.NewMockDB(t)

	defer close()

	repo := repository.NewCompetitionRepository(db)

	rows := sqlmock.NewRows([]string{"id", "competition_name", "competition_description"}).
		AddRow(1, "CTF", "Ini lomba ctf").
		AddRow(2, "DS", "ini lomba ds")

	mock.ExpectQuery(`SELECT \* FROM "competitions"`).WillReturnRows(rows)

	compes, err := repo.FetchAll(context.TODO())

	if err != nil {
		t.Fatal(err)
	}

	expected := []entity.Competition{
		{ID: 1, Name: "CTF", Desc: "Ini lomba ctf"},
		{ID: 2, Name: "DS", Desc: "ini lomba ds"},
	}

	assert.Equal(t, expected, compes, "Expecting same competitions result")
}

func TestFetchOneByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	tests := []struct {
		name        string
		args        args
		beforeTests func(sqlmock.Sqlmock)
		wantErr     bool
		wantErrMsg  string
		want        entity.Competition
	}{
		{
			name: "When fetching existing competition data , it should return the competition data",
			args: args{context.TODO(), 1},
			beforeTests: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{"id", "competition_name"}).
					AddRow(1, "Competition 1").
					AddRow(2, "Competition 2")

				mock.ExpectQuery(`SELECT \* FROM "competitions" WHERE "competitions"."id" = \$1 ORDER BY "competitions"."id" LIMIT \$2`).
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
			want:    entity.Competition{ID: 1, Name: "Competition 1"},
			wantErr: false,
		},
		{
			name: "When fetching non-existing competition data , it should return error",
			args: args{context.TODO(), 2},
			beforeTests: func(mock sqlmock.Sqlmock) {
				emptyRows := mock.NewRows([]string{"id", "competition_name"})

				mock.ExpectQuery(`SELECT \* FROM "competitions" WHERE "competitions"."id" = \$1 ORDER BY "competitions"."id" LIMIT \$2`).
					WithArgs(2, 1).
					WillReturnRows(emptyRows)
			},
			want:       entity.Competition{},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			repo := repository.NewCompetitionRepository(db)

			if test.beforeTests != nil {
				test.beforeTests(mock)
			}

			competition, err := repo.FetchOneByID(test.args.ctx, test.args.id)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error when no competition is not found")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting error message to be same")
			}

			assert.Equal(t, test.want, competition, "Expecting competition data to be equal")

		})
	}

}

func TestFetchOneWithRelations(t *testing.T) {
	type args struct {
		ctx       context.Context
		id        int
		relations []string
	}

	compeId := 1
	compeIdptr := &compeId

	tests := []struct {
		name        string
		args        args
		want        entity.Competition
		wantErr     bool
		wantErrMsg  string
		beforeTests func(mock sqlmock.Sqlmock)
	}{
		{
			name: "When fetching competitions with teams relation, it should return the competition with teams data",
			args: args{
				context.TODO(),
				1,
				[]string{"Teams"},
			},
			want: entity.Competition{
				ID:   1,
				Name: "CTF",
				Desc: "lomba ctf",
				Teams: []entity.Team{
					{
						ID:               "1",
						LeaderID:         uuid.MustParse("d20497aa-636c-496d-92d9-239fae6ab307"),
						CompetitionID:    1,
						Name:             "AcRtf",
						JoinToken:        "Jointknteam1",
						PaymentProofLink: "link.com",
						TwibbonProofLink: "link.com",
						Status:           enums.Verified,
						Phase:            enums.Elimination,
					},
					{
						ID:               "2",
						LeaderID:         uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
						CompetitionID:    1,
						Name:             "ctf!",
						JoinToken:        "Jointknteam2",
						PaymentProofLink: "link.com",
						TwibbonProofLink: "link.com",
						Status:           enums.Verified,
						Phase:            enums.Elimination,
					},
				},
			},
			wantErr: false,
			beforeTests: func(mock sqlmock.Sqlmock) {
				compRows := mock.NewRows([]string{"id", "competition_name", "competition_description"}).
					AddRow(1, "CTF", "lomba ctf")

				mock.ExpectQuery(`SELECT \* FROM "competitions" WHERE "competitions"."id" = \$1`).
					WillReturnRows(compRows)

				teamRows := mock.NewRows(
					[]string{
						"id",
						"leader_id",
						"competition_id",
						"team_name",
						"join_token",
						"payment_proof_link",
						"twibbon_proof_link",
						"status",
						"phase",
					},
				).AddRow(
					"1",
					uuid.MustParse("d20497aa-636c-496d-92d9-239fae6ab307"),
					1,
					"AcRtf",
					"Jointknteam1",
					"link.com",
					"link.com",
					"Verified",
					"Elimination",
				).AddRow(
					"2",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					1,
					"ctf!",
					"Jointknteam2",
					"link.com",
					"link.com",
					"Verified",
					"Elimination",
				)

				mock.ExpectQuery(`SELECT \* FROM "teams" WHERE "teams"."competition_id" = \$1`).
					WillReturnRows(teamRows)
			},
		},
		{
			name: "When fetching competitions with announcements relation, it should return the competition with announcements data",
			args: args{
				context.TODO(),
				1,
				[]string{"Announcements"},
			},
			want: entity.Competition{
				ID:   1,
				Name: "CTF",
				Desc: "lomba ctf",
				Announcements: []entity.Announcement{
					{
						ID:            1,
						AdminID:       uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
						CompetitionID: compeIdptr,
						Description:   "This is announcments 1 for this lomba",
					},
					{
						ID:            2,
						AdminID:       uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
						CompetitionID: compeIdptr,
						Description:   "This is announcments 2 for this lomba",
					},
				},
			},
			wantErr: false,
			beforeTests: func(mock sqlmock.Sqlmock) {
				compRows := mock.NewRows([]string{"id", "competition_name", "competition_description"}).
					AddRow(1, "CTF", "lomba ctf")

				mock.ExpectQuery(`SELECT \* FROM "competitions" WHERE "competitions"."id" = \$1`).
					WillReturnRows(compRows)

				announceRows := mock.NewRows(
					[]string{
						"id",
						"admin_id",
						"competition_id",
						"description",
					},
				).AddRow(
					1,
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					1,
					"This is announcments 1 for this lomba",
				).AddRow(
					2,
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					1,
					"This is announcments 2 for this lomba",
				)

				mock.ExpectQuery(`SELECT \* FROM "announcements" WHERE "announcements"."competition_id" = \$1`).
					WillReturnRows(announceRows)
			},
		},
		{
			name:       "When fetching non existing competition, it should return error item not found",
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
			beforeTests: func(mock sqlmock.Sqlmock) {
				compRows := mock.NewRows([]string{"id", "competition_name", "competition_description"})

				mock.ExpectQuery(`SELECT \* FROM "competitions"`).
					WillReturnRows(compRows)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			repo := repository.NewCompetitionRepository(db)

			if test.beforeTests != nil {
				test.beforeTests(mock)
			}

			compe, err := repo.FetchOneWithRelations(test.args.ctx, test.args.id, test.args.relations...)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")

			} else {
				assert.Nil(t, err, "error should not be expected ")
			}

			assert.Equal(t, test.want, compe, "Expecting result to be match with want")
		})
	}
}

func TestInsertCompe(t *testing.T) {
	type args struct {
		ctx   context.Context
		compe *entity.Competition
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(args args, sqlmock sqlmock.Sqlmock)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "When inserting a competition data, it should not return error",
			args: args{
				context.TODO(),
				&entity.Competition{
					Name:   "CTF",
					Desc:   "Ini deskripsi",
					LinkWA: "ini link",
				},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "competitions" (.+) RETURNING`).
					WithArgs(args.compe.Name, args.compe.Desc, args.compe.LinkWA).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			if test.beforeTest != nil {
				test.beforeTest(test.args, mock)
			}

			repo := repository.NewCompetitionRepository(db)

			err := repo.InsertCompe(test.args.ctx, test.args.compe)

			if !test.wantErr {
				assert.Nil(t, err, "Error should not be expected")
			} else {
				assert.NotNil(t, err, "Expecting error to be thrown")
			}
		})
	}
}

func TestUpdateCompe(t *testing.T) {
	type args struct {
		ctx   context.Context
		compe *entity.Competition
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(args args, sqlmock sqlmock.Sqlmock)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "When updating a existing competition data, it should not return error",
			args: args{
				context.TODO(),
				&entity.Competition{
					ID:     1,
					Name:   "CTF",
					Desc:   "Ini deskripsi",
					LinkWA: "ini link",
				},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE \"competitions\" SET .+`).
					WithArgs(
						args.compe.Name,
						args.compe.Desc,
						args.compe.LinkWA,
						args.compe.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "When updating a non-existing competition data, it should return error",
			args: args{
				context.TODO(),
				&entity.Competition{
					ID:     1,
					Name:   "CTF",
					Desc:   "Ini deskripsi",
					LinkWA: "ini link",
				},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE \"competitions\" SET .+`).
					WithArgs(
						args.compe.Name,
						args.compe.Desc,
						args.compe.LinkWA,
						args.compe.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectCommit()
			},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
		},
		{
			name: "When updating data but rows affected more than 1, it should return weird behaviour",
			args: args{
				context.TODO(),
				&entity.Competition{
					ID:     1,
					Name:   "CTF",
					Desc:   "Ini deskripsi",
					LinkWA: "ini link",
				},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE \"competitions\" SET .+`).
					WithArgs(
						args.compe.Name,
						args.compe.Desc,
						args.compe.LinkWA,
						args.compe.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 2))

				mock.ExpectCommit()
			},
			wantErr:    true,
			wantErrMsg: "weird behaviour. rows affected : 2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			if test.beforeTest != nil {
				test.beforeTest(test.args, mock)
			}

			repo := repository.NewCompetitionRepository(db)

			err := repo.UpdateCompe(test.args.ctx, test.args.compe)

			if !test.wantErr {
				assert.Nil(t, err, "Error should not be expected")
			} else {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error message")
			}
		})
	}
}

func TestDeleteCompe(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(args args, sqlmock sqlmock.Sqlmock)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "When updating a existing competition data, it should not return error",
			args: args{
				context.TODO(),
				1,
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`DELETE FROM "competitions" WHERE "competitions"."id" = \$1`).
					WithArgs().
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "When updating a non-existing competition data, it should return error",
			args: args{
				context.TODO(),
				1,
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`DELETE FROM "competitions" WHERE "competitions"."id" = \$1`).
					WithArgs(
						args.id,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectCommit()
			},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
		},
		{
			name: "When updating data but rows affected more than 1, it should return weird behaviour",
			args: args{
				context.TODO(),
				1,
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`DELETE FROM "competitions" WHERE "competitions"."id" = \$1`).
					WithArgs(
						args.id,
					).
					WillReturnResult(sqlmock.NewResult(0, 2))

				mock.ExpectCommit()
			},
			wantErr:    true,
			wantErrMsg: "weird behaviour. rows affected : 2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			if test.beforeTest != nil {
				test.beforeTest(test.args, mock)
			}

			repo := repository.NewCompetitionRepository(db)

			err := repo.DeleteCompe(test.args.ctx, test.args.id)

			if !test.wantErr {
				assert.Nil(t, err, "Error should not be expected")
			} else {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error message")
			}
		})
	}
}

func TestFetchAllTeamWithRelationAndCondition(t *testing.T) {
	type args struct {
		ctx       context.Context
		condition string
		args      []interface{}
		preload   []string
	}

	tests := []struct {
		name        string
		args        args
		beforeTest  func(sqlmock.Sqlmock)
		want        []entity.Competition
		wantErr     bool
		expectedErr error
	}{
		{
			name: "When fetching all competition without teams relation and condition, it should return the competition data",
			args: args{
				context.TODO(),
				"",
				[]interface{}{},
				[]string{},
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				compeRows := mock.NewRows([]string{"id", "competition_name", "competition_description"}).
					AddRow(1, "CTF", "Ini lomba ctf").
					AddRow(2, "DS", "ini lomba ds")

				mock.ExpectQuery(`SELECT \* FROM "competitions"`).WillReturnRows(compeRows)
			},
			want: []entity.Competition{
				{
					ID:   1,
					Name: "CTF",
					Desc: "Ini lomba ctf",
				},
				{
					ID:   2,
					Name: "DS",
					Desc: "ini lomba ds",
				},
			},
			wantErr: false,
		},
		{
			name: "When fetching all competition with teams relation and condition, it should return the competition data",
			args: args{
				context.TODO(),
				"competition_name = ?",
				[]interface{}{"CTF"},
				[]string{"Teams"},
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				compeRows := mock.NewRows([]string{"id", "competition_name", "competition_description"}).
					AddRow(1, "CTF", "Ini lomba ctf")

				mock.ExpectQuery(`SELECT \* FROM "competitions" WHERE (.+)`).
					WithArgs("CTF").
					WillReturnRows(compeRows)

				teamRows := mock.NewRows(
					[]string{
						"id",
						"leader_id",
						"competition_id",
						"team_name",
						"join_token",
						"payment_proof_link",
						"twibbon_proof_link",
						"status",
						"phase",
					},
				).AddRow(
					"1",
					uuid.MustParse("d20497aa-636c-496d-92d9-239fae6ab307"),
					1,
					"AcRtf",
					"Jointknteam1",
					"link.com",
					"link.com",
					"Verified",
					"Elimination",
				).AddRow(
					"2",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					1,
					"ctf!",
					"Jointknteam2",
					"link.com",
					"link.com",
					"Verified",
					"Elimination",
				)

				mock.ExpectQuery(`SELECT \* FROM "teams" WHERE (.+)`).
					WithArgs(1).
					WillReturnRows(teamRows)
			},
			want: []entity.Competition{
				{
					ID:   1,
					Name: "CTF",
					Desc: "Ini lomba ctf",
					Teams: []entity.Team{
						{
							ID:               "1",
							LeaderID:         uuid.MustParse("d20497aa-636c-496d-92d9-239fae6ab307"),
							CompetitionID:    1,
							Name:             "AcRtf",
							JoinToken:        "Jointknteam1",
							PaymentProofLink: "link.com",
							TwibbonProofLink: "link.com",
							Status:           enums.Verified,
							Phase:            enums.Elimination,
						},
						{
							ID:               "2",
							LeaderID:         uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
							CompetitionID:    1,
							Name:             "ctf!",
							JoinToken:        "Jointknteam2",
							PaymentProofLink: "link.com",
							TwibbonProofLink: "link.com",
							Status:           enums.Verified,
							Phase:            enums.Elimination,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "When fetching all competition with announcements relation and condition but there's an error, it should return internal server error",
			args: args{
				context.TODO(),
				"competition_name = ?",
				[]interface{}{"CTF"},
				[]string{"Announcements"},
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "competitions" WHERE (.+)`).
					WithArgs("CTF").
					WillReturnError(domain.ErrInternalServer)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "When fetching all competition with announcements relation and condition but it return 0 data, it should return not found error",
			args: args{
				context.TODO(),
				"competition_name = ?",
				[]interface{}{"CTF"},
				[]string{"Announcements"},
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				compeRows := mock.NewRows([]string{"id", "competition_name", "competition_description"})

				mock.ExpectQuery(`SELECT \* FROM "competitions" WHERE (.+)`).
					WithArgs("CTF").
					WillReturnRows(compeRows)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			if tt.beforeTest != nil {
				tt.beforeTest(mock)
			}

			repo := repository.NewCompetitionRepository(db)

			got, err := repo.FetchAllByConditionAndRelation(tt.args.ctx, tt.args.condition, tt.args.args, tt.args.preload...)

			if tt.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, tt.expectedErr, err, "Expecting same error message")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}

			assert.Equal(t, tt.want, got, "Expecting same result")
		})
	}
}

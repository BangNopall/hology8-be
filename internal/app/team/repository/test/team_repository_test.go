package test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/domain/enums"
	"github.com/hology8/hology-be/internal/app/team/repository"
	pgsqlMock "github.com/hology8/hology-be/internal/infra/database/mock"
)

func TestFetchOneByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}

	tests := []struct {
		name        string
		args        args
		beforeTests func(sqlmock.Sqlmock)
		wantErr     bool
		wantErrMsg  string
		want        entity.Team
	}{
		{
			name: "When fetching existing team data , it should return the team data",
			args: args{context.TODO(), "1"},
			beforeTests: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{"id", "team_name"}).
					AddRow("1", "Team 1").
					AddRow("2", "Team 2")

				mock.ExpectQuery(`SELECT \* FROM "teams" WHERE "teams"."id" = \$1 ORDER BY "teams"."id" LIMIT \$2`).
					WithArgs("1", 1).
					WillReturnRows(rows)
			},
			want:    entity.Team{ID: "1", Name: "Team 1"},
			wantErr: false,
		},
		{
			name: "When fetching non-existing team data , it should return error",
			args: args{context.TODO(), "2"},
			beforeTests: func(mock sqlmock.Sqlmock) {
				emptyRows := mock.NewRows([]string{"id", "team_name"})

				mock.ExpectQuery(`SELECT \* FROM "teams" WHERE "teams"."id" = \$1 ORDER BY "teams"."id" LIMIT \$2`).
					WithArgs("2", 1).
					WillReturnRows(emptyRows)
			},
			want:       entity.Team{},
			wantErr:    true,
			wantErrMsg: domain.ErrNotFound.Error(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			repo := repository.NewTeamRepository(db)

			if test.beforeTests != nil {
				test.beforeTests(mock)
			}

			team, err := repo.FetchOneByID(test.args.ctx, test.args.id)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error when no team is not found")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting error message to be same")
			}

			assert.Equal(t, test.want, team, "Expecting team data to be equal")

		})
	}

}

func TestFetchAllTeams(t *testing.T) {

	teamId1 := "1"
	teamId2 := "2"
	teamId1Ptr := &teamId1
	teamId2Ptr := &teamId2

	type args struct {
		ctx       context.Context
		params    *dto.TeamParams
		relations []string
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(args args, mock sqlmock.Sqlmock)
		want       interface{}
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "When fetching all teams, it should return all teams",
			args: args{
				context.TODO(),
				&dto.TeamParams{},
				[]string{},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				rows := mock.NewRows(
					[]string{
						"id",
						"leader_id",
						"competition_id",
						"team_name",
						"join_token",
						"payment_proof_link",
						"twibbon_proof_link",
						"proposal_doc_link",
						"video_link",
						"statement_letter_link",
						"status",
						"phase",
					},
				).AddRow(
					"1",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					1,
					"AcRtf",
					"joinToken1",
					"proof_link1.com",
					"twibbon_link1.com",
					"drive.google.com/doc/1",
					"youtube.com/watch?node=426",
					"hology.id",
					"Verified",
					"Elimination",
				).AddRow(
					"2",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					3,
					"Ingloryy rawr",
					"joinToken2",
					"proof_link2.com",
					"twibbon_link2.com",
					"drive.google.com/doc/2",
					"youtube.com/watch?node=466",
					"hology.id",
					"Verified",
					"Final",
				)

				mock.ExpectQuery(`SELECT \* FROM "teams"`).WillReturnRows(rows)
			},
			want: []entity.Team{
				{
					ID:                  "1",
					LeaderID:            uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					CompetitionID:       1,
					Name:                "AcRtf",
					JoinToken:           "joinToken1",
					PaymentProofLink:    "proof_link1.com",
					TwibbonProofLink:    "twibbon_link1.com",
					ProposalDocLink:     "drive.google.com/doc/1",
					StatementLetterLink: "hology.id",
					VideoLink:           "youtube.com/watch?node=426",
					Status:              "Verified",
					Phase:               "Elimination",
				},
				{
					ID:                  "2",
					LeaderID:            uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					CompetitionID:       3,
					Name:                "Ingloryy rawr",
					JoinToken:           "joinToken2",
					PaymentProofLink:    "proof_link2.com",
					TwibbonProofLink:    "twibbon_link2.com",
					ProposalDocLink:     "drive.google.com/doc/2",
					StatementLetterLink: "hology.id",
					VideoLink:           "youtube.com/watch?node=466",
					Status:              "Verified",
					Phase:               "Final",
				},
			},
			wantErr: false,
		},
		{
			name: "When fetching with several params, it should return teams filtered by params",
			args: args{
				context.TODO(),
				&dto.TeamParams{
					CompetitionID: 1,
					Status:        enums.Verified,
					Phase:         enums.Final,
				},
				[]string{},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				rows := mock.NewRows(
					[]string{
						"id",
						"leader_id",
						"competition_id",
						"team_name",
						"join_token",
						"payment_proof_link",
						"twibbon_proof_link",
						"proposal_doc_link",
						"video_link",
						"statement_letter_link",
						"status",
						"phase",
					},
				).AddRow(
					"1",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					1,
					"AcRtf",
					"joinToken1",
					"proof_link1.com",
					"twibbon_link1.com",
					"drive.google.com/doc/1",
					"youtube.com/watch?node=426",
					"hology.id",
					"Verified",
					"Final",
				).AddRow(
					"2",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					1,
					"Ingloryy rawr",
					"joinToken2",
					"proof_link2.com",
					"twibbon_link2.com",
					"drive.google.com/doc/2",
					"youtube.com/watch?node=466",
					"hology.id",
					"Verified",
					"Final",
				)

				mock.ExpectQuery(`SELECT \* FROM "teams" WHERE (.+)`).WillReturnRows(rows)
			},
			want: []entity.Team{
				{
					ID:                  "1",
					LeaderID:            uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					CompetitionID:       1,
					Name:                "AcRtf",
					JoinToken:           "joinToken1",
					PaymentProofLink:    "proof_link1.com",
					TwibbonProofLink:    "twibbon_link1.com",
					ProposalDocLink:     "drive.google.com/doc/1",
					VideoLink:           "youtube.com/watch?node=426",
					StatementLetterLink: "hology.id",
					Status:              "Verified",
					Phase:               "Final",
				},
				{
					ID:                  "2",
					LeaderID:            uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					CompetitionID:       1,
					Name:                "Ingloryy rawr",
					JoinToken:           "joinToken2",
					PaymentProofLink:    "proof_link2.com",
					TwibbonProofLink:    "twibbon_link2.com",
					ProposalDocLink:     "drive.google.com/doc/2",
					VideoLink:           "youtube.com/watch?node=466",
					StatementLetterLink: "hology.id",
					Status:              "Verified",
					Phase:               "Final",
				},
			},
			wantErr: false,
		},
		{
			name: "When fetching with several params with relations param, it should return teams filtered by params and preload all the relations",
			args: args{
				context.TODO(),
				&dto.TeamParams{
					CompetitionID: 1,
					Status:        enums.Verified,
					Phase:         enums.Final,
				},
				[]string{"Announcements", "Members", "Competition"},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				teamRows := mock.NewRows(
					[]string{
						"id",
						"leader_id",
						"competition_id",
						"team_name",
						"join_token",
						"payment_proof_link",
						"twibbon_proof_link",
						"proposal_doc_link",
						"video_link",
						"statement_letter_link",
						"status",
						"phase",
					},
				).AddRow(
					"1",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					1,
					"AcRtf",
					"joinToken1",
					"proof_link1.com",
					"twibbon_link1.com",
					"drive.google.com/doc/1",
					"youtube.com/watch?node=426",
					"hology.id",
					"Verified",
					"Final",
				).AddRow(
					"2",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					2,
					"Ingloryy rawr",
					"joinToken2",
					"proof_link2.com",
					"twibbon_link2.com",
					"drive.google.com/doc/2",
					"youtube.com/watch?node=466",
					"hology.id",
					"Verified",
					"Final",
				)

				announcementRows := mock.NewRows(
					[]string{
						"id",
						"admin_id",
						"team_id",
						"description",
					},
				).AddRow(
					1,
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					"1",
					"Pengumuman untuk tim AcRtf",
				).AddRow(
					1,
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					"2",
					"Pengumuman untuk tim Ingloryy rawr",
				)

				memberRows := mock.NewRows(
					[]string{
						"user_id",
						"team_id",
					},
				).AddRow(
					uuid.MustParse("c14c1f1c-e590-426c-b70b-18f8652af0c4"),
					"1",
				).AddRow(
					uuid.MustParse("35d0f316-9f7f-4df8-9c64-9bf2dfdf5342"),
					"1",
				).AddRow(
					uuid.MustParse("0bce5fc5-2e23-4f0b-9091-baff648d1cd6"),
					"2",
				).AddRow(
					uuid.MustParse("916489f6-f156-4cb9-97ae-d41847980f10"),
					"2",
				)

				compeRows := mock.
					NewRows([]string{"id", "competition_name", "competition_description"}).
					AddRow(
						1,
						"CTF",
						"Ini lomba ctf",
					).
					AddRow(
						2,
						"UI UX",
						"Ini lomba ui ux",
					)

				mock.ExpectQuery(`SELECT \* FROM "teams" WHERE (.+)`).
					WillReturnRows(teamRows)

				mock.ExpectQuery(`SELECT \* FROM "announcements" WHERE "announcements"\."team_id" IN \(\$1,\$2\)`).
					WillReturnRows(announcementRows)

				mock.ExpectQuery(`SELECT \* FROM "competitions" WHERE "competitions"\."id" IN \(\$1,\$2\)`).
					WillReturnRows(compeRows)

				mock.ExpectQuery(`SELECT \* FROM "detail_teams" WHERE "detail_teams"\."team_id" IN \(\$1,\$2\)`).
					WillReturnRows(memberRows)

			},
			want: []entity.Team{
				{
					ID:                  "1",
					LeaderID:            uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					CompetitionID:       1,
					Name:                "AcRtf",
					JoinToken:           "joinToken1",
					PaymentProofLink:    "proof_link1.com",
					TwibbonProofLink:    "twibbon_link1.com",
					ProposalDocLink:     "drive.google.com/doc/1",
					VideoLink:           "youtube.com/watch?node=426",
					StatementLetterLink: "hology.id",
					Status:              "Verified",
					Phase:               "Final",
					Announcements: []entity.Announcement{
						{
							ID:          1,
							AdminID:     uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
							TeamID:      teamId1Ptr,
							Description: "Pengumuman untuk tim AcRtf",
						},
					},
					Members: []entity.DetailTeams{
						{
							UserID: uuid.MustParse("c14c1f1c-e590-426c-b70b-18f8652af0c4"),
							TeamID: "1",
						},
						{
							UserID: uuid.MustParse("35d0f316-9f7f-4df8-9c64-9bf2dfdf5342"),
							TeamID: "1",
						},
					},
					Competition: entity.Competition{
						ID:   1,
						Name: "CTF",
						Desc: "Ini lomba ctf",
					},
				},
				{
					ID:                  "2",
					LeaderID:            uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					CompetitionID:       2,
					Name:                "Ingloryy rawr",
					JoinToken:           "joinToken2",
					PaymentProofLink:    "proof_link2.com",
					TwibbonProofLink:    "twibbon_link2.com",
					ProposalDocLink:     "drive.google.com/doc/2",
					VideoLink:           "youtube.com/watch?node=466",
					StatementLetterLink: "hology.id",
					Status:              "Verified",
					Phase:               "Final",
					Announcements: []entity.Announcement{
						{
							ID:          1,
							AdminID:     uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
							TeamID:      teamId2Ptr,
							Description: "Pengumuman untuk tim Ingloryy rawr",
						},
					},
					Members: []entity.DetailTeams{
						{
							UserID: uuid.MustParse("0bce5fc5-2e23-4f0b-9091-baff648d1cd6"),
							TeamID: "2",
						},
						{
							UserID: uuid.MustParse("916489f6-f156-4cb9-97ae-d41847980f10"),
							TeamID: "2",
						},
					},
					Competition: entity.Competition{
						ID:   2,
						Name: "UI UX",
						Desc: "Ini lomba ui ux",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			repo := repository.NewTeamRepository(db)

			if test.beforeTest != nil {
				test.beforeTest(test.args, mock)
			}

			teams, err := repo.FetchAllTeams(test.args.ctx, test.args.params, test.args.relations...)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}

			assert.Equal(t, test.want, teams, "Expecting teams result to be equal")
		})
	}

}

func TestFetchOneByParams(t *testing.T) {
	teamId1 := "1"
	teamId1Ptr := &teamId1

	type args struct {
		ctx       context.Context
		params    *dto.TeamParams
		relations []string
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(args args, mock sqlmock.Sqlmock)
		want       interface{}
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "When fetching team with id param, it should return a team",
			args: args{
				context.TODO(),
				&dto.TeamParams{ID: "1"},
				[]string{},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				rows := mock.NewRows(
					[]string{
						"id",
						"leader_id",
						"competition_id",
						"team_name",
						"join_token",
						"payment_proof_link",
						"twibbon_proof_link",
						"proposal_doc_link",
						"video_link",
						"statement_letter_link",
						"status",
						"phase",
					},
				).AddRow(
					"1",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					1,
					"AcRtf",
					"joinToken1",
					"proof_link1.com",
					"twibbon_link1.com",
					"drive.google.com/doc/1",
					"youtube.com/watch?node=426",
					"hology.id",
					"Verified",
					"Elimination",
				)

				mock.ExpectQuery(`SELECT \* FROM "teams" WHERE (.+)`).
					WithArgs(args.params.ID).
					WillReturnRows(rows)
			},
			want: entity.Team{
				ID:                  "1",
				LeaderID:            uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
				CompetitionID:       1,
				Name:                "AcRtf",
				JoinToken:           "joinToken1",
				PaymentProofLink:    "proof_link1.com",
				TwibbonProofLink:    "twibbon_link1.com",
				ProposalDocLink:     "drive.google.com/doc/1",
				VideoLink:           "youtube.com/watch?node=426",
				StatementLetterLink: "hology.id",
				Status:              "Verified",
				Phase:               "Elimination",
			},
			wantErr: false,
		},
		{
			name: "When fetching with several params, it should return teams filtered by params",
			args: args{
				context.TODO(),
				&dto.TeamParams{
					Name:          "AcRtf",
					CompetitionID: 1,
				},
				[]string{},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				rows := mock.NewRows(
					[]string{
						"id",
						"leader_id",
						"competition_id",
						"team_name",
						"join_token",
						"payment_proof_link",
						"twibbon_proof_link",
						"proposal_doc_link",
						"video_link",
						"statement_letter_link",
						"status",
						"phase",
					},
				).AddRow(
					"1",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					1,
					"AcRtf",
					"joinToken1",
					"proof_link1.com",
					"twibbon_link1.com",
					"drive.google.com/doc/1",
					"youtube.com/watch?node=426",
					"hology.id",
					"Verified",
					"Final",
				)

				mock.ExpectQuery(`SELECT \* FROM "teams" WHERE (.+)`).WillReturnRows(rows)
			},
			want: entity.Team{
				ID:                  "1",
				LeaderID:            uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
				CompetitionID:       1,
				Name:                "AcRtf",
				JoinToken:           "joinToken1",
				PaymentProofLink:    "proof_link1.com",
				TwibbonProofLink:    "twibbon_link1.com",
				ProposalDocLink:     "drive.google.com/doc/1",
				VideoLink:           "youtube.com/watch?node=426",
				StatementLetterLink: "hology.id",
				Status:              "Verified",
				Phase:               "Final",
			},
			wantErr: false,
		},
		{
			name: "When fetching with several params with relations param, it should return teams filtered by params and preload all the relations",
			args: args{
				context.TODO(),
				&dto.TeamParams{
					CompetitionID: 1,
					Status:        enums.Verified,
					Phase:         enums.Final,
				},
				[]string{"Announcements", "Members", "Competition"},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				teamRows := mock.NewRows(
					[]string{
						"id",
						"leader_id",
						"competition_id",
						"team_name",
						"join_token",
						"payment_proof_link",
						"twibbon_proof_link",
						"proposal_doc_link",
						"video_link",
						"statement_letter_link",
						"status",
						"phase",
					},
				).AddRow(
					"1",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					1,
					"AcRtf",
					"joinToken1",
					"proof_link1.com",
					"twibbon_link1.com",
					"drive.google.com/doc/1",
					"youtube.com/watch?node=426",
					"hology.id",
					"Verified",
					"Final",
				)

				announcementRows := mock.NewRows(
					[]string{
						"id",
						"admin_id",
						"team_id",
						"description",
					},
				).AddRow(
					1,
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					"1",
					"Pengumuman untuk tim AcRtf",
				)

				memberRows := mock.NewRows(
					[]string{
						"user_id",
						"team_id",
					},
				).AddRow(
					uuid.MustParse("c14c1f1c-e590-426c-b70b-18f8652af0c4"),
					"1",
				).AddRow(
					uuid.MustParse("35d0f316-9f7f-4df8-9c64-9bf2dfdf5342"),
					"1",
				)

				compeRows := mock.
					NewRows([]string{"id", "competition_name", "competition_description"}).
					AddRow(
						1,
						"CTF",
						"Ini lomba ctf",
					)

				mock.ExpectQuery(`SELECT \* FROM "teams" WHERE (.+)`).
					WillReturnRows(teamRows)

				mock.ExpectQuery(`SELECT \* FROM "announcements" WHERE (.+)`).
					WillReturnRows(announcementRows)

				mock.ExpectQuery(`SELECT \* FROM "competitions" WHERE (.+)`).
					WillReturnRows(compeRows)

				mock.ExpectQuery(`SELECT \* FROM "detail_teams" WHERE (.+)`).
					WillReturnRows(memberRows)

			},
			want: entity.Team{
				ID:                  "1",
				LeaderID:            uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
				CompetitionID:       1,
				Name:                "AcRtf",
				JoinToken:           "joinToken1",
				PaymentProofLink:    "proof_link1.com",
				TwibbonProofLink:    "twibbon_link1.com",
				ProposalDocLink:     "drive.google.com/doc/1",
				VideoLink:           "youtube.com/watch?node=426",
				StatementLetterLink: "hology.id",
				Status:              "Verified",
				Phase:               "Final",
				Announcements: []entity.Announcement{
					{
						ID:          1,
						AdminID:     uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
						TeamID:      teamId1Ptr,
						Description: "Pengumuman untuk tim AcRtf",
					},
				},
				Members: []entity.DetailTeams{
					{
						UserID: uuid.MustParse("c14c1f1c-e590-426c-b70b-18f8652af0c4"),
						TeamID: "1",
					},
					{
						UserID: uuid.MustParse("35d0f316-9f7f-4df8-9c64-9bf2dfdf5342"),
						TeamID: "1",
					},
				},
				Competition: entity.Competition{
					ID:   1,
					Name: "CTF",
					Desc: "Ini lomba ctf",
				},
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			repo := repository.NewTeamRepository(db)

			if test.beforeTest != nil {
				test.beforeTest(test.args, mock)
			}

			team, err := repo.FetchOneByParams(test.args.ctx, test.args.params, test.args.relations...)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}

			assert.Equal(t, test.want, team, "Expecting teams result to be equal")
		})
	}
}

func TestFetchTeamMember(t *testing.T) {
	type args struct {
		ctx    context.Context
		teamId string
		userId uuid.UUID
	}

	currTime := time.Now()

	tests := []struct {
		name       string
		args       args
		beforeTest func(args args, sqlmock sqlmock.Sqlmock)
		want       interface{}
		wantErr    bool
	}{
		{
			name: "When fetch an existing team member, it should return member with preloaded relationship",
			args: args{
				context.TODO(),
				"teamId77FF",
				uuid.MustParse("771aaa85-1421-4730-9ca0-77f551990081"),
			},
			beforeTest: func(args args, sqlmock sqlmock.Sqlmock) {
				detailRows := sqlmock.NewRows([]string{"user_id", "team_id", "created_at"}).
					AddRow(uuid.MustParse("771aaa85-1421-4730-9ca0-77f551990081"), "teamId77FF", currTime)

				userRows := sqlmock.NewRows([]string{"id", "email", "password", "fullname"}).
					AddRow(uuid.MustParse("771aaa85-1421-4730-9ca0-77f551990081"), "devan@gmail.com", "1234", "Devan")

				teamRows := sqlmock.NewRows([]string{"id", "team_name", "join_token"}).
					AddRow("teamId77FF", "SoftAf", "joinThisTeam")

				sqlmock.ExpectQuery(`SELECT \* FROM "detail_teams" WHERE (.+)`).
					WillReturnRows(detailRows)

				sqlmock.ExpectQuery(`SELECT \* FROM "teams" WHERE (.+)`).
					WillReturnRows(teamRows)

				sqlmock.ExpectQuery(`SELECT \* FROM "users" WHERE (.+)`).
					WillReturnRows(userRows)
			},
			want: entity.DetailTeams{
				UserID:    uuid.MustParse("771aaa85-1421-4730-9ca0-77f551990081"),
				TeamID:    "teamId77FF",
				CreatedAt: currTime,
				User: entity.User{
					ID:       uuid.MustParse("771aaa85-1421-4730-9ca0-77f551990081"),
					Email:    "devan@gmail.com",
					Password: "1234",
					Fullname: "Devan",
				},
				Team: entity.Team{
					ID:        "teamId77FF",
					Name:      "SoftAf",
					JoinToken: "joinThisTeam",
				},
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			repo := repository.NewTeamRepository(db)

			if test.beforeTest != nil {
				test.beforeTest(test.args, mock)
			}

			member, err := repo.FetchTeamMember(test.args.ctx, test.args.teamId, test.args.userId)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}

			assert.Equal(t, test.want, member, "Expecting result to be equal")
		})
	}
}

func TestInsertTeam(t *testing.T) {
	type args struct {
		ctx  context.Context
		team *entity.Team
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(args args, sqlmock sqlmock.Sqlmock)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "When inserting a team data, it should not return error",
			args: args{
				context.TODO(),
				&entity.Team{
					ID:                  "2",
					LeaderID:            uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					CompetitionID:       2,
					Name:                "Ingloryy rawr",
					JoinToken:           "joinToken2",
					SenderPaymentName:   "Devan",
					BankAccountNumber:   "1234567890",
					BankName:            "BCA",
					PaymentProofLink:    "proof_link2.com",
					TwibbonProofLink:    "twibbon_link.com",
					ProposalDocLink:     "drive.google.com",
					VideoLink:           "youtube.com",
					PPTLink:             "ppt.com",
					StatementLetterLink: "hology.id",
					Status:              enums.Waiting,
					Phase:               enums.Elimination,
					WinnerPlace:         enums.Default,
					CreatedAt:           time.Now(),
					UpdatedAt:           time.Now(),
				},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO \"teams\" (.+) VALUES (.+) RETURNING \"university_id\"").
					WithArgs(
						args.team.ID,
						args.team.LeaderID,
						args.team.CompetitionID,
						args.team.Name,
						args.team.JoinToken,
						args.team.SenderPaymentName,
						args.team.BankAccountNumber,
						args.team.BankName,
						args.team.PaymentProofLink,
						args.team.TwibbonProofLink,
						args.team.ProposalDocLink,
						args.team.StatementLetterLink,
						args.team.VideoLink,
						args.team.PPTLink,
						args.team.Status,
						args.team.Phase,
						args.team.WinnerPlace,
						args.team.CreatedAt,
						args.team.UpdatedAt,
					).
					WillReturnRows(sqlmock.NewRows([]string{"university_id"}).AddRow(args.team.UniversityID)).
					WillReturnError(nil)

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

			repo := repository.NewTeamRepository(db)

			err := repo.InsertTeam(test.args.ctx, test.args.team)

			if !test.wantErr {
				assert.Nil(t, err, "Error should not be expected")
			} else {
				assert.NotNil(t, err, "Expecting error to be thrown")
			}
		})
	}
}

func TestInsertTeamMember(t *testing.T) {
	type args struct {
		ctx    context.Context
		member *entity.DetailTeams
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(args args, mock sqlmock.Sqlmock)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "When inserting a team member, it should not return error",
			args: args{
				context.TODO(),
				&entity.DetailTeams{
					UserID:    uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					TeamID:    "80",
					CreatedAt: time.Now(),
				},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "detail_teams" (.+) RETURNING`).
					WithArgs(
						args.member.UserID,
						args.member.TeamID,
						args.member.CreatedAt,
					).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(args.member.UserID))

				mock.ExpectCommit()
			},
		},
		{
			name: "When inserting a team member but the entry already exist, it should return error",
			args: args{
				context.TODO(),
				&entity.DetailTeams{
					UserID:    uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					TeamID:    "80",
					CreatedAt: time.Now(),
				},
			},
			wantErr:    true,
			wantErrMsg: domain.ErrDuplicateEntry.Error(),
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "detail_teams" (.+) RETURNING`).
					WithArgs(
						args.member.UserID,
						args.member.TeamID,
						args.member.CreatedAt,
					).
					WillReturnError(gorm.ErrDuplicatedKey)

				mock.ExpectRollback()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			if test.beforeTest != nil {
				test.beforeTest(test.args, mock)
			}

			repo := repository.NewTeamRepository(db)

			err := repo.InsertTeamMember(test.args.ctx, test.args.member)

			if !test.wantErr {
				assert.Nil(t, err, "Error should not be expected")
			} else {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecing same error message")
			}
		})
	}
}

func TestUpdateTeam(t *testing.T) {
	type args struct {
		ctx  context.Context
		id   string
		team *dto.TeamUpdate
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(args args, sqlmock sqlmock.Sqlmock)
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "When updating a existing team data, it should not return error",
			args: args{
				context.TODO(),
				"1",
				&dto.TeamUpdate{
					PaymentProofLink:    "proof_link2.com",
					TwibbonProofLink:    "twibbon_link2.com",
					ProposalDocLink:     "drive.google.com/doc/2",
					VideoLink:           "youtube.com/watch?node=466",
					StatementLetterLink: "hology.id",
					Status:              "Verified",
					Phase:               "Final",
				},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE \"teams\" SET .+`).
					WithArgs(
						args.team.PaymentProofLink,
						args.team.TwibbonProofLink,
						args.team.ProposalDocLink,
						args.team.StatementLetterLink,
						args.team.VideoLink,
						args.team.Status,
						args.team.Phase,
						args.id,
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
				"1",
				&dto.TeamUpdate{
					PaymentProofLink:    "proof_link2.com",
					TwibbonProofLink:    "twibbon_link2.com",
					ProposalDocLink:     "drive.google.com/doc/2",
					VideoLink:           "youtube.com/watch?node=466",
					StatementLetterLink: "hology.id",
					Status:              "Verified",
					Phase:               "Final",
				},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE \"teams\" SET .+`).
					WithArgs(
						args.team.PaymentProofLink,
						args.team.TwibbonProofLink,
						args.team.ProposalDocLink,
						args.team.StatementLetterLink,
						args.team.VideoLink,
						args.team.Status,
						args.team.Phase,
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
				"1",
				&dto.TeamUpdate{
					UniversityID:        1,
					SenderPaymentName:   "Devan",
					BankAccountNumber:   "299",
					BankName:            "BCA",
					PaymentProofLink:    "proof_link2.com",
					TwibbonProofLink:    "twibbon_link2.com",
					ProposalDocLink:     "drive.google.com/doc/2",
					VideoLink:           "youtube.com/watch?node=466",
					StatementLetterLink: "hology.id",
					Status:              "Verified",
					Phase:               "Final",
				},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE \"teams\" SET .+`).
					WithArgs(
						args.team.UniversityID,
						args.team.SenderPaymentName,
						args.team.BankAccountNumber,
						args.team.BankName,
						args.team.PaymentProofLink,
						args.team.TwibbonProofLink,
						args.team.ProposalDocLink,
						args.team.StatementLetterLink,
						args.team.VideoLink,
						args.team.Status,
						args.team.Phase,
						args.id,
					).
					WillReturnResult(sqlmock.NewResult(0, 2))

				mock.ExpectCommit()
			},
			wantErr:    true,
			wantErrMsg: "weird behaviour. rows affected : 2",
		},
		{
			name: "When updating data with no existing univeristy, it should return error university not found",
			args: args{
				context.TODO(),
				"1",
				&dto.TeamUpdate{
					UniversityID:        1,
					SenderPaymentName:   "Devan",
					BankAccountNumber:   "299",
					BankName:            "BCA",
					PaymentProofLink:    "proof_link2.com",
					TwibbonProofLink:    "twibbon_link2.com",
					ProposalDocLink:     "drive.google.com/doc/2",
					VideoLink:           "youtube.com/watch?node=466",
					StatementLetterLink: "hology.id",
					Status:              "Verified",
					Phase:               "Final",
				},
			},
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE \"teams\" SET .+`).
					WithArgs(
						args.team.UniversityID,
						args.team.SenderPaymentName,
						args.team.BankAccountNumber,
						args.team.BankName,
						args.team.PaymentProofLink,
						args.team.TwibbonProofLink,
						args.team.ProposalDocLink,
						args.team.StatementLetterLink,
						args.team.VideoLink,
						args.team.Status,
						args.team.Phase,
						args.id,
					).
					WillReturnError(gorm.ErrForeignKeyViolated)

				mock.ExpectRollback()
			},
			wantErr:    true,
			wantErrMsg: domain.ErrUniversityNotFound.Error(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			if test.beforeTest != nil {
				test.beforeTest(test.args, mock)
			}

			repo := repository.NewTeamRepository(db)

			err := repo.UpdateTeam(test.args.ctx, test.args.id, test.args.team)

			if !test.wantErr {
				assert.Nil(t, err, "Error should not be expected")
			} else {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.wantErrMsg, err.Error(), "Expecting same error message")
			}
		})
	}
}

func TestDeleteTeamMember(t *testing.T) {
	type args struct {
		ctx    context.Context
		member *entity.DetailTeams
	}

	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
		beforeTest  func(args args, mock sqlmock.Sqlmock)
	}{
		{
			name: "When deleting a member from a team, it should return no error",
			args: args{
				context.TODO(),
				&entity.DetailTeams{
					UserID: uuid.MustParse("d614a90b-202c-4560-ad07-ee342cb66e41"),
					TeamID: "77FF",
				},
			},
			wantErr: false,
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`DELETE FROM "detail_teams" WHERE (.+)`).
					WithArgs(args.member.TeamID, args.member.UserID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "When deleting a non existing member from a team, it should return error item not found",
			args: args{
				context.TODO(),
				&entity.DetailTeams{
					UserID: uuid.MustParse("d614a90b-202c-4560-ad07-ee342cb66e41"),
					TeamID: "77FF",
				},
			},
			wantErr:     true,
			expectedErr: domain.ErrNotFound,
			beforeTest: func(args args, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(`DELETE FROM "detail_teams" WHERE (.+)`).
					WithArgs(args.member.TeamID, args.member.UserID).
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectCommit()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			if test.beforeTest != nil {
				test.beforeTest(test.args, mock)
			}

			repo := repository.NewTeamRepository(db)

			err := repo.DeleteTeamMember(test.args.ctx, test.args.member)

			if test.wantErr {
				assert.NotNil(t, err, "Expected error to be thrown")
				assert.Equal(t, test.expectedErr, err, "Expected same error")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}
		})
	}
}

func TestFetchAllTeamWithRelationAndCondition(t *testing.T) {
	type args struct {
		ctx       context.Context
		condition string
		relations []string
		args      []interface{}
	}

	tests := []struct {
		name        string
		args        args
		beforeTest  func(mock sqlmock.Sqlmock)
		want        []entity.Team
		wantErr     bool
		expectedErr error
	}{
		{
			name: "When fetching all teams without relation and condition, it should return teams",
			args: args{
				context.TODO(),
				"",
				nil,
				[]interface{}{},
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				teamRows := mock.NewRows(
					[]string{
						"id",
						"leader_id",
						"competition_id",
						"team_name",
						"join_token",
						"payment_proof_link",
						"twibbon_proof_link",
						"statement_letter_link",
					},
				).AddRow(
					"1",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					1,
					"AcRtf",
					"joinToken1",
					"proof_link1.com",
					"twibbon_link1.com",
					"hology.id",
				).AddRow(
					"2",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					2,
					"Ingloryy rawr",
					"joinToken2",
					"proof_link2.com",
					"twibbon_link2.com",
					"hology.id",
				).AddRow(
					"3",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					3,
					"Ingloryy rawr",
					"joinToken3",
					"proof_link3.com",
					"twibbon_link3.com",
					"hology.id",
				)

				mock.ExpectQuery(`SELECT \* FROM "teams"`).
					WillReturnRows(teamRows)
			},
			want: []entity.Team{
				{
					ID:                  "1",
					LeaderID:            uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					CompetitionID:       1,
					Name:                "AcRtf",
					JoinToken:           "joinToken1",
					PaymentProofLink:    "proof_link1.com",
					TwibbonProofLink:    "twibbon_link1.com",
					StatementLetterLink: "hology.id",
				},
				{
					ID:                  "2",
					LeaderID:            uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					CompetitionID:       2,
					Name:                "Ingloryy rawr",
					JoinToken:           "joinToken2",
					PaymentProofLink:    "proof_link2.com",
					TwibbonProofLink:    "twibbon_link2.com",
					StatementLetterLink: "hology.id",
				},
				{
					ID:                  "3",
					LeaderID:            uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					CompetitionID:       3,
					Name:                "Ingloryy rawr",
					JoinToken:           "joinToken3",
					PaymentProofLink:    "proof_link3.com",
					TwibbonProofLink:    "twibbon_link3.com",
					StatementLetterLink: "hology.id",
				},
			},
			wantErr: false,
		},
		{
			name: "when fetching all teams with condition and relation, it should return teams and the relation",
			args: args{
				ctx:       context.TODO(),
				condition: "team_name = ?",
				args:      []interface{}{"AcRtf"},
				relations: []string{"Members"},
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				teamRows := mock.NewRows(
					[]string{
						"id",
						"leader_id",
						"competition_id",
						"team_name",
						"join_token",
						"payment_proof_link",
						"twibbon_proof_link",
						"statement_letter_link",
					},
				).AddRow(
					"1",
					uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					1,
					"AcRtf",
					"joinToken1",
					"proof_link1.com",
					"twibbon_link1.com",
					"hology.id",
				)

				memberRows := mock.NewRows(
					[]string{
						"user_id",
						"team_id",
					},
				).AddRow(
					uuid.MustParse("c14c1f1c-e590-426c-b70b-18f8652af0c4"),
					"1",
				).AddRow(
					uuid.MustParse("35d0f316-9f7f-4df8-9c64-9bf2dfdf5342"),
					"1",
				)

				mock.ExpectQuery(`SELECT \* FROM "teams" WHERE (.+)`).
					WillReturnRows(teamRows)

				mock.ExpectQuery(`SELECT \* FROM "detail_teams" WHERE (.+)`).
					WillReturnRows(memberRows)
			},
			want: []entity.Team{
				{
					ID:                  "1",
					LeaderID:            uuid.MustParse("ad172b18-f829-400f-8d75-cb16538c8503"),
					CompetitionID:       1,
					Name:                "AcRtf",
					JoinToken:           "joinToken1",
					PaymentProofLink:    "proof_link1.com",
					TwibbonProofLink:    "twibbon_link1.com",
					StatementLetterLink: "hology.id",
					Members: []entity.DetailTeams{
						{
							UserID: uuid.MustParse("c14c1f1c-e590-426c-b70b-18f8652af0c4"),
							TeamID: "1",
						},
						{
							UserID: uuid.MustParse("35d0f316-9f7f-4df8-9c64-9bf2dfdf5342"),
							TeamID: "1",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "when fetching all teams with condition and relation but there's an error, it should return internal server error",
			args: args{
				ctx:       context.TODO(),
				condition: "team_name = ?",
				args:      []interface{}{"AcRtf"},
				relations: []string{"Members"},
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "teams" WHERE (.+)`).
					WillReturnError(assert.AnError)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when fetching all teams with condition and relation but row affected is 0, it should return error not found",
			args: args{
				ctx:       context.TODO(),
				condition: "team_name = ?",
				args:      []interface{}{"AcRtf"},
				relations: []string{"Members"},
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				teamRows := mock.NewRows(
					[]string{
						"id",
						"leader_id",
						"competition_id",
						"team_name",
						"join_token",
						"payment_proof_link",
						"twibbon_proof_link",
						"statement_letter_link",
					},
				)

				mock.ExpectQuery(`SELECT \* FROM "teams" WHERE (.+)`).
					WillReturnRows(teamRows)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			if test.beforeTest != nil {
				test.beforeTest(mock)
			}

			repo := repository.NewTeamRepository(db)

			teams, _, err := repo.FetchAllByConditionAndRelation(test.args.ctx, test.args.condition, test.args.args, "created_at DESC", nil, test.args.relations...)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, test.expectedErr, err, "Expecting same error message")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}

			assert.Equal(t, test.want, teams, "Expecting teams result to be equal")
		})
	}
}

func TestCount(t *testing.T) {
	type args struct {
		condition string
		args      []interface{}
		groupBy   string
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(mock sqlmock.Sqlmock)
		want       int64
		wantErr    bool
	}{
		{
			name: "when counting all teams, it should return the count",
			args: args{
				condition: "1 = 1",
				args:      []interface{}{},
				groupBy:   "",
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{"count"}).AddRow(int64(3))

				mock.ExpectQuery(`SELECT count\(\*\) FROM "teams" WHERE 1 = 1`).
					WillReturnRows(rows)
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "when counting all teams but there's an error, it should return internal server error",
			args: args{
				condition: "1 = 1",
				args:      []interface{}{},
				groupBy:   "",
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "teams" WHERE 1 = 1`).
					WillReturnError(assert.AnError)
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, close, mock := pgsqlMock.NewMockDB(t)

			defer close()

			if test.beforeTest != nil {
				test.beforeTest(mock)
			}

			repo := repository.NewTeamRepository(db)

			count, err := repo.Count(test.args.condition, test.args.args, test.args.groupBy)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}

			assert.Equal(t, test.want, count, "Expecting count result to be equal")
		})
	}
}

package test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/internal/app/user/repository"
	dbMock "github.com/hology8/hology-be/internal/infra/database/mock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	type args struct {
		user entity.User
	}

	uuid := uuid.New()

	mockArgs := args{
		user: entity.User{
			ID:              uuid,
			Email:           "test@gmail.com",
			Password:        "randomstring",
			EmailIsVerified: true,
		},
	}

	test := []struct {
		name        string
		args        args
		beforeTest  func(sqlmock.Sqlmock)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when success register user, it should success to register user without error",
			args: mockArgs,
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				mocksql.ExpectBegin()
				mocksql.ExpectQuery("INSERT INTO \"users\" (.+) VALUES (.+)").
					WithArgs(
						mockArgs.user.ID,
						mockArgs.user.Email,
						mockArgs.user.Password,
						mockArgs.user.Fullname,
						mockArgs.user.BirthDate,
						mockArgs.user.EmailVerifiedToken,
						mockArgs.user.ForgotPasswordToken,
						mockArgs.user.EmailIsVerified,
						mockArgs.user.KtmImageLink,
						mockArgs.user.FollowProofLink,
						mockArgs.user.ShareProofLink,
						mockArgs.user.City,
						mockArgs.user.HomeAddress,
						mockArgs.user.ExpiredToken,
						mockArgs.user.ExpiredTokenForgot,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockArgs.user.ID)).
					WillReturnError(nil)
				mocksql.ExpectCommit()
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "When register user, it should return internal server error",
			args: mockArgs,
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				mocksql.ExpectBegin()
				mocksql.ExpectQuery(`INSERT INTO "users" \("id","email","password","fullname","birth_date","email_verified_token","forgot_password_token","email_is_verified","ktm_image_link","follow_proof_link","share_proof_link","city","home_address","expired_token"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9,\$10,\$11,\$12,\$13,\$14\) RETURNING .*`).
					WithArgs(
						mockArgs.user.ID,
						mockArgs.user.Email,
						mockArgs.user.Password,
						mockArgs.user.Fullname,
						mockArgs.user.BirthDate,
						mockArgs.user.EmailVerifiedToken,
						mockArgs.user.ForgotPasswordToken,
						mockArgs.user.EmailIsVerified,
						mockArgs.user.KtmImageLink,
						mockArgs.user.FollowProofLink,
						mockArgs.user.ShareProofLink,
						mockArgs.user.City,
						mockArgs.user.HomeAddress,
						mockArgs.user.ExpiredToken,
						mockArgs.user.ExpiredTokenForgot,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid))
				mocksql.ExpectRollback()
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when register user, it should return error duplicate entry",
			args: mockArgs,
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				mocksql.ExpectBegin()
				mocksql.ExpectQuery(`INSERT INTO "users" \("id","email","password","fullname","birth_date","email_verified_token","forgot_password_token","email_is_verified","ktm_image_link","follow_proof_link","share_proof_link","city","home_address","expired_token","expired_token_forgot"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9,\$10,\$11,\$12,\$13,\$14,\$15\) RETURNING .*`).
					WithArgs(
						mockArgs.user.ID,
						mockArgs.user.Email,
						mockArgs.user.Password,
						mockArgs.user.Fullname,
						mockArgs.user.BirthDate,
						mockArgs.user.EmailVerifiedToken,
						mockArgs.user.ForgotPasswordToken,
						mockArgs.user.EmailIsVerified,
						mockArgs.user.KtmImageLink,
						mockArgs.user.FollowProofLink,
						mockArgs.user.ShareProofLink,
						mockArgs.user.City,
						mockArgs.user.HomeAddress,
						mockArgs.user.ExpiredToken,
						mockArgs.user.ExpiredTokenForgot,
					).
					WillReturnError(gorm.ErrDuplicatedKey)
				mocksql.ExpectRollback()
			},
			wantErr:     true,
			expectedErr: domain.ErrDuplicateEntry,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			gormDb, close, mock := dbMock.NewMockDB(t)
			defer close()

			u := repository.NewUserRepository(gormDb)

			if tt.beforeTest != nil {
				tt.beforeTest(mock)
			}

			err := u.CreateUser(&tt.args.user)

			if tt.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, tt.expectedErr, err, "Expecting error to be %v", tt.expectedErr)
			} else {
				assert.Nil(t, err, "Expecting no error")
			}
		})
	}
}

func TestFindUser(t *testing.T) {
	type args struct {
		user      entity.User
		userParam dto.UserParam
	}

	uuid := uuid.New()

	mockArgs := args{
		user: entity.User{
			ID:       uuid,
			Email:    "test@email.com",
			Password: "randomstring",
			Fullname: "testFullname",
		},
		userParam: dto.UserParam{
			ID:    uuid,
			Email: "test@email.com",
		},
	}

	tests := []struct {
		name        string
		args        args
		beforeTest  func(sqlmock.Sqlmock)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when find user, it should success to find user without error",
			args: mockArgs,
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				mocksql.ExpectQuery(`SELECT \* FROM "users" WHERE
				\("users"\."id" = \$1 AND "users"\."email" = \$2\) AND "users"\."id" = \$3 ORDER BY "users"\."id" LIMIT \$4`).
					WithArgs(mockArgs.userParam.ID, mockArgs.userParam.Email, mockArgs.userParam.ID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "fullname"}).
						AddRow(mockArgs.user.ID, mockArgs.user.Email, mockArgs.user.Password, mockArgs.user.Fullname))
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "when find user, it should return internal server error",
			args: mockArgs,
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				mocksql.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."email" = \$2 ORDER BY "users"\."id" LIMIT 1`).
					WithArgs(mockArgs.userParam.ID, mockArgs.userParam.Email).
					WillReturnError(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when find user, it should return not found error",
			args: mockArgs,
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				mocksql.ExpectQuery(`SELECT \* FROM "users" WHERE \("users"\."id" = \$1 AND "users"\."email" = \$2\) AND "users"\."id" = \$3 ORDER BY "users"\."id" LIMIT \$4`).
					WithArgs(mockArgs.userParam.ID, mockArgs.userParam.Email, mockArgs.userParam.ID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr:     true,
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDb, close, mock := dbMock.NewMockDB(t)
			defer close()

			u := repository.NewUserRepository(gormDb)

			if tt.beforeTest != nil {
				tt.beforeTest(mock)
			}

			err := u.FindUser(&tt.args.user, &tt.args.userParam)

			if tt.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, tt.expectedErr, err, "Expecting error to be %v", tt.expectedErr)
			} else {
				assert.Nil(t, err, "Expecting no error")
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {

	type args struct {
		userUpdate dto.UserUpdate
		userId     uuid.UUID
	}

	uuid := uuid.New()
	expiredTime := time.Now()

	mockArgs := args{
		userUpdate: dto.UserUpdate{
			Password:           "testNewPassword",
			Fullname:           "testNewFullname",
			BirthDate:          "01-01-2000",
			WANumber:           "08123456789",
			LineID:             "testLineID",
			DiscordID:          "testDiscordID",
			StudentID:          "testStudentID",
			City:               "testCity",
			UniversityID:       1,
			KtmImageLink:       "testKtmImageLink",
			EmailIsVerified:    true,
			EmailVerifiedToken: "testEmailVerifiedToken",
			ExpiredToken:       expiredTime,
		},
		userId: uuid,
	}

	test := []struct {
		name        string
		args        args
		beforeTest  func(sqlmock.Sqlmock)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when update user, it should success to update user without error",
			args: mockArgs,
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				mocksql.ExpectBegin()
				mocksql.ExpectExec(`UPDATE "users" SET 
				"password"=\$1,"fullname"=\$2,"birth_date"=\$3,"wa_number"=\$4,"line_id"=\$5,"discord_id"=\$6,"student_id"=\$7,"email_verified_token"=\$8,"email_is_verified"=\$9,"ktm_image_link"=\$10,"city"=\$11,"university_id"=\$12,"expired_token"=\$13
				WHERE id = \$14`).
					WithArgs(
						mockArgs.userUpdate.Password,
						mockArgs.userUpdate.Fullname,
						mockArgs.userUpdate.BirthDate,
						mockArgs.userUpdate.WANumber,
						mockArgs.userUpdate.LineID,
						mockArgs.userUpdate.DiscordID,
						mockArgs.userUpdate.StudentID,
						mockArgs.userUpdate.EmailVerifiedToken,
						mockArgs.userUpdate.EmailIsVerified,
						mockArgs.userUpdate.KtmImageLink,
						mockArgs.userUpdate.City,
						mockArgs.userUpdate.UniversityID,
						mockArgs.userUpdate.ExpiredToken,
						mockArgs.userId,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mocksql.ExpectCommit()
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "when update user, it should return internal server error",
			args: mockArgs,
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				mocksql.ExpectBegin()
				mocksql.ExpectExec(`UPDATE "users" SET 
				"password"=\$1,"fullname"=\$2,"birth_date"=\$3,"wa_number"=\$4,"line_id"=\$5,"discord_id"=\$6,"student_id"=\$7,"email_verified_token"=\$8,"email_is_verified"=\$9,"ktm_image_link"=\$10,"city"=\$11,"university_id"=\$12,"expired_token"=\$13
				WHERE id = \$14`).
					WithArgs(
						mockArgs.userUpdate.Password,
						mockArgs.userUpdate.Fullname,
						mockArgs.userUpdate.BirthDate,
						mockArgs.userUpdate.WANumber,
						mockArgs.userUpdate.LineID,
						mockArgs.userUpdate.DiscordID,
						mockArgs.userUpdate.StudentID,
						mockArgs.userUpdate.EmailVerifiedToken,
						mockArgs.userUpdate.EmailIsVerified,
						mockArgs.userUpdate.KtmImageLink,
						mockArgs.userUpdate.City,
						mockArgs.userUpdate.UniversityID,
						mockArgs.userUpdate.ExpiredToken,
						mockArgs.userId,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mocksql.ExpectRollback()
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			gormDb, close, mock := dbMock.NewMockDB(t)
			defer close()

			u := repository.NewUserRepository(gormDb)

			if tt.beforeTest != nil {
				tt.beforeTest(mock)
			}

			err := u.UpdateUser(&tt.args.userUpdate, tt.args.userId)

			if tt.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, tt.expectedErr, err, "Expecting error to be %v", tt.expectedErr)
			} else {
				assert.Nil(t, err, "Expecting no error")
			}
		})
	}
}

func TestFetchAllByConditionAndRelation(t *testing.T) {
	type args struct {
		condition string
		args      []interface{}
		relation  []string
	}

	dummyUUID := uuid.New()

	tests := []struct {
		name        string
		args        args
		beforeTest  func(sqlmock.Sqlmock)
		want        []entity.User
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when fetch all without relation and condition, it should success to fetch all user",
			args: args{
				condition: "",
				args:      nil,
				relation:  nil,
			},
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				userRows := mocksql.NewRows(
					[]string{
						"id",
						"email",
						"password",
						"fullname",
					}).
					AddRow(
						dummyUUID,
						"test1@email.com",
						"randomstring",
						"test1Fullname",
					).
					AddRow(
						dummyUUID,
						"test2@email.com",
						"randomstring",
						"test2Fullname",
					)

				mocksql.ExpectQuery(`SELECT \* FROM "users"`).
					WillReturnRows(userRows)
			},
			want: []entity.User{
				{
					ID:       dummyUUID,
					Email:    "test1@email.com",
					Password: "randomstring",
					Fullname: "test1Fullname",
				},
				{
					ID:       dummyUUID,
					Email:    "test2@email.com",
					Password: "randomstring",
					Fullname: "test2Fullname",
				},
			},
			wantErr: false,
		},
		{
			name: "when fetch all with relation and condition, it should success to fetch all user",
			args: args{
				condition: "id = ?",
				args:      []interface{}{dummyUUID},
				relation:  []string{"Teams"},
			},
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				teamRows := mocksql.NewRows(
					[]string{
						"id",
						"email",
						"password",
						"fullname",
					}).
					AddRow(
						dummyUUID,
						"test1@email.com",
						"randomstring",
						"test1Fullname",
					)

				detailTeamsRows := mocksql.NewRows(
					[]string{
						"user_id",
						"team_id",
					}).
					AddRow(
						dummyUUID,
						"1",
					).
					AddRow(
						dummyUUID,
						"2",
					)

				mocksql.ExpectQuery(`SELECT \* FROM "users" WHERE (.+)`).
					WithArgs(dummyUUID).
					WillReturnRows(teamRows)

				mocksql.ExpectQuery(`SELECT \* FROM "detail_teams" WHERE (.+)`).
					WithArgs(dummyUUID).
					WillReturnRows(detailTeamsRows)
			},
			want: []entity.User{
				{
					ID:       dummyUUID,
					Email:    "test1@email.com",
					Password: "randomstring",
					Fullname: "test1Fullname",
					Teams: []entity.DetailTeams{
						{
							UserID: dummyUUID,
							TeamID: "1",
						},
						{
							UserID: dummyUUID,
							TeamID: "2",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "when fetch all with relation and condition but there's an error, it should return internal server error",
			args: args{
				condition: "id = ?",
				args:      []interface{}{dummyUUID},
				relation:  []string{"Teams"},
			},
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				mocksql.ExpectQuery(`SELECT \* FROM "users" WHERE (.+)`).
					WithArgs(dummyUUID).
					WillReturnError(domain.ErrInternalServer)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when fetch all with relation and condition but it return 0 data, it should return not found error",
			args: args{
				condition: "id = ?",
				args:      []interface{}{dummyUUID},
				relation:  []string{"Teams"},
			},
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				userRows := mocksql.NewRows(
					[]string{
						"id",
						"email",
						"password",
						"fullname",
					})

				mocksql.ExpectQuery(`SELECT \* FROM "users" WHERE (.+)`).
					WithArgs(dummyUUID).
					WillReturnRows(userRows)
			},
			want:        nil,
			wantErr:     true,
			expectedErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDb, close, mock := dbMock.NewMockDB(t)
			defer close()

			u := repository.NewUserRepository(gormDb)

			if tt.beforeTest != nil {
				tt.beforeTest(mock)
			}

			users, _, err := u.FetchAllByConditionAndRelation(tt.args.condition, tt.args.args, nil, nil, tt.args.relation...)

			if tt.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, tt.expectedErr, err, "Expecting error to be %v", tt.expectedErr)
			} else {
				assert.Nil(t, err, "Expecting no error")
				assert.Equal(t, tt.want, users, "Expecting result to be %v", tt.want)
			}
		})
	}
}

func TestDeleteUnverifiedUsers(t *testing.T) {
	gormDb, close, mock := dbMock.NewMockDB(t)
	defer close()

	mock.ExpectBegin()

	mock.ExpectExec(`DELETE FROM "users" WHERE email_is_verified = ?`).
		WithArgs(false).WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectCommit()

	repo := repository.NewUserRepository(gormDb)

	err := repo.DeleteUnverifiedUser()

	assert.Nil(t, err, "expecting no error")
}

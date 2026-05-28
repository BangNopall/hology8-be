package test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/internal/app/admin/repository"
	dbMock "github.com/hology8/hology-be/internal/infra/database/mock"
	"github.com/stretchr/testify/assert"
)

func TestFindAdmin(t *testing.T) {
	type args struct {
		admin      entity.Admin
		adminParam dto.AdminParam
	}

	mockArgs := args{
		admin: entity.Admin{
			Username: "testUsername",
			Password: "randomstring",
		},
		adminParam: dto.AdminParam{
			Username: "testUsername",
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
			name: "when find admin, it should success without error",
			args: mockArgs,
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				mocksql.ExpectQuery(`SELECT \* FROM \"admins\" WHERE \"admins\"\.\"username\" = \$1 ORDER BY \"admins\"\.\"id\" LIMIT \$2`).
					WithArgs(mockArgs.adminParam.Username, 1).
					WillReturnRows(sqlmock.NewRows([]string{"username", "password"}).AddRow(mockArgs.admin.Username, mockArgs.admin.Password))
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "when find admin, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				mocksql.ExpectQuery(`SELECT \* FROM \"admins\" WHERE \"admins\"\.\"username\" = \$1 ORDER BY \"admins\"\.\"id\" LIMIT \$2`).
					WithArgs(mockArgs.adminParam.Username, 1).
					WillReturnError(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when find admin, it should failed with error not found",
			args: mockArgs,
			beforeTest: func(mocksql sqlmock.Sqlmock) {
				mocksql.ExpectQuery(`SELECT \* FROM \"admins\" WHERE \"admins\"\.\"username\" = \$1 ORDER BY \"admins\"\.\"id\" LIMIT \$2`).
					WithArgs(mockArgs.adminParam.Username, 1).
					WillReturnError(domain.ErrNotFound)
			},
			wantErr:     true,
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			gormDb, close, mock := dbMock.NewMockDB(t)
			defer close()

			a := repository.NewAdminRepository(gormDb)

			if tt.beforeTest != nil {
				tt.beforeTest(mock)
			}

			err := a.FindAdmin(&tt.args.admin, &tt.args.adminParam)

			if tt.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, tt.expectedErr, err, "Expecting error to be %v", tt.expectedErr)
			} else {
				assert.Nil(t, err, "Expecting no error")
			}
		})
	}
}

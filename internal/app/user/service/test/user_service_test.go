package test

import (
	"context"
	"mime/multipart"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
	mockUserRepo "github.com/BangNopall/hology8-be/internal/app/user/repository/mock"
	"github.com/BangNopall/hology8-be/internal/app/user/service"
	mockAws "github.com/BangNopall/hology8-be/pkg/aws/mock"
	"github.com/BangNopall/hology8-be/pkg/bcrypt"
	mockBcrypt "github.com/BangNopall/hology8-be/pkg/bcrypt/mock"
	mockGoMail "github.com/BangNopall/hology8-be/pkg/gomail/mock"
	mockJwt "github.com/BangNopall/hology8-be/pkg/jwt/mock"
	mockRedis "github.com/BangNopall/hology8-be/pkg/redis/mock"
	mockTime "github.com/BangNopall/hology8-be/pkg/time/mock"
	mockUUID "github.com/BangNopall/hology8-be/pkg/uuid/mock"
)

type mockObjects struct {
	userRepo *mockUserRepo.MockUserRepository
	jwt      *mockJwt.MockJwtInterface
	uuid     *mockUUID.MockUUIDInterface
	bcrypt   *mockBcrypt.MockBcryptInterface
	time     *mockTime.MockTimeInterface
	goMail   *mockGoMail.MockGoMailInterface
	storage  *mockAws.MockCloudStorage
	redis    *mockRedis.MockRedisInterface
}

func TestLoginRegisterOauth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mockUserRepo.NewMockUserRepository(ctrl)
	jwt := mockJwt.NewMockJwtInterface(ctrl)
	uuidMock := mockUUID.NewMockUUIDInterface(ctrl)

	mockObject := mockObjects{
		userRepo: r,
		jwt:      jwt,
		uuid:     uuidMock,
	}

	type args struct {
		userOauth      dto.UserOauth
		unverifiedUser entity.User
		verifiedUser   entity.User
		loginResponse  dto.UserLoginResponse
		userParam      dto.UserParam
	}

	uuid := uuid.New()

	mockArgs := args{
		dto.UserOauth{
			Email:      "test@email.com",
			Name:       "testName",
			GivenName:  "testGivenName",
			FamilyName: "testFamilyName",
			Id:         "testId",
			Picture:    "testPicture",
		},
		entity.User{
			ID:       uuid,
			Email:    "test@email.com",
			Password: "randomstring",
		},
		entity.User{
			ID:              uuid,
			Email:           "test@email.com",
			EmailIsVerified: true,
		},
		dto.UserLoginResponse{
			Token: "testtoken",
		},
		dto.UserParam{
			ID: uuid,
		},
	}

	tests := []struct {
		name        string
		args        args
		beforeTest  func(mockObject mockObjects)
		want        dto.UserLoginResponse
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when register user, it should success without error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.jwt.EXPECT().GenerateToken(uuid, gomock.Any(), gomock.Any()).Return(mockArgs.loginResponse.Token, nil)
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userOauth.Email}).Return(domain.ErrNotFound)
				mockObject.userRepo.EXPECT().CreateUser(&entity.User{ID: uuid, Email: mockArgs.userOauth.Email, EmailIsVerified: true}).Return(nil)
				mockObject.uuid.EXPECT().New().Return(uuid, nil)
			},
			want:        mockArgs.loginResponse,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "when register user, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userOauth.Email}).Return(domain.ErrNotFound)
				mockObject.uuid.EXPECT().New().Return(uuid, nil)
				mockObject.userRepo.EXPECT().CreateUser(&entity.User{ID: uuid, Email: mockArgs.userOauth.Email, EmailIsVerified: true}).Return(nil)
				mockObject.jwt.EXPECT().GenerateToken(uuid, gomock.Any(), gomock.Any()).Return("", domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when register user, it should error when generate uuid with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userOauth.Email}).Return(domain.ErrNotFound)
				mockObject.uuid.EXPECT().New().Return(uuid, domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when register user, it should failed when create user with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userOauth.Email}).Return(domain.ErrNotFound)
				mockObject.uuid.EXPECT().New().Return(uuid, nil)
				mockObject.userRepo.EXPECT().CreateUser(&entity.User{ID: uuid, Email: mockArgs.userOauth.Email, EmailIsVerified: true}).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when login user with oauth, it should success without error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.jwt.EXPECT().GenerateToken(uuid, gomock.Any(), gomock.Any()).Return(mockArgs.loginResponse.Token, nil)
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userOauth.Email}).SetArg(0, mockArgs.verifiedUser).Return(nil)
			},
			want:        mockArgs.loginResponse,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "when login user with oauth, it should failed when generate token with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userOauth.Email}).SetArg(0, mockArgs.verifiedUser).Return(nil)
				mockObject.jwt.EXPECT().GenerateToken(uuid, gomock.Any(), gomock.Any()).Return("", domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when login user with oauth, it should verified the unverified user without error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.jwt.EXPECT().GenerateToken(uuid, gomock.Any(), gomock.Any()).Return(mockArgs.loginResponse.Token, nil)
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userOauth.Email}).SetArg(0, mockArgs.unverifiedUser).Return(nil)
				mockObject.userRepo.EXPECT().UpdateUser(&dto.UserUpdate{EmailIsVerified: true}, mockArgs.unverifiedUser.ID).Return(nil)
			},
			want:        mockArgs.loginResponse,
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "when login user with oauth, it should failed when update user with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userOauth.Email}).SetArg(0, mockArgs.unverifiedUser).Return(nil)
				mockObject.userRepo.EXPECT().UpdateUser(&dto.UserUpdate{EmailIsVerified: true}, mockArgs.unverifiedUser.ID).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when login user with oauth, it should failed when generate token with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userOauth.Email}).SetArg(0, mockArgs.unverifiedUser).Return(nil)
				mockObject.userRepo.EXPECT().UpdateUser(&dto.UserUpdate{EmailIsVerified: true}, mockArgs.unverifiedUser.ID).Return(nil)
				mockObject.jwt.EXPECT().GenerateToken(uuid, gomock.Any(), gomock.Any()).Return("", domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when login user with oauth, it should failed when find user with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userOauth.Email}).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when login user with oauth, it should failed with time out",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userOauth.Email}).SetArg(0, mockArgs.verifiedUser).Return(nil)
				mockObject.jwt.EXPECT().GenerateToken(uuid, gomock.Any(), gomock.Any()).
					DoAndReturn(func(any1 any, any2 any, any3 any) (string, error) {
						time.Sleep(time.Millisecond * 1000)
						return mockArgs.loginResponse.Token, nil
					})
			},
			wantErr:     true,
			expectedErr: domain.ErrTimeout,
		},
		{
			name: "when register user, it should failed with time out",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userOauth.Email}).Return(domain.ErrNotFound)
				mockObject.userRepo.EXPECT().CreateUser(&entity.User{ID: uuid, Email: mockArgs.userOauth.Email, EmailIsVerified: true}).Return(nil)
				mockObject.uuid.EXPECT().New().Return(uuid, nil)
				mockObject.jwt.EXPECT().GenerateToken(uuid, gomock.Any(), gomock.Any()).
					DoAndReturn(func(any1 any, any2 any, any3 any) (string, error) {
						time.Sleep(time.Millisecond * 1000)
						return mockArgs.loginResponse.Token, nil
					})
			},
			wantErr:     true,
			expectedErr: domain.ErrTimeout,
		},
		{
			name: "when login and verifying user, it should failed with time out",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userOauth.Email}).SetArg(0, mockArgs.unverifiedUser).Return(nil)
				mockObject.userRepo.EXPECT().UpdateUser(&dto.UserUpdate{EmailIsVerified: true}, mockArgs.unverifiedUser.ID).Return(nil)
				mockObject.jwt.EXPECT().GenerateToken(uuid, gomock.Any(), gomock.Any()).
					DoAndReturn(func(any1 any, any2 any, any3 any) (string, error) {
						time.Sleep(time.Millisecond * 1000)
						return mockArgs.loginResponse.Token, nil
					})
			},
			wantErr:     true,
			expectedErr: domain.ErrTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := service.NewUserService(mockObject.userRepo, mockObject.uuid, nil, nil, nil, mockObject.jwt, nil, time.Millisecond*500, nil)

			if tt.beforeTest != nil {
				tt.beforeTest(mockObject)
			}

			got, err := w.LoginRegisterOauth(context.Background(), tt.args.userOauth)

			if tt.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, tt.expectedErr, err, "Expecting error to be %v", tt.expectedErr)
			} else {
				assert.Nil(t, err, "Expecting no error")
				assert.Equal(t, tt.want, got, "Expecting response to be %v", tt.want)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mockUserRepo.NewMockUserRepository(ctrl)
	bcrypt := mockBcrypt.NewMockBcryptInterface(ctrl)
	goMail := mockGoMail.NewMockGoMailInterface(ctrl)
	uuidMock := mockUUID.NewMockUUIDInterface(ctrl)
	timeMock := mockTime.NewMockTimeInterface(ctrl)

	mockObject := mockObjects{
		userRepo: r,
		bcrypt:   bcrypt,
		goMail:   goMail,
		uuid:     uuidMock,
		time:     timeMock,
	}

	type args struct {
		user                dto.UserRegister
		notExpiredTokenUser entity.User
		expiredTokenUser    entity.User
		newUser             entity.User
	}

	uuid := uuid.New()
	currentTime := time.Now()
	expiredToken := time.Now().Add(time.Hour * 1)

	mockArgs := args{
		dto.UserRegister{
			Email:           "testEmail",
			Password:        "testPassword",
			ConfirmPassword: "testPassword",
		},
		entity.User{
			ID:              uuid,
			Email:           "testEmail",
			Password:        "testPassword",
			EmailIsVerified: false,
			ExpiredToken:    expiredToken,
		},
		entity.User{
			ID:              uuid,
			Email:           "testEmail",
			Password:        "testPassword",
			EmailIsVerified: false,
			ExpiredToken:    expiredToken,
		},
		entity.User{
			ID:                 uuid,
			Email:              "testEmail",
			Password:           "successHashingPassword",
			EmailIsVerified:    false,
			EmailVerifiedToken: "successHashingVerPassword",
			ExpiredToken:       expiredToken,
		},
	}

	tests := []struct {
		name        string
		args        args
		beforeTest  func(mockObject mockObjects)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when register user, it should success without error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return("successHashingPassword", nil)
				mockObject.bcrypt.EXPECT().Hash(gomock.Any()).Return("successHashingVerPassword", nil)
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).Return(domain.ErrNotFound)
				mockObject.uuid.EXPECT().New().Return(uuid, nil)
				mockObject.userRepo.EXPECT().CreateUser(&mockArgs.newUser).Return(nil)
				mockObject.goMail.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
			},
			wantErr: false,
		},
		{
			name: "when register user's password doesnt match with confirm password, it should return error password not match",
			args: args{
				user: dto.UserRegister{
					Email:           "testEmail",
					Password:        "testPassword",
					ConfirmPassword: "testPassword1",
				},
			},
			beforeTest:  func(mockObject mockObjects) {},
			wantErr:     true,
			expectedErr: domain.ErrConfirmPasswordNotMatch,
		},
		{
			name: "when register user and failed hashing password, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return("", domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when register user but failed hashing email verification password, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return("successHashingPassword", nil)
				mockObject.bcrypt.EXPECT().Hash(gomock.Any()).Return("", domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when register user but failed to find user, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return("successHashingPassword", nil)
				mockObject.bcrypt.EXPECT().Hash(gomock.Any()).Return("successHashingVerPassword", nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when register user but the token is not expired, it should failed with error please check email",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return("successHashingPassword", nil)
				mockObject.bcrypt.EXPECT().Hash(gomock.Any()).Return("successHashingVerPassword", nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).SetArg(0, mockArgs.notExpiredTokenUser).Return(nil)
			},
			wantErr:     true,
			expectedErr: domain.ErrCheckEmail,
		},
		{
			name: "when register user but failed to update user, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return("successHashingPassword", nil)
				mockObject.bcrypt.EXPECT().Hash(gomock.Any()).Return("successHashingVerPassword", nil)
				mockObject.time.EXPECT().Now().Return(currentTime.Add(time.Hour * 2))
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).SetArg(0, mockArgs.expiredTokenUser).Return(nil)
				mockObject.goMail.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockObject.userRepo.EXPECT().UpdateUser(&dto.UserUpdate{Password: "successHashingPassword", EmailVerifiedToken: "successHashingVerPassword", ExpiredToken: expiredToken}, mockArgs.expiredTokenUser.ID).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when register user but failed to send email on expired token, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return("successHashingPassword", nil)
				mockObject.bcrypt.EXPECT().Hash(gomock.Any()).Return("successHashingVerPassword", nil)
				mockObject.time.EXPECT().Now().Return(currentTime.Add(time.Hour * 2))
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).SetArg(0, mockArgs.expiredTokenUser).Return(nil)
				mockObject.goMail.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when success register user on expired token, it should success without error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return("successHashingPassword", nil)
				mockObject.bcrypt.EXPECT().Hash(gomock.Any()).Return("successHashingVerPassword", nil)
				mockObject.time.EXPECT().Now().Return(currentTime.Add(time.Hour * 2))
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).SetArg(0, mockArgs.expiredTokenUser).Return(nil)
				mockObject.goMail.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockObject.userRepo.EXPECT().UpdateUser(&dto.UserUpdate{Password: "successHashingPassword", EmailVerifiedToken: "successHashingVerPassword", ExpiredToken: expiredToken}, mockArgs.expiredTokenUser.ID).Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "when register user but failed to generate uuid, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return("successHashingPassword", nil)
				mockObject.bcrypt.EXPECT().Hash(gomock.Any()).Return("successHashingVerPassword", nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).Return(domain.ErrNotFound)
				mockObject.uuid.EXPECT().New().Return(uuid, domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when register user but failed to create user, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return("successHashingPassword", nil)
				mockObject.bcrypt.EXPECT().Hash(gomock.Any()).Return("successHashingVerPassword", nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).Return(domain.ErrNotFound)
				mockObject.uuid.EXPECT().New().Return(uuid, nil)
				mockObject.goMail.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockObject.userRepo.EXPECT().CreateUser(&mockArgs.newUser).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when register user but failed to send email for new user, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return("successHashingPassword", nil)
				mockObject.bcrypt.EXPECT().Hash(gomock.Any()).Return("successHashingVerPassword", nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).Return(domain.ErrNotFound)
				mockObject.uuid.EXPECT().New().Return(uuid, nil)
				mockObject.goMail.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when register user but operation time exceeded, it should failed with error time out",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return("successHashingPassword", nil)
				mockObject.bcrypt.EXPECT().Hash(gomock.Any()).Return("successHashingVerPassword", nil)
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).Return(domain.ErrNotFound)
				mockObject.uuid.EXPECT().New().Return(uuid, nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
				mockObject.goMail.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(any1 any, any2 any, any3 any) error {
						time.Sleep(time.Millisecond * 1500)
						return nil
					})
			},
			wantErr:     true,
			expectedErr: domain.ErrTimeout,
		},
		{
			name: "when resend email to user but operation time exceeded, it should failed with error time out",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return("successHashingPassword", nil)
				mockObject.bcrypt.EXPECT().Hash(gomock.Any()).Return("successHashingVerPassword", nil)
				mockObject.time.EXPECT().Now().Return(currentTime.Add(time.Hour * 2))
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).SetArg(0, mockArgs.expiredTokenUser).Return(nil)
				mockObject.goMail.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(any1 any, any2 any, any3 any) error {
						time.Sleep(time.Millisecond * 1500)
						return nil
					})
			},
			wantErr:     true,
			expectedErr: domain.ErrTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := service.NewUserService(mockObject.userRepo, mockObject.uuid, mockObject.bcrypt, mockObject.time, mockObject.goMail, nil, nil, time.Millisecond*1000, nil)

			if tt.beforeTest != nil {
				tt.beforeTest(mockObject)
			}

			err := w.Register(context.Background(), tt.args.user, "https://example.com/")

			if tt.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, tt.expectedErr, err, "Expecting error to be %v", tt.expectedErr)
			} else {
				assert.Nil(t, err, "Expecting no error")
			}
		})
	}
}

func TestVerifyEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mockUserRepo.NewMockUserRepository(ctrl)
	bcrypt := mockBcrypt.NewMockBcryptInterface(ctrl)
	timeMock := mockTime.NewMockTimeInterface(ctrl)

	mockObject := mockObjects{
		userRepo: r,
		bcrypt:   bcrypt,
		time:     timeMock,
	}

	type args struct {
		email        string
		emailVerPass string
		user         entity.User
	}

	uuid := uuid.New()
	currentTime := time.Now()
	expiredToken := time.Now().Add(time.Hour * 1)

	mockArgs := args{
		"testEmail",
		"testEmailVerPass",
		entity.User{
			ID:                 uuid,
			Email:              "testEmail",
			Password:           "testPassword",
			EmailIsVerified:    false,
			EmailVerifiedToken: "testHashedEmailVerPass",
			ExpiredToken:       expiredToken,
		},
	}

	tests := []struct {
		name        string
		args        args
		beforeTest  func(mockObject mockObjects)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when verify email, it should success without error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.email}).SetArg(0, mockArgs.user).Return(nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.bcrypt.EXPECT().Compare(mockArgs.emailVerPass, mockArgs.user.EmailVerifiedToken).Return(true)
				mockObject.userRepo.EXPECT().UpdateUser(&dto.UserUpdate{EmailIsVerified: true}, mockArgs.user.ID).Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "when verify email but failed to find user, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.email}).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when verify email but failed to compare email verification password, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.email}).SetArg(0, mockArgs.user).Return(nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.bcrypt.EXPECT().Compare(mockArgs.emailVerPass, mockArgs.user.EmailVerifiedToken).Return(false)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when verify email but the token expired, it should failed with error invalid token",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.email}).SetArg(0, mockArgs.user).Return(nil)
				mockObject.time.EXPECT().Now().Return(currentTime.Add(time.Hour * 2))
			},
			wantErr:     true,
			expectedErr: domain.ErrInvalidToken,
		},
		{
			name: "when verify email but failed to update user, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.email}).SetArg(0, mockArgs.user).Return(nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.bcrypt.EXPECT().Compare(mockArgs.emailVerPass, mockArgs.user.EmailVerifiedToken).Return(true)
				mockObject.userRepo.EXPECT().UpdateUser(&dto.UserUpdate{EmailIsVerified: true}, mockArgs.user.ID).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when verify email but operation time exceeded, it should failed with error time out",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.email}).SetArg(0, mockArgs.user).Return(nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.bcrypt.EXPECT().Compare(mockArgs.emailVerPass, mockArgs.user.EmailVerifiedToken).Return(true)
				mockObject.userRepo.EXPECT().UpdateUser(&dto.UserUpdate{EmailIsVerified: true}, mockArgs.user.ID).
					DoAndReturn(func(any any, any2 any) error {
						time.Sleep(time.Millisecond * 1000)
						return nil
					})
			},
			wantErr:     true,
			expectedErr: domain.ErrTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := service.NewUserService(mockObject.userRepo, nil, mockObject.bcrypt, mockObject.time, nil, nil, nil, time.Millisecond*500, nil)

			if tt.beforeTest != nil {
				tt.beforeTest(mockObject)
			}

			err := w.VerifyEmail(context.Background(), tt.args.email, tt.args.emailVerPass)

			if tt.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, tt.expectedErr, err, "Expecting error to be %v", tt.expectedErr)
			} else {
				assert.Nil(t, err, "Expecting no error")
			}
		})
	}
}

func TestLoginWithEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mockUserRepo.NewMockUserRepository(ctrl)
	jwt := mockJwt.NewMockJwtInterface(ctrl)
	bcrypt := mockBcrypt.NewMockBcryptInterface(ctrl)

	mockObject := mockObjects{
		userRepo: r,
		jwt:      jwt,
		bcrypt:   bcrypt,
	}

	type args struct {
		userLogin     dto.UserLogin
		user          entity.User
		loginResponse dto.UserLoginResponse
	}

	uuid := uuid.New()

	mockArgs := args{
		dto.UserLogin{
			Email:    "testEmail",
			Password: "testPassword",
		},
		entity.User{
			ID:              uuid,
			Email:           "testEmail",
			Password:        "testHashedPassword",
			EmailIsVerified: true,
		},
		dto.UserLoginResponse{
			Token: "testToken",
		},
	}

	tests := []struct {
		name        string
		args        args
		beforeTest  func(mockObject mockObjects)
		want        dto.UserLoginResponse
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when login with email, it should success without error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userLogin.Email}).SetArg(0, mockArgs.user).Return(nil)
				mockObject.bcrypt.EXPECT().Compare(mockArgs.userLogin.Password, mockArgs.user.Password).Return(true)
				mockObject.jwt.EXPECT().GenerateToken(uuid, gomock.Any(), gomock.Any()).Return(mockArgs.loginResponse.Token, nil)
			},
			want:    mockArgs.loginResponse,
			wantErr: false,
		},
		{
			name: "when login with email but failed to find user, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userLogin.Email}).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when login with email but failed to find user, it should failed with wrong password/username error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userLogin.Email}).Return(domain.ErrNotFound)
			},
			wantErr:     true,
			expectedErr: domain.ErrWrongEmailOrPassword,
		},
		{
			name: "when login with email but failed to compare password, it should failed with wrong password/username error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userLogin.Email}).SetArg(0, mockArgs.user).Return(nil)
				mockObject.bcrypt.EXPECT().Compare(mockArgs.userLogin.Password, mockArgs.user.Password).Return(false)
			},
			wantErr:     true,
			expectedErr: domain.ErrWrongEmailOrPassword,
		},
		{
			name: "when login with email but failed to generate token, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userLogin.Email}).SetArg(0, mockArgs.user).Return(nil)
				mockObject.bcrypt.EXPECT().Compare(mockArgs.userLogin.Password, mockArgs.user.Password).Return(true)
				mockObject.jwt.EXPECT().GenerateToken(uuid, gomock.Any(), gomock.Any()).Return("", domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when login with email but operation time exceeded, it should failed with error time out",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.userLogin.Email}).SetArg(0, mockArgs.user).Return(nil)
				mockObject.bcrypt.EXPECT().Compare(mockArgs.userLogin.Password, mockArgs.user.Password).Return(true)
				mockObject.jwt.EXPECT().GenerateToken(uuid, gomock.Any(), gomock.Any()).
					DoAndReturn(func(any any, any2 any, any3 any) (string, error) {
						time.Sleep(time.Millisecond * 1000)
						return mockArgs.loginResponse.Token, nil
					})
			},
			wantErr:     true,
			expectedErr: domain.ErrTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := service.NewUserService(mockObject.userRepo, nil, mockObject.bcrypt, nil, nil, mockObject.jwt, nil, time.Millisecond*500, nil)

			if tt.beforeTest != nil {
				tt.beforeTest(mockObject)
			}

			got, err := w.LoginWithEmail(context.Background(), tt.args.userLogin)

			if tt.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, tt.expectedErr, err, "Expecting error to be %v", tt.expectedErr)
			} else {
				assert.Nil(t, err, "Expecting no error")
				assert.Equal(t, tt.want, got, "Expecting response to be %v", tt.want)
			}
		})
	}
}

func TestResetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mockUserRepo.NewMockUserRepository(ctrl)
	bcryptMock := mockBcrypt.NewMockBcryptInterface(ctrl)
	timeMock := mockTime.NewMockTimeInterface(ctrl)

	mockObject := mockObjects{
		userRepo: r,
		bcrypt:   bcryptMock,
		time:     timeMock,
	}

	type args struct {
		user                dto.UserResetPassword
		forgotPasswordToken string
	}

	uuid := uuid.New()
	currentTime := time.Now()
	expiredToken := time.Now().Add(time.Hour * 1)
	hashedPassword, _ := bcrypt.Bcrypt.Hash("dummyPassword")
	mockArgs := args{
		dto.UserResetPassword{
			Password:        "dummyPassword",
			ConfirmPassword: "dummyPassword",
		},
		"dummyForgotPasswordToken",
	}

	tests := []struct {
		name        string
		args        args
		beforeTest  func(mockObject mockObjects)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when reset password, it should success without error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{ForgotPasswordToken: mockArgs.forgotPasswordToken}).
					SetArg(0, entity.User{ID: uuid, ExpiredTokenForgot: expiredToken}).Return(nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return(hashedPassword, nil)
				mockObject.userRepo.EXPECT().UpdateUser(gomock.Any(), uuid).Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "when register user's password doesnt match with confirm password, it should return error password not match",
			args: args{
				user: dto.UserResetPassword{
					Password:        "testPassword",
					ConfirmPassword: "testPassword1",
				},
			},
			beforeTest:  func(mockObject mockObjects) {},
			wantErr:     true,
			expectedErr: domain.ErrConfirmPasswordNotMatch,
		},
		{
			name: "when reset password but there is no user found, it should failed with error invalid token",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{ForgotPasswordToken: mockArgs.forgotPasswordToken}).Return(domain.ErrNotFound)
			},
			wantErr:     true,
			expectedErr: domain.ErrInvalidToken,
		},
		{
			name: "when reset password but the token expired, it should failed with error invalid token",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{ForgotPasswordToken: mockArgs.forgotPasswordToken}).SetArg(0, entity.User{}).Return(nil)
				mockObject.time.EXPECT().Now().Return(expiredToken)
			},
			wantErr:     true,
			expectedErr: domain.ErrInvalidToken,
		},
		{
			name: "when reset password but failed to update user, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{ForgotPasswordToken: mockArgs.forgotPasswordToken}).
					SetArg(0, entity.User{ID: uuid, ExpiredTokenForgot: expiredToken}).Return(nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return(hashedPassword, nil)
				mockObject.userRepo.EXPECT().UpdateUser(gomock.Any(), uuid).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when reset password but operation time exceeded, it should failed with error time out",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{ForgotPasswordToken: mockArgs.forgotPasswordToken}).
					SetArg(0, entity.User{ID: uuid, ExpiredTokenForgot: expiredToken}).Return(nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.bcrypt.EXPECT().Hash(mockArgs.user.Password).Return(hashedPassword, nil)
				mockObject.userRepo.EXPECT().UpdateUser(gomock.Any(), uuid).
					DoAndReturn(func(any any, any2 any) error {
						time.Sleep(time.Millisecond * 600)
						return nil
					})
			},
			wantErr:     true,
			expectedErr: domain.ErrTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := service.NewUserService(mockObject.userRepo, mockObject.uuid, mockObject.bcrypt, mockObject.time, mockObject.goMail, nil, nil, time.Millisecond*500, nil)

			if tt.beforeTest != nil {
				tt.beforeTest(mockObject)
			}

			err := w.ResetPassword(context.Background(), tt.args.user, tt.args.forgotPasswordToken)

			if tt.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, tt.expectedErr, err, "Expecting error to be %v", tt.expectedErr)
			} else {
				assert.Nil(t, err, "Expecting no error")
			}
		})
	}
}

func TestForgotPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mockUserRepo.NewMockUserRepository(ctrl)
	bcryptMock := mockBcrypt.NewMockBcryptInterface(ctrl)
	timeMock := mockTime.NewMockTimeInterface(ctrl)
	goMail := mockGoMail.NewMockGoMailInterface(ctrl)

	mockObject := mockObjects{
		userRepo: r,
		bcrypt:   bcryptMock,
		time:     timeMock,
		goMail:   goMail,
	}

	type args struct {
		user dto.UserForgotPassword
	}

	uuid := uuid.New()
	currentTime := time.Now()
	expiredToken := time.Now().Add(time.Hour * 1)

	mockArgs := args{
		dto.UserForgotPassword{
			Email: "dummy@gmail.com",
		},
	}

	tests := []struct {
		name        string
		args        args
		beforeTest  func(mockObject mockObjects)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when forgot password, it should success without error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).
					SetArg(0, entity.User{ID: uuid}).Return(nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
				mockObject.userRepo.EXPECT().UpdateUser(gomock.Any(), uuid).Return(nil)
				mockObject.goMail.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "when forgot password but failed to find user, it should failed with error not found error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).Return(domain.ErrNotFound)
			},
			wantErr:     true,
			expectedErr: domain.ErrUserNotFound,
		},
		{
			name: "when forgot password but token is not expired, it should failed with error check email error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).
					SetArg(0, entity.User{ID: uuid, ExpiredTokenForgot: currentTime.Add(time.Hour * 2)}).Return(nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
			},
			wantErr:     true,
			expectedErr: domain.ErrCheckEmail,
		},
		{
			name: "when forgot password but failed to update user, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).
					SetArg(0, entity.User{ID: uuid}).Return(nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
				mockObject.userRepo.EXPECT().UpdateUser(gomock.Any(), uuid).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when forgot password but failed to send email, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{Email: mockArgs.user.Email}).
					SetArg(0, entity.User{ID: uuid}).Return(nil)
				mockObject.time.EXPECT().Now().Return(currentTime)
				mockObject.time.EXPECT().Add(time.Hour * 1).Return(expiredToken)
				mockObject.userRepo.EXPECT().UpdateUser(gomock.Any(), uuid).Return(nil)
				mockObject.goMail.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := service.NewUserService(mockObject.userRepo, mockObject.uuid, mockObject.bcrypt, mockObject.time, mockObject.goMail, nil, nil, time.Millisecond*500, nil)

			if tt.beforeTest != nil {
				tt.beforeTest(mockObject)
			}

			err := w.ForgotPassword(context.Background(), tt.args.user, "https://example.com")

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mockUserRepo.NewMockUserRepository(ctrl)

	mockObject := mockObjects{
		userRepo: r,
	}

	type args struct {
		userUpdate dto.UserUpdate
		userId     uuid.UUID
	}

	uuid := uuid.New()

	mockArgs := args{
		dto.UserUpdate{
			Fullname:  "testNewFullname",
			BirthDate: "20-08-2024",
		},
		uuid,
	}

	tests := []struct {
		name        string
		args        args
		beforeTest  func(mockObject mockObjects)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when update user, it should success without error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().UpdateUser(gomock.Any(), mockArgs.userId).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "when update user but failed to update user, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().UpdateUser(gomock.Any(), mockArgs.userId).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := service.NewUserService(mockObject.userRepo, nil, nil, nil, nil, nil, nil, time.Millisecond*500, nil)

			if tt.beforeTest != nil {
				tt.beforeTest(mockObject)
			}

			err := w.UpdateUser(context.Background(), tt.args.userId, tt.args.userUpdate)

			if tt.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, tt.expectedErr, err, "Expecting error to be %v", tt.expectedErr)
			} else {
				assert.Nil(t, err, "Expecting no error")
			}
		})
	}
}

func TestUploadKtmImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mockUserRepo.NewMockUserRepository(ctrl)
	storage := mockAws.NewMockCloudStorage(ctrl)

	mockObject := mockObjects{
		userRepo: r,
		storage:  storage,
	}

	type args struct {
		userId           uuid.UUID
		file             *multipart.FileHeader
		BigFile          *multipart.FileHeader
		dir              string
		userWithPhoto    entity.User
		userWtihoutPhoto entity.User
	}

	uuid := uuid.New()

	mockArgs := args{
		uuid,
		&multipart.FileHeader{
			Filename: "testFilename",
			Size:     1000,
		},
		&multipart.FileHeader{
			Filename: "testFilename",
			Size:     10000000,
		},
		"ktm-image/" + uuid.String() + "/testFilename",
		entity.User{
			ID:           uuid,
			Email:        "testEmail",
			Password:     "testPassword",
			KtmImageLink: "testKtmImage",
		},
		entity.User{
			ID:       uuid,
			Email:    "testEmail",
			Password: "testPassword",
		},
	}

	tests := []struct {
		name        string
		args        args
		beforeTest  func(mockObject mockObjects)
		wantErr     bool
		expectedErr error
	}{
		{
			name: "when upload ktm image, it should success without error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{ID: mockArgs.userId}).SetArg(0, mockArgs.userWtihoutPhoto).Return(nil)
				mockObject.storage.EXPECT().Upload(gomock.Any(), mockArgs.file).Return(mockArgs.dir, nil)
				mockObject.userRepo.EXPECT().UpdateUser(gomock.Any(), mockArgs.userId).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "when upload ktm image but file size is too big, it should failed with error file size too big",
			args:        mockArgs,
			beforeTest:  func(mockObject mockObjects) {},
			wantErr:     true,
			expectedErr: domain.ErrFileTooBig,
		},
		{
			name: "when upload ktm image but failed to find user, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{ID: mockArgs.userId}).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when upload ktm(update), it should success without error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{ID: mockArgs.userId}).SetArg(0, mockArgs.userWithPhoto).Return(nil)
				mockObject.storage.EXPECT().Update(gomock.Any(), mockArgs.file, mockArgs.userWithPhoto.KtmImageLink).Return(mockArgs.dir, nil)
				mockObject.userRepo.EXPECT().UpdateUser(gomock.Any(), mockArgs.userId).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "when upload ktm(update) but failed to update ktm, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{ID: mockArgs.userId}).SetArg(0, mockArgs.userWithPhoto).Return(nil)
				mockObject.storage.EXPECT().Update(gomock.Any(), mockArgs.file, mockArgs.userWithPhoto.KtmImageLink).Return(mockArgs.dir, domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when upload ktm but failed to upload, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{ID: mockArgs.userId}).SetArg(0, mockArgs.userWtihoutPhoto).Return(nil)
				mockObject.storage.EXPECT().Upload(gomock.Any(), mockArgs.file).Return(mockArgs.dir, domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when upload ktm but failed to update user, it should failed with error internal server error",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{ID: mockArgs.userId}).SetArg(0, mockArgs.userWtihoutPhoto).Return(nil)
				mockObject.storage.EXPECT().Upload(gomock.Any(), mockArgs.file).Return(mockArgs.dir, nil)
				mockObject.userRepo.EXPECT().UpdateUser(gomock.Any(), mockArgs.userId).Return(domain.ErrInternalServer)
			},
			wantErr:     true,
			expectedErr: domain.ErrInternalServer,
		},
		{
			name: "when upload ktm but operation time exceeded, it should failed with error time out",
			args: mockArgs,
			beforeTest: func(mockObject mockObjects) {
				mockObject.userRepo.EXPECT().FindUser(&entity.User{}, &dto.UserParam{ID: mockArgs.userId}).SetArg(0, mockArgs.userWtihoutPhoto).Return(nil)
				mockObject.storage.EXPECT().Upload(gomock.Any(), mockArgs.file).Return(mockArgs.dir, nil)
				mockObject.userRepo.EXPECT().UpdateUser(gomock.Any(), mockArgs.userId).
					DoAndReturn(func(any any, any2 any) error {
						time.Sleep(time.Millisecond * 1000)
						return nil
					})
			},
			wantErr:     true,
			expectedErr: domain.ErrTimeout,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := service.NewUserService(mockObject.userRepo, nil, nil, nil, nil, nil, nil, time.Millisecond*500, mockObject.storage)

			if tt.beforeTest != nil {
				tt.beforeTest(mockObject)
			}

			var err error

			if i == 1 {
				err = w.UploadKtmImage(context.Background(), tt.args.userId, tt.args.BigFile)
			} else {
				err = w.UploadKtmImage(context.Background(), tt.args.userId, tt.args.file)
			}

			if tt.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
				assert.Equal(t, tt.expectedErr, err, "Expecting error to be %v", tt.expectedErr)
			} else {
				assert.Nil(t, err, "Expecting no error")
			}

		})
	}
}

func TestFetchByParam(t *testing.T) {
	type args struct {
		ctx   context.Context
		user  entity.User
		param dto.UserParam
	}

	gomockCtr := gomock.NewController(t)

	defer gomockCtr.Finish()

	mockObj := mockObjects{
		userRepo: mockUserRepo.NewMockUserRepository(gomockCtr),
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		want       interface{}
		beforeTest func(mockObj mockObjects, args args)
	}{
		{
			name: "When fetch a user by id, it should not return error",
			args: args{
				context.TODO(),
				entity.User{},
				dto.UserParam{
					ID: uuid.MustParse("e15ab4ea-8fc8-4d81-bd7e-4e99eb57e96e"),
				},
			},
			wantErr: false,
			beforeTest: func(mockObj mockObjects, args args) {
				mockObj.userRepo.EXPECT().
					FetchAllByConditionAndRelation("id = ?", []interface{}{args.param.ID}, []string{}, nil,
						"Province",
						"University",
					).
					Return([]entity.User{
						{
							Fullname: "Devan",
						},
					}, dto.PaginationResponse{}, nil)
			},
			want: dto.UserResponse{
				Fullname: "Devan",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := service.NewUserService(mockObj.userRepo, nil, nil, nil, nil, nil, nil, time.Millisecond*500, nil)

			if test.beforeTest != nil {
				test.beforeTest(mockObj, test.args)
			}

			user, err := w.FetchByParam(test.args.ctx, &test.args.param)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}

			assert.Equal(t, test.want, user)
		})
	}
}

func TestLogout(t *testing.T) {
	type args struct {
		ctx      context.Context
		jwtToken string
	}

	gomockCtr := gomock.NewController(t)

	defer gomockCtr.Finish()

	mockObj := mockObjects{
		redis: mockRedis.NewMockRedisInterface(gomockCtr),
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeTest func(mockObj mockObjects, args args)
	}{
		{
			name: "When user logout, it should not return error",
			args: args{
				context.TODO(),
				"jwttoken",
			},
			wantErr: false,
			beforeTest: func(mockObj mockObjects, args args) {
				mockObj.redis.EXPECT().
					Set(gomock.Any(), args.jwtToken, "LOGGED OUT", gomock.Any())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := service.NewUserService(nil, nil, nil, nil, nil, nil, mockObj.redis, time.Millisecond*500, nil)

			if test.beforeTest != nil {
				test.beforeTest(mockObj, test.args)
			}

			err := w.Logout(test.args.ctx, test.args.jwtToken)

			if test.wantErr {
				assert.NotNil(t, err, "Expecting error to be thrown")
			} else {
				assert.Nil(t, err, "Error should not be expected")
			}
		})
	}
}

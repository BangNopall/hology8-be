package service

import (
	"context"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/internal/infra/env"
	"github.com/hology8/hology-be/pkg/aws"
	"github.com/hology8/hology-be/pkg/bcrypt"
	"github.com/hology8/hology-be/pkg/gomail"
	"github.com/hology8/hology-be/pkg/helpers"
	html_content "github.com/hology8/hology-be/pkg/html"
	"github.com/hology8/hology-be/pkg/jwt"
	"github.com/hology8/hology-be/pkg/log"
	"github.com/hology8/hology-be/pkg/redis"
	timePkg "github.com/hology8/hology-be/pkg/time"
	uuidPkg "github.com/hology8/hology-be/pkg/uuid"
)

type userService struct {
	userRepo contracts.UserRepository
	uuid     uuidPkg.UUIDInterface
	bcrypt   bcrypt.BcryptInterface
	time     timePkg.TimeInterface
	goMail   gomail.GoMailInterface
	jwt      jwt.JwtInterface
	redis    redis.RedisInterface
	timeout  time.Duration
	aws      aws.CloudStorage
}

func NewUserService(
	userRepo contracts.UserRepository,
	uuid uuidPkg.UUIDInterface,
	bcrypt bcrypt.BcryptInterface,
	time timePkg.TimeInterface,
	goMail gomail.GoMailInterface,
	jwt jwt.JwtInterface,
	redis redis.RedisInterface,
	timeout time.Duration,
	aws aws.CloudStorage,
) contracts.UserService {
	return &userService{
		userRepo,
		uuid,
		bcrypt,
		time,
		goMail,
		jwt,
		redis,
		timeout,
		aws,
	}
}

func (s *userService) LoginRegisterOauth(ctx context.Context, user dto.UserOauth) (dto.UserLoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var registeredUser entity.User
	err := s.userRepo.FindUser(&registeredUser, &dto.UserParam{Email: user.Email})

	if err == domain.ErrNotFound {
		uuid, err := s.uuid.New()
		if err != nil {
			return dto.UserLoginResponse{}, err
		}

		err = s.userRepo.CreateUser(&entity.User{ID: uuid, Email: user.Email, EmailIsVerified: true})
		if err != nil {
			return dto.UserLoginResponse{}, err
		}

		tokenString, err := s.jwt.GenerateToken(uuid, env.AppEnv.JwtUserRole, "")

		select {
		case <-ctx.Done():
			return dto.UserLoginResponse{}, domain.ErrTimeout
		default:
			return dto.UserLoginResponse{Token: tokenString}, err
		}
	}

	if err == domain.ErrInternalServer {
		return dto.UserLoginResponse{}, err
	}

	if !registeredUser.EmailIsVerified {
		updateUser := dto.UserUpdate{
			EmailIsVerified: true,
		}

		err = s.userRepo.UpdateUser(&updateUser, registeredUser.ID)
		if err != nil {
			return dto.UserLoginResponse{}, err
		}

		tokenString, err := s.jwt.GenerateToken(registeredUser.ID, env.AppEnv.JwtUserRole, "")

		select {
		case <-ctx.Done():
			return dto.UserLoginResponse{}, domain.ErrTimeout
		default:
			return dto.UserLoginResponse{Token: tokenString}, err
		}
	}

	tokenString, err := s.jwt.GenerateToken(registeredUser.ID, env.AppEnv.JwtUserRole, "")

	select {
	case <-ctx.Done():
		return dto.UserLoginResponse{}, domain.ErrTimeout
	default:
		return dto.UserLoginResponse{Token: tokenString}, err
	}
}

func (s *userService) Register(ctx context.Context, user dto.UserRegister, referer string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if user.Password != user.ConfirmPassword {
		return domain.ErrConfirmPasswordNotMatch
	}

	hashPassword, err := s.bcrypt.Hash(user.Password)
	if err != nil {
		return err
	}

	emailVerPassword := helpers.GenerateRandomString(64)

	emailVerPWhash, err := s.bcrypt.Hash(emailVerPassword)
	if err != nil {
		return err
	}

	currentTime := s.time.Now()
	expiredToken := s.time.Add(time.Hour * 1)

	link := "https://hology.ub.ac.id/" + "auth/verify-email/" + user.Email + "/" + emailVerPassword

	// api to call: host + /api/v1/users/verify-email/ + user.Email + / + emailVerPassword

	subject := "Verifikasi Email Akun Hology 8.0"
	HTMLbody := html_content.GetEmailVerifHTML(link)

	sendEmail := func(email string) <-chan error {
		errCh := make(chan error, 1)

		go func() {
			defer close(errCh)
			errCh <- s.goMail.SendEmail(subject, HTMLbody, email)
		}()

		return errCh
	}

	var registeredUser entity.User
	err = s.userRepo.FindUser(&registeredUser, &dto.UserParam{Email: user.Email})

	if err == domain.ErrInternalServer {
		return domain.ErrInternalServer
	}

	expiredTime := time.Date(
		registeredUser.ExpiredToken.Year(),
		registeredUser.ExpiredToken.Month(),
		registeredUser.ExpiredToken.Day(),
		registeredUser.ExpiredToken.Hour(),
		registeredUser.ExpiredToken.Minute(),
		registeredUser.ExpiredToken.Second(),
		registeredUser.ExpiredToken.Nanosecond(),
		time.Local)

	if currentTime.Before(expiredTime) && !registeredUser.EmailIsVerified && registeredUser.Email != "" {
		return domain.ErrCheckEmail
	}

	if currentTime.After(expiredTime) && !registeredUser.EmailIsVerified && registeredUser.Email != "" {
		updateUser := dto.UserUpdate{
			Password:           hashPassword,
			EmailVerifiedToken: emailVerPWhash,
			ExpiredToken:       expiredToken,
		}

		select {
		case <-ctx.Done():
			return domain.ErrTimeout
		case err := <-sendEmail(user.Email):
			if err != nil {
				return err
			}

			err = s.userRepo.UpdateUser(&updateUser, registeredUser.ID)
			if err != nil {
				return err
			}
		}

		return nil
	}

	uuid, err := s.uuid.New()
	if err != nil {
		return err
	}

	newUser := entity.User{
		ID:                 uuid,
		Email:              user.Email,
		Password:           hashPassword,
		EmailVerifiedToken: emailVerPWhash,
		ExpiredToken:       expiredToken,
	}

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	case err := <-sendEmail(user.Email):
		if err != nil {
			return err
		}
		err = s.userRepo.CreateUser(&newUser)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *userService) VerifyEmail(ctx context.Context, email string, emailVerPass string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var registeredUser entity.User
	err := s.userRepo.FindUser(&registeredUser, &dto.UserParam{Email: email})
	if err != nil {
		return err
	}

	expiredTime := time.Date(
		registeredUser.ExpiredToken.Year(),
		registeredUser.ExpiredToken.Month(),
		registeredUser.ExpiredToken.Day(),
		registeredUser.ExpiredToken.Hour(),
		registeredUser.ExpiredToken.Minute(),
		registeredUser.ExpiredToken.Second(),
		registeredUser.ExpiredToken.Nanosecond(),
		time.Local)

	currentTime := s.time.Now()

	if currentTime.After(expiredTime) {
		return domain.ErrInvalidToken
	}

	valid := s.bcrypt.Compare(emailVerPass, registeredUser.EmailVerifiedToken)

	if !valid {
		return domain.ErrInternalServer
	}

	updateUser := dto.UserUpdate{EmailIsVerified: true}

	err = s.userRepo.UpdateUser(&updateUser, registeredUser.ID)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *userService) LoginWithEmail(ctx context.Context, user dto.UserLogin) (dto.UserLoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var registeredUser entity.User

	err := s.userRepo.FindUser(&registeredUser, &dto.UserParam{Email: user.Email})
	if err == domain.ErrNotFound {
		return dto.UserLoginResponse{}, domain.ErrWrongEmailOrPassword
	}

	if err == domain.ErrInternalServer {
		return dto.UserLoginResponse{}, domain.ErrInternalServer

	}

	valid := s.bcrypt.Compare(user.Password, registeredUser.Password)
	if !valid {
		return dto.UserLoginResponse{}, domain.ErrWrongEmailOrPassword

	}

	if !registeredUser.EmailIsVerified {
		return dto.UserLoginResponse{}, domain.ErrCheckEmail
	}

	tokenString, err := s.jwt.GenerateToken(registeredUser.ID, env.AppEnv.JwtUserRole, "")

	select {
	case <-ctx.Done():
		return dto.UserLoginResponse{}, domain.ErrTimeout
	default:
		return dto.UserLoginResponse{Token: tokenString}, err
	}
}

func (s *userService) ResetPassword(ctx context.Context, user dto.UserResetPassword, forgotPasswordToken string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if user.Password != user.ConfirmPassword {
		return domain.ErrConfirmPasswordNotMatch
	}

	var registeredUser entity.User
	err := s.userRepo.FindUser(&registeredUser, &dto.UserParam{ForgotPasswordToken: forgotPasswordToken})
	if err != nil {
		return domain.ErrInvalidToken
	}

	expiredTime := time.Date(
		registeredUser.ExpiredTokenForgot.Year(),
		registeredUser.ExpiredTokenForgot.Month(),
		registeredUser.ExpiredTokenForgot.Day(),
		registeredUser.ExpiredTokenForgot.Hour(),
		registeredUser.ExpiredTokenForgot.Minute(),
		registeredUser.ExpiredTokenForgot.Second(),
		registeredUser.ExpiredTokenForgot.Nanosecond(),
		time.Local)

	currentTime := s.time.Now()

	if currentTime.After(expiredTime) {
		return domain.ErrInvalidToken
	}

	hashPassword, err := s.bcrypt.Hash(user.Password)
	if err != nil {
		return err
	}

	updateUser := dto.UserUpdate{
		Password:           hashPassword,
		ExpiredTokenForgot: time.Now(),
	}

	err = s.userRepo.UpdateUser(&updateUser, registeredUser.ID)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *userService) ForgotPassword(ctx context.Context, user dto.UserForgotPassword, referer string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var registeredUser entity.User
	err := s.userRepo.FindUser(&registeredUser, &dto.UserParam{Email: user.Email})
	if err == domain.ErrNotFound {
		return domain.ErrUserNotFound
	}

	if err == domain.ErrInternalServer {
		return domain.ErrInternalServer
	}

	forgotPasswordToken := helpers.GenerateRandomString(64)

	currentTime := s.time.Now()
	expiredToken := s.time.Add(time.Hour * 1)

	link := "https://hology.ub.ac.id" + `/auth/reset-password/` + forgotPasswordToken

	subject := "Forgot Password"
	HTMLbody := html_content.GetEmailForgotPassword(link)

	expiredTime := time.Date(
		registeredUser.ExpiredTokenForgot.Year(),
		registeredUser.ExpiredTokenForgot.Month(),
		registeredUser.ExpiredTokenForgot.Day(),
		registeredUser.ExpiredTokenForgot.Hour(),
		registeredUser.ExpiredTokenForgot.Minute(),
		registeredUser.ExpiredTokenForgot.Second(),
		registeredUser.ExpiredTokenForgot.Nanosecond(),
		time.Local)

	if currentTime.Before(expiredTime) {
		return domain.ErrCheckEmail
	}

	updateUser := dto.UserUpdate{
		ForgotPasswordToken: forgotPasswordToken,
		ExpiredTokenForgot:  expiredToken,
	}

	err = s.userRepo.UpdateUser(&updateUser, registeredUser.ID)
	if err != nil {
		return err
	}

	err = s.goMail.SendEmail(subject, HTMLbody, user.Email)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *userService) UpdateUser(ctx context.Context, userID uuid.UUID, userUpdate dto.UserUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var (
		err error
	)

	if userUpdate.Password != "" {
		userUpdate.Password, err = s.bcrypt.Hash(userUpdate.Password)

		if err != nil {
			return err
		}
	}

	if userUpdate.BirthDate != "" {

		parsedDate, err := time.Parse("02-01-2006", userUpdate.BirthDate)
		if err != nil {
			return domain.ErrBadRequest
		}

		userUpdate.BirthDate = parsedDate.Format("2006-01-02")
	}

	err = s.userRepo.UpdateUser(&userUpdate, userID)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *userService) UploadKtmImage(ctx context.Context, userId uuid.UUID, ktmImage *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if ktmImage.Size > 1048576 {
		return domain.ErrFileTooBig
	}

	dir := "ktm-image/" + userId.String() + "/" + uuid.NewString() + helpers.GenerateRandomString(10)

	var registeredUser entity.User

	err := s.userRepo.FindUser(&registeredUser, &dto.UserParam{ID: userId})

	if err != nil {
		return err
	}

	var link string

	if registeredUser.KtmImageLink != "" {
		link, err = s.aws.Update(dir, ktmImage, registeredUser.KtmImageLink)
	} else {
		link, err = s.aws.Upload(dir, ktmImage)
	}

	if err != nil {
		return err
	}

	updateUser := dto.UserUpdate{KtmImageLink: link}

	err = s.userRepo.UpdateUser(&updateUser, userId)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *userService) UploadProofImage(ctx context.Context, userId uuid.UUID, proofName string, file *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if proofName != "follow-proof" && proofName != "share-proof" {
		return domain.ErrInvalidProofType
	}

	dir := proofName + "/" + userId.String() + "/" + uuid.NewString() + helpers.GenerateRandomString(10)

	var (
		registeredUser entity.User
		updatedUser    dto.UserUpdate
		link           string
	)

	err := s.userRepo.FindUser(&registeredUser, &dto.UserParam{ID: userId})

	if err != nil {
		return err
	}

	if proofName == "follow-proof" && registeredUser.FollowProofLink != "" {
		link, err = s.aws.Update(dir, file, registeredUser.FollowProofLink)
	} else if proofName == "share-proof" && registeredUser.ShareProofLink != "" {
		link, err = s.aws.Update(dir, file, registeredUser.ShareProofLink)
	} else {
		link, err = s.aws.Upload(dir, file)
	}

	if err != nil {
		return err
	}

	if proofName == "follow-proof" {
		updatedUser = dto.UserUpdate{FollowProofLink: link}
	} else {
		updatedUser = dto.UserUpdate{ShareProofLink: link}
	}

	err = s.userRepo.UpdateUser(&updatedUser, userId)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *userService) FetchByParam(ctx context.Context, userParam *dto.UserParam) (dto.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var (
		args      = []interface{}{userParam.ID}
		condition = "id = ?"
		joins     = []string{}
	)

	preloads := []string{
		"Province",
		"University",
	}

	users, _, err := s.userRepo.FetchAllByConditionAndRelation(
		condition,
		args,
		joins,
		nil,
		preloads...,
	)

	if err != nil {
		return dto.UserResponse{}, err
	}

	if len(users) < 1 {
		return dto.UserResponse{}, domain.ErrNotFound
	}

	userResp := dto.ConvertUserEntityToResponseDto(&users[0])

	select {
	case <-ctx.Done():
		return dto.UserResponse{}, domain.ErrTimeout
	default:
		return *userResp, nil
	}
}

func (s *userService) DeleteUnverifiedUsers() {
	s.userRepo.DeleteUnverifiedUser()
	log.Info(nil, "[USER SERVICE][DeleteUnverifiedUsers] deleted unverified users")
}

func (s *userService) Logout(ctx context.Context, jwtToken string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	expTime := 1
	var err error

	if env.AppEnv.JwtExpireTime != "" {
		expTime, err = strconv.Atoi(env.AppEnv.JwtExpireTime)
	}

	if err != nil {
		return err
	}

	err = s.redis.Set(ctx, jwtToken, "LOGGED OUT", time.Hour*time.Duration(expTime))

	if err != nil {
		return err
	}

	return nil
}

func (s *userService) FetchAll(ctx context.Context, userParam *dto.UserParam, pageParam *dto.PaginationRequest) ([]dto.UserResponse, dto.PaginationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var (
		args      = []interface{}{}
		condition = ""
		joins     = []string{}
	)

	preloads := []string{
		"Province",
		"University",
	}

	users, pageResp, err := s.userRepo.FetchAllByConditionAndRelation(
		condition,
		args,
		joins,
		pageParam,
		preloads...,
	)

	if err != nil {
		if err == domain.ErrNotFound {
			return []dto.UserResponse{}, pageResp, nil
		}

		return nil, pageResp, err
	}

	res := make([]dto.UserResponse, 0)

	for _, user := range users {
		userResp := dto.ConvertUserEntityToResponseDto(&user)

		res = append(res, *userResp)
	}

	select {
	case <-ctx.Done():
		return nil, dto.PaginationResponse{}, domain.ErrTimeout
	default:
		return res, pageResp, nil
	}
}

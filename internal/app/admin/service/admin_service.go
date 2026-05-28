package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hology8/hology-be/domain"
	"github.com/hology8/hology-be/domain/contracts"
	"github.com/hology8/hology-be/domain/dto"
	"github.com/hology8/hology-be/domain/entity"
	"github.com/hology8/hology-be/internal/infra/env"
	"github.com/hology8/hology-be/pkg/bcrypt"
	"github.com/hology8/hology-be/pkg/gomail"
	html_content "github.com/hology8/hology-be/pkg/html"
	"github.com/hology8/hology-be/pkg/jwt"
	"github.com/hology8/hology-be/pkg/log"
)

type adminService struct {
	adminRepo       contracts.AdminRepository
	userRepo        contracts.UserRepository
	competitionRepo contracts.CompetitionRepository
	teamRepo        contracts.TeamRepository
	bcrypt          bcrypt.BcryptInterface
	jwt             jwt.JwtInterface
	goMail          gomail.GoMailInterface
	timeout         time.Duration
}

func NewAdminService(
	adminRepo contracts.AdminRepository,
	userRepo contracts.UserRepository,
	competitionRepo contracts.CompetitionRepository,
	teamRepo contracts.TeamRepository,
	bcrypt bcrypt.BcryptInterface,
	jwt jwt.JwtInterface,
	goMail gomail.GoMailInterface,
	timeout time.Duration,
) contracts.AdminService {
	return &adminService{
		adminRepo,
		userRepo,
		competitionRepo,
		teamRepo,
		bcrypt,
		jwt,
		goMail,
		timeout,
	}
}

func (s *adminService) Login(ctx context.Context, adminLogin dto.AdminLogin) (dto.AdminLoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var admin entity.Admin
	err := s.adminRepo.FindAdmin(&admin, &dto.AdminParam{Username: adminLogin.Username})

	if err != nil {
		if err == domain.ErrNotFound {
			return dto.AdminLoginResponse{}, domain.ErrWrongEmailOrPassword
		}

		return dto.AdminLoginResponse{}, err
	}

	valid := s.bcrypt.Compare(adminLogin.Password, admin.Password)

	if !valid {
		return dto.AdminLoginResponse{}, domain.ErrWrongEmailOrPassword
	}

	token, err := s.jwt.GenerateToken(admin.ID, env.AppEnv.JwtAdminRole, admin.RoleID)

	select {
	case <-ctx.Done():
		return dto.AdminLoginResponse{}, domain.ErrTimeout
	default:
		return dto.AdminLoginResponse{Token: token}, err
	}
}

func (s *adminService) SendEmail(ctx context.Context, to string, emailMessage dto.EmailMessage) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var args []interface{}
	var emails []string

	switch to {
	case "all":
		args = append(args, true)

		users, _, err := s.userRepo.FetchAllByConditionAndRelation("email_is_verified = (?)", args, nil, nil)
		if err != nil {
			return err
		}

		for _, u := range users {
			emails = append(emails, u.Email)
		}

	case "competition":
		if emailMessage.Name == nil {
			return domain.ErrMissingAttribute
		}

		for _, name := range emailMessage.Name {
			args = append(args, name)
		}

		competitions, err := s.competitionRepo.FetchAllByConditionAndRelation(ctx, "competition_name in (?)", args, "Teams.Members.User")
		if err != nil {
			return err
		}

		var leaderID []interface{}
		for _, c := range competitions {
			for _, t := range c.Teams {
				leaderID = append(leaderID, t.LeaderID)
				for _, m := range t.Members {
					emails = append(emails, m.User.Email)
				}
			}
		}

		users, _, err := s.userRepo.FetchAllByConditionAndRelation("id in (?)", leaderID, nil, nil)
		if err != nil {
			return err
		}

		for _, u := range users {
			emails = append(emails, u.Email)
		}

	case "team":
		if emailMessage.Name == nil {
			return domain.ErrMissingAttribute
		}

		for _, name := range emailMessage.Name {
			args = append(args, name)
		}

		teams, _, err := s.teamRepo.FetchAllByConditionAndRelation(ctx, "team_name in (?)", args, "created_at ASC", nil)
		if err != nil {
			return err
		}

		var leaderID []interface{}
		for _, t := range teams {
			leaderID = append(leaderID, t.LeaderID)
		}

		users, _, err := s.userRepo.FetchAllByConditionAndRelation("id in (?)", leaderID, nil, nil)
		if err != nil {
			return err
		}

		for _, u := range users {
			emails = append(emails, u.Email)
		}
	default:
		return domain.ErrMissingAttribute
	}

	subject := emailMessage.Subject
	htmlBody := html_content.GetBasicEmail(emailMessage.Content)

	s.sendEmailBatch(subject, htmlBody, emails)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return nil
	}
}

func (s *adminService) sendEmailBatch(subject string, htmlBody string, emails []string) {
	batchSize := 50
	dailyLimit := 1000
	var wg sync.WaitGroup

	errChan := make(chan error, 1)

	totalEmails := len(emails)
	immediateBatches := totalEmails
	if totalEmails > dailyLimit {
		immediateBatches = dailyLimit

	}

	for i := 0; i < immediateBatches; i += batchSize {
		end := i + batchSize
		if end > immediateBatches {
			end = immediateBatches
		}

		info := fmt.Sprintf("[ADMIN SERVICE][SendEmailBatch]Sending immediate batch: %d - %d", i, end)
		log.Info(log.LogInfo{
			"Info": "Sending email to 1000 users",
		}, info)

		wg.Add(1)
		go s.sendEmail(subject, htmlBody, emails[i:end], &wg, errChan)
	}

	remaingEmails := emails[immediateBatches:]
	delay := time.Hour*24 + time.Minute*1
	for len(remaingEmails) > 0 {
		batch := remaingEmails
		if len(batch) > dailyLimit {
			batch = remaingEmails[:dailyLimit]
			remaingEmails = remaingEmails[dailyLimit:]
		} else {
			remaingEmails = nil
		}

		s.scheduleEmailBatch(subject, htmlBody, batch, batchSize, delay, &wg, errChan)
		delay += time.Hour*24 + time.Minute*1
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()
}

func (s *adminService) sendEmail(
	subject string,
	htmlBody string,
	emails []string,
	wg *sync.WaitGroup,
	errChan chan<- error) {
	defer wg.Done()
	err := s.goMail.SendEmails(subject, htmlBody, emails)
	if err != nil {
		log.Warn(log.LogInfo{
			"error": err.Error(),
		}, "[ADMIN SERVICE][Send Email Batch] failed to send email")
		errChan <- err
	} else {
		log.Info(log.LogInfo{
			"Info": "Success to send email",
		}, "[ADMIN SERVICE][Send Email Batch] Success to send email")
	}
}

func (s *adminService) scheduleEmailBatch(
	subject string,
	htmlBody string,
	emails []string,
	batchSize int,
	delay time.Duration,
	wg *sync.WaitGroup,
	errChan chan<- error) {
	wg.Add(1)

	time.AfterFunc(delay, func() {
		for i := 0; i < len(emails); i += batchSize {
			end := i + batchSize
			if end > len(emails) {
				end = len(emails)
			}

			info := fmt.Sprintf("[ADMIN SERVICE][Send Email Batch] Scheduling batch: %d - %d, delay: %v", i, end, delay)
			log.Info(log.LogInfo{
				"Info": "Scheduling batch",
			}, info)

			wg.Add(1)
			go s.sendEmail(subject, htmlBody, emails[i:end], wg, errChan)
		}
		wg.Done()
	})
}

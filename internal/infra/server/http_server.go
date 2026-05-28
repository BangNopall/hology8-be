package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	_ "github.com/BangNopall/hology8-be/docs"

	announcementCtr "github.com/BangNopall/hology8-be/internal/app/announcement/controller"
	announcementRepo "github.com/BangNopall/hology8-be/internal/app/announcement/repository"
	announcementSvc "github.com/BangNopall/hology8-be/internal/app/announcement/service"
	"github.com/BangNopall/hology8-be/internal/infra/env"

	competitionCtr "github.com/BangNopall/hology8-be/internal/app/competition/controller"
	competitionRepo "github.com/BangNopall/hology8-be/internal/app/competition/repository"
	competitionSvc "github.com/BangNopall/hology8-be/internal/app/competition/service"

	teamCtr "github.com/BangNopall/hology8-be/internal/app/team/controller"
	teamRepo "github.com/BangNopall/hology8-be/internal/app/team/repository"
	teamSvc "github.com/BangNopall/hology8-be/internal/app/team/service"

	adminCtr "github.com/BangNopall/hology8-be/internal/app/admin/controller"
	adminRepo "github.com/BangNopall/hology8-be/internal/app/admin/repository"
	adminSvc "github.com/BangNopall/hology8-be/internal/app/admin/service"

	userCtr "github.com/BangNopall/hology8-be/internal/app/user/controller"
	userRepo "github.com/BangNopall/hology8-be/internal/app/user/repository"
	userSvc "github.com/BangNopall/hology8-be/internal/app/user/service"

	univCtr "github.com/BangNopall/hology8-be/internal/app/university/controller"
	univRepo "github.com/BangNopall/hology8-be/internal/app/university/repository"
	univSvc "github.com/BangNopall/hology8-be/internal/app/university/service"

	provCtr "github.com/BangNopall/hology8-be/internal/app/province/controller"
	provRepo "github.com/BangNopall/hology8-be/internal/app/province/repository"
	provSvc "github.com/BangNopall/hology8-be/internal/app/province/service"

	logCtr "github.com/BangNopall/hology8-be/internal/app/log/controller"
	logRepo "github.com/BangNopall/hology8-be/internal/app/log/repository"
	logSvc "github.com/BangNopall/hology8-be/internal/app/log/service"

	partnerCtr "github.com/BangNopall/hology8-be/internal/app/partner/controller"
	partnerRepo "github.com/BangNopall/hology8-be/internal/app/partner/repository"
	partnerSvc "github.com/BangNopall/hology8-be/internal/app/partner/service"

	voucherCtr "github.com/BangNopall/hology8-be/internal/app/voucher/controller"
	voucherRepo "github.com/BangNopall/hology8-be/internal/app/voucher/repository"
	voucherSvc "github.com/BangNopall/hology8-be/internal/app/voucher/service"

	presenceCtr "github.com/BangNopall/hology8-be/internal/app/presence/controller"
	presenceRepo "github.com/BangNopall/hology8-be/internal/app/presence/repository"
	presenceSvc "github.com/BangNopall/hology8-be/internal/app/presence/service"

	utilsCtr "github.com/BangNopall/hology8-be/internal/app/utils/controller"

	"github.com/BangNopall/hology8-be/internal/middlewares"
	"github.com/BangNopall/hology8-be/pkg/aws"
	"github.com/BangNopall/hology8-be/pkg/bcrypt"
	validators "github.com/BangNopall/hology8-be/pkg/gin"
	"github.com/BangNopall/hology8-be/pkg/gomail"
	"github.com/BangNopall/hology8-be/pkg/jwt"
	"github.com/BangNopall/hology8-be/pkg/log"
	"github.com/BangNopall/hology8-be/pkg/oauth"
	"github.com/BangNopall/hology8-be/pkg/redis"
	timePkg "github.com/BangNopall/hology8-be/pkg/time"
	"github.com/BangNopall/hology8-be/pkg/uuid"
)

type Server interface {
	Start(port string)
	MountMiddlewares()
	RegistCustomValidation()
	MountRoutes(db *gorm.DB)
}

type httpServer struct {
	app       *gin.Engine
	scheduler *cron.Cron
}

func NewHttpServer() Server {

	if env.AppEnv.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	app := gin.Default()
	scheduler := cron.New()

	return &httpServer{app, scheduler}
}

func (s *httpServer) Start(port string) {
	if port[0] != ':' {
		port = ":" + port
	}

	err := s.app.Run(port)

	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[SERVER][Start] failed to start server")
	}
}

func (s *httpServer) MountMiddlewares() {

	if env.AppEnv.AppEnv == "development" {
		url := ginSwagger.URL(`/swagger/doc.json`)
		s.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	s.app.Use(middlewares.CORS())
	s.app.Use(middlewares.ApiKey())
}

func (s *httpServer) RegistCustomValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("alphnumsympace", validators.Alphnumsympace)
		v.RegisterValidation("plusnumeric", validators.Plusnumeric)
		v.RegisterValidation("validdate", validators.DateValidation)
		v.RegisterValidation("validpassword", validators.PasswordValidation)
	}
}

func (s *httpServer) MountRoutes(db *gorm.DB) {
	oauth := oauth.Oauth
	uuid := uuid.UUID
	bcrypt := bcrypt.Bcrypt
	gomail := gomail.Gomail
	jwt := jwt.Jwt
	timePkg := timePkg.Time
	redis := redis.NewRedisClient()
	storage := aws.NewS3Storage()

	// Bootstrap repository, service and controller in here

	// repositories
	announcementRepo := announcementRepo.NewAnnouncementRepository(db)
	competitionRepo := competitionRepo.NewCompetitionRepository(db)
	teamRepo := teamRepo.NewTeamRepository(db)
	adminRepo := adminRepo.NewAdminRepository(db)
	userRepo := userRepo.NewUserRepository(db)
	univRepo := univRepo.NewUniversityRepository(db)
	provRepo := provRepo.NewProvinceRepository(db)
	logRepo := logRepo.NewLogRepository(db)
	partnerRepo := partnerRepo.NewPartnerRepository(db)
	voucherRepo := voucherRepo.NewVoucherRepository(db)
	presenceRepo := presenceRepo.NewPresenceRepository(db)

	// middleware
	middleware := middlewares.NewMiddleware(
		jwt,
		adminRepo,
		teamRepo,
		userRepo,
		redis,
	)

	// services
	univSvc := univSvc.NewUniversityService(univRepo, time.Second*15)
	provSvc := provSvc.NewProvinceService(provRepo, time.Second*15)
	userSvc := userSvc.NewUserService(userRepo, uuid, bcrypt, timePkg, gomail, jwt, redis, time.Second*15, storage)
	adminService := adminSvc.NewAdminService(adminRepo, userRepo, competitionRepo, teamRepo, bcrypt, jwt, gomail, time.Second*15)
	teamService := teamSvc.NewTeamService(teamRepo, competitionRepo, userRepo, time.Second*15, storage)
	competitionService := competitionSvc.NewCompetitionService(competitionRepo, time.Second*15)
	announcementService := announcementSvc.NewAnnouncementService(announcementRepo, teamRepo, competitionRepo, time.Second*15)
	logService := logSvc.NewLogService(logRepo, time.Second*15) // log service
	partnerService := partnerSvc.NewPartnerService(partnerRepo, time.Second*15, storage)
	voucherService := voucherSvc.NewVoucherService(voucherRepo, teamRepo, time.Second*15, db)
	presenceService := presenceSvc.NewPresenceService(presenceRepo)

	// controllers
	univCtr.InitUniversityController(univSvc, s.app)
	provCtr.InitProvinceController(provSvc, s.app)
	utilsCtr.InitUtilsController(s.app)
	userCtr.InitUserController(userSvc, s.app, oauth, middleware, redis)
	adminCtr.InitAdminController(adminService, logService, s.app, middleware)
	teamCtr.InitTeamController(teamService, logService, s.app, middleware)
	logCtr.InitLogController(logService, s.app, middleware)
	announcementCtr.InitAnnouncementController(announcementService, logService, s.app, middleware)
	competitionCtr.InitCompetitionController(competitionService, logService, s.app, middleware)
	partnerCtr.InitPartnerController(partnerService, s.app, middleware)
	voucherCtr.InitVoucherController(voucherService, s.app, middleware)
	presenceCtr.InitPresenceController(presenceService, logService, s.app, middleware)

	// now running cronjobs
	_, err := s.scheduler.AddFunc("0 0 * * 0", userSvc.DeleteUnverifiedUsers)

	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[HTTP SERVER][Mount routes] failed to add cron job delete unverified users")
	}

	s.scheduler.Start()
}

package main

import (
	"os"

	"github.com/BangNopall/hology8-be/internal/infra/database"
	"github.com/BangNopall/hology8-be/internal/infra/env"
	"github.com/BangNopall/hology8-be/internal/infra/server"
)

// @title						Hology 8 API
// @version					1.0
// @description				This is Hology 8 API Documentation
// @host						api.hology.id
// @schemes					https
// @BasePath 				/api/v1
// @securityDefinitions.apiKey	UserAuth
// @in							header
// @name						Authorization
// @description				API Key for accessing protected user and admin endpoints. Type: Bearer TOKEN
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						x-api-key
// @description				API Key for accessing all endpoints. Type: Key TOKEN
func main() {
	server := server.NewHttpServer()
	pgsqldb := database.NewPgsqlConn()

	database.Migrate(pgsqldb, os.Args)
	database.Seeder(pgsqldb, os.Args)

	server.MountMiddlewares()
	server.MountRoutes(pgsqldb)
	server.RegistCustomValidation()
	server.Start(env.AppEnv.AppPort)
}

package database

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/BangNopall/hology8-be/domain/entity"
	"github.com/BangNopall/hology8-be/internal/infra/env"
	"github.com/BangNopall/hology8-be/pkg/bcrypt"
	"github.com/BangNopall/hology8-be/pkg/helpers"
	"github.com/BangNopall/hology8-be/pkg/helpers/flag"
	"github.com/BangNopall/hology8-be/pkg/log"
)

const SEEDERS_FILE_PATH = "data/seeders/"
const SEEDERS_DEV_PATH = SEEDERS_FILE_PATH + "dev/"
const SEEDERS_PROD_PATH = SEEDERS_FILE_PATH + "prod/"

func NewPgsqlConn() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		env.AppEnv.DBHost,
		env.AppEnv.DBUser,
		env.AppEnv.DBPass,
		env.AppEnv.DBName,
		env.AppEnv.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})

	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[PGSQL CONN][NewPgsqlConn] Failed to connect to database")
	}

	return db
}

func Migrate(db *gorm.DB, args []string) {
	if flag.FlagVars.Fresh {

		if env.AppEnv.AppEnv == "production" {
			var choice string
			fmt.Print("Application is on production. Are you sure you want to do fresh migration ? (y/n): ")
			fmt.Scan(&choice)

			if choice != "y" {
				fmt.Print("Exiting...\n")
				os.Exit(0)
			}
		}

		log.Info(nil, "[PGSQL CONN][Migrate] Dropping All Tables")
		db.Migrator().DropTable(getInterfaces()...)

	}

	log.Info(nil, "[PGSQL CONN][Migrate] Auto Migrating Tables")

	db.Exec(`
		DO $$ BEGIN
			CREATE TYPE status AS ENUM (
				'Waiting',
				'Rejected',
				'Verified'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;

		DO $$ BEGIN
			CREATE TYPE phase AS ENUM (
				'Elimination',
				'Final',
				'Disqualified',
				'Eliminated'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;

		DO $$ BEGIN
			CREATE TYPE winner_place AS ENUM (
				'First',
				'Second',
				'Third',
				'Default'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`)

	db.AutoMigrate(getInterfaces()...)
}

func Seeder(db *gorm.DB, args []string) {
	if !flag.FlagVars.Seeder {
		return
	}

	seedersPath := SEEDERS_PROD_PATH

	if env.AppEnv.AppEnv == "development" {
		seedersPath = SEEDERS_DEV_PATH
	}

	models := []string{
		"role",
		"admin",
		"province",
		"university",
		"competition",
		"partner_type",
		"voucher",
	}

	seeders := map[string][]interface{}{
		// roles data
		"role": makeSeeders(
			seedersPath+"roles.txt",
			reflect.TypeOf(entity.Role{}),
			",",
			"ID",
			"Name",
		),

		// admins data
		"admin": makeSeeders(
			seedersPath+"admins.txt",
			reflect.TypeOf(entity.Admin{}),
			",",
			"Fullname",
			"Username",
			"Password",
			"RoleID",
		),

		// provinces data
		"province": makeSeeders(
			seedersPath+"provinces.txt",
			reflect.TypeOf(entity.Province{}),
			"",
			"Name",
		),

		// universities data
		"university": makeSeeders(
			seedersPath+"universities.txt",
			reflect.TypeOf(entity.University{}),
			"",
			"Name",
		),

		// competitions
		"competition": makeSeeders(
			seedersPath+"competitions.txt",
			reflect.TypeOf(entity.Competition{}),
			",",
			"Name",
			"Desc",
		),

		// partner types
		"partner_type": makeSeeders(
			seedersPath+"partner_types.txt",
			reflect.TypeOf(entity.PartnerType{}),
			"",
			"Name",
		),

		//vouchers
		"voucher": makeSeeders(
			seedersPath+"vouchers.txt",
			reflect.TypeOf(entity.Voucher{}),
			"",
			"ID",
		),
	}

	if len(models) != len(seeders) {
		log.Fatal(nil, "[PGSQL CONN] models length differs with seeders length")
	}

	if flag.FlagVars.SeederModel != "" {
		models := strings.Split(flag.FlagVars.SeederModel, ",")
		toRun := make([]string, 0)

		for _, model := range models {
			_, exists := seeders[model]

			if !exists {
				log.Warn(nil, fmt.Sprintf("[PGSQL CONN][Seeder] cant find model %v in seeders data", model))
				continue
			}

			toRun = append(toRun, model)
		}

		for seedModel := range seeders {
			erase := true

			for _, model := range toRun {
				if seedModel == model {
					erase = false
					break
				}
			}

			if erase {
				seeders[seedModel] = []interface{}{}
			}
		}
	}

	createSeeders(
		db,
		models,
		seeders,
	)
}

func makeSeeders(
	filepath string,
	modelType reflect.Type,
	separator string,
	fields ...string,
) []interface{} {
	seeders := make([]interface{}, 0)

	seedersData, _ := helpers.ReadFile(filepath, separator)

	for _, data := range seedersData {
		seeder := reflect.New(modelType).Interface()

		dataFields := []string{data}

		if separator != "" {
			dataFields = strings.Split(data, separator)
		}

		for idx, field := range fields {
			seederField := reflect.ValueOf(seeder).Elem().FieldByName(field)

			if !seederField.IsValid() || !seederField.CanSet() {
				log.Warn(nil, fmt.Sprintf("Model type %s does not have field %s\n", modelType, field))

			} else {
				trimmed := strings.TrimSpace(dataFields[idx])

				// Handle special case lol
				if field == "Password" {
					trimmed, _ = bcrypt.Bcrypt.Hash(trimmed)
				}

				seederField.SetString(trimmed)
			}
		}

		seeders = append(seeders, seeder)
	}

	return seeders
}

func createSeeders(db *gorm.DB, models []string, allSeeders map[string][]interface{}) {

	for _, model := range models {

		if len(allSeeders[model]) == 0 {
			continue
		}

		log.Info(nil, fmt.Sprintf("[PGSQL CONN][createSeeders] Generating seeders for model : %v", model))

		for idx, seed := range allSeeders[model] {
			res := db.Create(seed)

			if res.Error != nil {
				log.Error(log.LogInfo{
					"error": res.Error.Error(),
				}, fmt.Sprintf("[PGSQL CONN][createSeeders] failed to create seeder at index data : %v", idx))
			}
		}

	}
}

func getInterfaces() []interface{} {
	return []interface{}{
		&entity.Province{},
		&entity.University{},
		&entity.Role{},
		&entity.Competition{},
		&entity.Admin{},
		&entity.User{},
		&entity.Log{},
		&entity.PartnerType{},
		&entity.Partner{},
		&entity.Voucher{},
		&entity.Team{},
		&entity.Announcement{},
		&entity.DetailTeams{},
		&entity.Presence{},
	}
}

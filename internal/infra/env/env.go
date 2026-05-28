package env

import (
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/hology8/hology-be/pkg/log"
)

type Env struct {
	AppEnv             string `mapstructure:"APP_ENV"`
	AppPort            string `mapstructure:"APP_PORT"`
	ApiKey             string `mapstructure:"API_KEY"`
	DBHost             string `mapstructure:"DB_HOST"`
	DBPort             string `mapstructure:"DB_PORT"`
	DBUser             string `mapstructure:"DB_USER"`
	DBPass             string `mapstructure:"DB_PASS"`
	DBName             string `mapstructure:"DB_NAME"`
	SupaURL            string `mapstructure:"SUPABASE_URL"`
	SupaKey            string `mapstructure:"SUPABASE_KEY"`
	SupaBucket         string `mapstructure:"SUPABASE_BUCKET"`
	CsrfKey            string `mapstructure:"CSRF_KEY"`
	StoreSecret        string `mapstructure:"STORE_SECRET"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoMailPort         string `mapstructure:"GOMAIL_PORT"`
	GoMailHost         string `mapstructure:"GOMAIL_HOST"`
	GoMailUsername     string `mapstructure:"GOMAIL_USERNAME"`
	GoMailPassword     string `mapstructure:"GOMAIL_PASSWORD"`
	FireBucket         string `mapstructure:"FIREBASE_BUCKET"`
	FireCredPath       string `mapstructure:"FIREBASE_CREDENTIALS_PATH"`
	RedisHost          string `mapstructure:"REDIS_HOST"`
	RedisPort          string `mapstructure:"REDIS_PORT"`
	RedisPassword      string `mapstructure:"REDIS_PASS"`
	JwtSecretKey       string `mapstructure:"JWT_SECRET_KEY"`
	JwtExpireTime      string `mapstructure:"JWT_EXP_TIME"`
	JwtUserRole        string `mapstructure:"JWT_USER_ROLE"`
	JwtAdminRole       string `mapstructure:"JWT_ADMIN_ROLE"`
	AwsBaseEndpoint    string `mapstructure:"AWS_BASE_ENDPOINT"`
	AwsBucket          string `mapstructure:"AWS_BUCKET"`
	AwsRegion          string `mapstructure:"AWS_REGION"`
	AwsAccessKey       string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecret          string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
}

var AppEnv = getEnv()

func getEnv() *Env {
	env := &Env{}

	// instantly return env in testing mode
	if strings.Contains(os.Args[0], ".test") {
		return env
	}

	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[ENV][getEnv] failed to read config file")
	}

	if err := viper.Unmarshal(env); err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[ENV][getEnv] failed to unmarshal to struct")
	}

	switch env.AppEnv {
	case "development":
		log.Info(nil, "Application is running on development mode")
	case "production":
		log.Info(nil, "Application is running on production mode")
	default:
		log.Fatal(nil, "[ENV][getEnv] app_mode variables is undefined")
	}

	return env
}

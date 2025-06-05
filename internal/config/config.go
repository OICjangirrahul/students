package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// HTTPサーバー設定：サーバーのアドレスとタイムアウト設定を管理
type HTTPServer struct {
	// サーバーのリッスンアドレス（例：":8082"）
	Addr    string        `yaml:"addr" env:"HTTP_ADDR"`
	Timeout time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT"`
}

// データベース設定：PostgreSQLデータベースへの接続情報を管理
type DBConfig struct {
	// データベースホスト
	Host string `yaml:"host" env:"DB_HOST"`
	// データベースポート
	Port string `yaml:"port" env:"DB_PORT"`
	// データベースユーザー名
	User string `yaml:"user" env:"DB_USER"`
	// データベースパスワード
	Password string `yaml:"password" env:"DB_PASSWORD"`
	// データベース名
	Name string `yaml:"name" env:"DB_NAME"`
}

// JWT設定：JSON Web Tokenの生成と検証に関する設定を管理
type JWTConfig struct {
	// トークン署名用の秘密鍵
	Secret string `yaml:"secret" env:"JWT_SECRET"`
	// トークンの有効期限
	Expiration string `yaml:"expiration" env:"JWT_EXPIRATION"`
}

// AWS設定：AWSサービスへのアクセス設定を管理
type AWSConfig struct {
	// AWSリージョン
	Region string `yaml:"region" env:"AWS_REGION"`
	// S3バケット名
	S3Bucket string `yaml:"s3_bucket" env:"AWS_S3_BUCKET"`
	// DynamoDBテーブル名
	DynamoTable string `yaml:"dynamo_table" env:"AWS_DYNAMO_TABLE"`
}

// ログ設定：アプリケーションのログ出力設定を管理
type LogConfig struct {
	// ログレベル（debug, info, warn, error）
	Level string `yaml:"level" env:"LOG_LEVEL"`
	// ログフォーマット（json, text）
	Format string `yaml:"format" env:"LOG_FORMAT"`
}

// アプリケーション全体の設定を管理する構造体
type Config struct {
	// 実行環境（dev, prod, test）
	Env string `yaml:"env" env:"ENV"`
	// データベース設定
	DB DBConfig `yaml:"db"`
	// HTTPサーバー設定
	HTTPServer HTTPServer `yaml:"http_server"`
	// JWT設定
	JWT JWTConfig `yaml:"jwt"`
	// AWS設定
	AWS AWSConfig `yaml:"aws"`
	// ログ設定
	Log LogConfig `yaml:"log"`
}

// 設定をロードする
// YAMLファイルと環境変数から設定を読み込む
// 環境変数はYAMLの設定より優先される
func LoadConfig(path string) (*Config, error) {
	// .envファイルが存在する場合は読み込む
	if err := godotenv.Load(); err != nil {
		// .envファイルが存在しない場合はエラーを無視
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	// デフォルト設定を初期化
	cfg := &Config{
		Env: getEnv("ENV", "dev"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "students"),
		},
		HTTPServer: HTTPServer{
			Addr:    getEnv("HTTP_ADDR", ":8082"),
			Timeout: time.Duration(getEnvAsInt("HTTP_TIMEOUT", 30)) * time.Second,
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-here"),
			Expiration: getEnv("JWT_EXPIRATION", "24h"),
		},
		AWS: AWSConfig{
			Region:      getEnv("AWS_REGION", "us-west-2"),
			S3Bucket:    getEnv("AWS_S3_BUCKET", ""),
			DynamoTable: getEnv("AWS_DYNAMO_TABLE", ""),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "debug"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	// YAMLパスが指定されている場合、YAMLファイルを読み込んでマージ
	if path != "" {
		if err := loadYAMLConfig(path, cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// YAMLファイルから設定を読み込み、既存の設定とマージする
func loadYAMLConfig(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	return nil
}

// AWS SDKの設定をロードする
func LoadAWSConfig(cfg *Config) (aws.Config, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.AWS.Region),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	return awsCfg, nil
}

// 環境変数から値を取得し、存在しない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// 環境変数から整数値を取得し、存在しない場合はデフォルト値を返す
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := time.ParseDuration(value); err == nil {
			return int(intValue.Seconds())
		}
	}
	return defaultValue
}

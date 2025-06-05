package repositories

import (
	"fmt"

	"github.com/OICjangirrahul/students/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// データベース接続を作成する
// 設定情報を受け取り、PostgreSQLデータベースへの接続を確立する
func NewDB(cfg *config.Config) (*gorm.DB, error) {
	// データベース接続文字列を生成
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)

	// GORMを使用してデータベースに接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 接続を検証
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// データベースにPingを送信して接続を確認
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

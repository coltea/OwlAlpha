package bootstrap

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/coltea/owlalpha/backend/internal/model/entity"
)

type Config struct {
	Server struct {
		Address   string `json:"address"`
		JWTSecret string `json:"jwtSecret"`
	} `json:"server"`
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
		SSLMode  string `json:"sslMode"`
	} `json:"database"`
	Redis struct {
		Address  string `json:"address"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
	OpenAI struct {
		BaseURL string `json:"baseUrl"`
		APIKey  string `json:"apiKey"`
		Model   string `json:"model"`
	} `json:"openai"`
}

type Dependencies struct {
	Config *Config
	DB     *gorm.DB
	Redis  *redis.Client
}

func New(ctx context.Context) (*Dependencies, error) {
	var cfg Config
	for _, item := range []struct {
		key string
		dst any
	}{
		{"server", &cfg.Server},
		{"database", &cfg.Database},
		{"redis", &cfg.Redis},
		{"openai", &cfg.OpenAI},
	} {
		if v := g.Cfg().MustGet(ctx, item.key); v != nil {
			if err := v.Scan(item.dst); err != nil {
				return nil, fmt.Errorf("scan config.%s: %w", item.key, err)
			}
		}
	}

	applyEnvOverrides(&cfg)

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	if err := db.AutoMigrate(&entity.User{}, &entity.Report{}, &entity.ModelConfig{}); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	return &Dependencies{
		Config: &cfg,
		DB:     db,
		Redis:  rdb,
	}, nil
}

func applyEnvOverrides(cfg *Config) {
	if value := getenv("SERVER_ADDRESS", cfg.Server.Address); value != "" {
		cfg.Server.Address = value
	}
	if value := getenv("SERVER_JWTSECRET", cfg.Server.JWTSecret); value != "" {
		cfg.Server.JWTSecret = value
	}

	cfg.Database.Host = getenv("DATABASE_HOST", cfg.Database.Host)
	cfg.Database.User = getenv("DATABASE_USER", cfg.Database.User)
	cfg.Database.Password = getenv("DATABASE_PASSWORD", cfg.Database.Password)
	cfg.Database.Name = getenv("DATABASE_NAME", cfg.Database.Name)
	cfg.Database.SSLMode = getenv("DATABASE_SSLMODE", cfg.Database.SSLMode)
	if value := getenv("DATABASE_PORT", ""); value != "" {
		if port, err := strconv.Atoi(value); err == nil {
			cfg.Database.Port = port
		}
	}

	cfg.Redis.Address = getenv("REDIS_ADDRESS", cfg.Redis.Address)
	cfg.Redis.Password = getenv("REDIS_PASSWORD", cfg.Redis.Password)
	if value := getenv("REDIS_DB", ""); value != "" {
		if db, err := strconv.Atoi(value); err == nil {
			cfg.Redis.DB = db
		}
	}

	cfg.OpenAI.BaseURL = getenv("OPENAI_BASEURL", cfg.OpenAI.BaseURL)
	cfg.OpenAI.APIKey = getenv("OPENAI_APIKEY", cfg.OpenAI.APIKey)
	cfg.OpenAI.Model = getenv("OPENAI_MODEL", cfg.OpenAI.Model)
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

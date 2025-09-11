package config

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/alimosavifard/zyros-backend/migrations"
	"github.com/alimosavifard/zyros-backend/utils"
)

var DB *gorm.DB
var RedisClient *redis.Client
var onceDB, onceRedis sync.Once

var logger = utils.InitLogger()

func ConnectDB(cfg *Config) *gorm.DB {
	onceDB.Do(func() {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.DB_HOST,
			cfg.DB_USER,
			cfg.DB_PASSWORD,
			cfg.DB_NAME,
			cfg.DB_PORT)
		var err error
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to connect to database")
		}
		sqlDB, _ := DB.DB()
		maxOpenConns, _ := strconv.Atoi(cfg.DB_MAX_OPEN_CONNS)
		maxIdleConns, _ := strconv.Atoi(cfg.DB_MAX_IDLE_CONNS)
		sqlDB.SetMaxOpenConns(maxOpenConns)
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetConnMaxLifetime(0)
		logger.Info().Msg("Database connection established successfully")

		if err := migrations.RunMigrations(DB); err != nil {
			logger.Fatal().Err(err).Msg("Failed to run migrations")
		}
	})
	return DB
}

func ConnectRedis(cfg *Config) *redis.Client {
	onceRedis.Do(func() {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.REDIS_ADDR,
			Password: cfg.REDIS_PASSWORD,
			DB:       0,
		})
		_, err := RedisClient.Ping(context.Background()).Result()
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to connect to Redis")
		}
		logger.Info().Msg("Redis connection established successfully")
	})
	return RedisClient
}
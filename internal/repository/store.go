package repository

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/kelein/trove-fiber/pkg/log"
)

// NewRedis creates a new Redis client
func NewRedis(conf *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		DB:       conf.GetInt("data.redis.db"),
		Addr:     conf.GetString("data.redis.addr"),
		Password: conf.GetString("data.redis.password"),
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}
	return rdb
}

// NewDB creates a new GORM database connection
func NewDB(conf *viper.Viper) *gorm.DB {
	var db *gorm.DB
	var err error
	driver := conf.GetString("data.db.user.driver")
	dsn := conf.GetString("data.db.user.dsn")

	switch driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			DisableAutomaticPing: false,
			Logger:               log.DefaultOrmlogger(),
		})

	case "postgres":
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})

	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			DisableAutomaticPing: false,
			Logger:               log.DefaultOrmlogger(),
		})

	default:
		panic("unknown db driver")
	}
	if err != nil {
		panic(err)
	}
	db = db.Debug()

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}

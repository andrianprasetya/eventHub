package server

import (
	"fmt"
	"github.com/andrianprasetya/eventHub/database"
	"github.com/andrianprasetya/eventHub/internal/shared/redisser"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"os"
	"strconv"
)

var (
	Database    *gorm.DB
	RedisClient redisser.RedisClient
)

func InitDatabase() *gorm.DB {
	Database = database.GetConnection() //Migrate DB
	database.MigrateDatabase()
	//Seeding DB
	database.SeedUsers(Database)
	return Database
}

// InitRedis redis
func InitRedis() redisser.RedisClient {

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	username := os.Getenv("REDIS_USERNAME")
	dbName := os.Getenv("REDIS_DB")
	db, err := strconv.Atoi(dbName)

	if err != nil {
		panic(fmt.Sprintf("Invalid REDIS_DB value must integer: %v", err))
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Username: username, // no user set
		Password: password, // no password set
		DB:       db,       // use default DB
	})
	RedisClient = redisser.NewRedisClient(rdb)
	return RedisClient
}

package database

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func RD_Init() *redis.Client {
	rd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password
		DB:       0,                // Default DB
	})

	fmt.Println("Success : Success to connect redis.")

	return rd
}

package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
)

var redisDB *redis.Client

const CloudfalrePrefix = "st:domain:cloudflare:"

func ConnectRedis() {
	pong, err := connect()
	fmt.Println(pong, err)
	if err != nil {
		log.Println("Failed to connect redis: ", err)
		log.Println("Wait for 5 second , Try to reconnect redis")
		time.Sleep(5 * time.Second)
		pong, err = connect()
		if err != nil {
			log.Println("Give up to connect to redis")
			return
		}
	}

	log.Println("Successfully connected to redis.")
}

func connect() (string, error) {

	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		redisAddress = "127.0.0.1:6379"
		//redisAddress = "redis-master.db:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisDB = redis.NewClient(&redis.Options{
		// Addr: "172.17.0.2:6379",
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       0,
	})

	return redisDB.Ping().Result()
}

func WriteTestKey(key string, value string) {
	setRes := redisDB.Set(key, value, 0)
	fmt.Println(key, setRes.Val())
	getRes := redisDB.Get(key)
	fmt.Println(key, getRes.Val())
}

func SetKey(key string, value string) string {
	setRes := redisDB.Set(key, value, 0)
	fmt.Println(key, setRes.Val())
	getRes := redisDB.Get(key)
	fmt.Println(key, getRes.Val())
	return getRes.Val()
}

func GetKey(key string) string {
	getRes := redisDB.Get(key)
	fmt.Println(key, getRes.Val())
	return getRes.Val()
}

func WriteDomains() {

}

func Keys(keyPattern string) []string {
	res := redisDB.Keys(keyPattern)
	keyCount := len(res.Val())
	fmt.Println(keyCount)
	return res.Val()
}

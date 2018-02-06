package redisTemple

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/panenming/go-im/libs/config"
)

// 实例
type RedisInstance struct {
	Addr     string
	Password string
	DB       int
}

// redis client
var RedisClient *redis.Client

func init() {
	file := "config/config.json"
	var conf RedisInstance
	err := config.New(file, &conf)
	if err != nil {
		panic(err)
	}

	conf.string()
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})
}

func (instance *RedisInstance) string() {
	log.Printf("{'Addr':'%s','Password':'%s','DB':'%d'} \n", instance.Addr, instance.Password, instance.DB)
}

// TODO

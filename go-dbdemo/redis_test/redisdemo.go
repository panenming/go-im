package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "10.39.35.38:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	err = client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("key = ", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

	// client 有连接池

	err = client.HSet("hash", "f1", 1).Err()
	if err != nil {
		panic(err)
	}

	map1, err := client.HGetAll("hash").Result()
	if err == redis.Nil {
		fmt.Println("hash does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("f1 = ", map1["f1"])
	}

	err = client.ZAdd("zset", redis.Z{Score: 1,
		Member: "name"}).Err()

	if err != nil {
		panic(err)
	}

	zset1, err := client.ZRangeByScore("zset", redis.ZRangeBy{
		Min: "-1",
		Max: "0",
	}).Result()
	if err == redis.Nil {
		fmt.Println("hash does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("zset = ", zset1)
	}

}

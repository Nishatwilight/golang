package redisconfig

import (
	"fmt"
	

	"github.com/go-redis/redis"
)

func RedisConnect() (redis.Client, error) {

	// tlsConfig := &tls.Config{MaxVersion: tls.VersionTLS11, MinVersion: tls.VersionTLS10}

	client := redis.NewClient(&redis.Options{
		// Addr: "localhost:6379",
		Addr:     "srt-cache.redis.cache.windows.net:6379",
		Password: "QhWF9orVyQobOT2vjNPiQjl0tfnWSqZcCAzCaBPFv6k=",
		DB:       0,
		// TLSConfig: tlsConfig,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		return *client, err
	}
	fmt.Println(pong)

	return *client, nil
}




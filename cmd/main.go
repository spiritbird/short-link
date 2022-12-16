package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	S "github.com/spiritbird/short-link"
	"time"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.10.131:6379",
		Password: "123456",
		DB:       2,
	})
	serv := S.NewShortLinkServ(client, "")

	ss, err := serv.Shorten("http://glosku-mall.com", 1*time.Hour)
	fmt.Println(ss, err)
}

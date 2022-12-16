package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spiritbird/short-link"
	"time"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.10.131:6379",
		Password: "123456",
		DB:       0,
	})
	var serv = short_link.NewShortLinkServ(client, "")
	ss, err := serv.Shorten("https://google.com", 1*time.Hour)
	fmt.Println(ss, err)
}

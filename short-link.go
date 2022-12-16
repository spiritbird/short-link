package short_link

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spiritbird/short-link/tools"
	"strconv"
	"strings"
	"time"
)

const (
	mask    = 0x3fffffff
	charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	ShortLinkKey = "shortLink:%s"

	URLHashKey = "urlHash:%s"
)

type ShortLinkServ interface {
	Shorten(string, time.Duration) (string, error)
	GetTarget(string) (string, error)
}

type shortLink struct {
	Client *redis.Client
	Host   string
}

func NewShortLinkServ(connect *redis.Client, host string) *shortLink {
	return &shortLink{
		Client: connect,
		Host:   host,
	}
}

func (serv *shortLink) Shorten(longUrl string, expr time.Duration) (string, error) {
	// 计算长网址的 MD5 签名串
	signature := tools.GetMD5Encode(longUrl)
	cachedShortUrl, _ := serv.Client.Get(context.Background(), signature).Result()

	if cachedShortUrl != "" {
		return cachedShortUrl, nil
	}

	// 将 MD5 签名串分成 4 段，每段 8 个字节
	parts := []string{
		signature[:8],
		signature[8:16],
		signature[16:24],
		signature[24:],
	}
	shortedUrl := "" //长网址转换后的短字符

	// 循环处理每一段
	for _, part := range parts {
		// 将 8 个字节看成 16 进制串，并与 0x3fffffff 进行与操作
		hexString := fmt.Sprintf("%x", part)
		hexInt, _ := strconv.Atoi(hexString)
		n := hexInt & mask

		// 将 30 位数字分成 6 段，每 5 位的数字作为字母表的索引取得特定字符
		var shortUrl strings.Builder
		for j := 0; j < 6; j++ {
			idx := n & 0x0000003d
			shortUrl.WriteByte(charset[idx])
			n = n >> 5
		}
		shortedUrl = shortUrl.String()
	}
	pip := serv.Client.Pipeline()

	pip.Set(context.Background(), fmt.Sprintf(ShortLinkKey, shortedUrl), longUrl, expr)
	pip.Set(context.Background(), fmt.Sprintf(URLHashKey, signature), shortedUrl, expr)

	_, err := pip.Exec(context.Background())
	if err != nil {
		return "", ErrRedisConnect
	}
	return shortedUrl, nil
}

func (serv *shortLink) GetTarget(shortUrl string) (string, error) {
	cachedLongUrl, err := serv.Client.Get(context.Background(), shortUrl).Result()
	if err != nil {
		return "", ErrNotFound
	}
	return cachedLongUrl, nil
}

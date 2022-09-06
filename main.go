package main

import (
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

var redisConnection redis.Conn

func main() {
	endpoint := os.Getenv("REDIS_ENDPOINT")
	password := os.Getenv("REDIS_PASSWORD")
	options := redis.DialPassword(password)
	redisConnection, err := redis.Dial("tcp", endpoint, options)
	if err != nil {
		panic(err)
	}
	defer redisConnection.Close()
	r := gin.Default()
	r.GET("/shorten", shorten)
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func shorten(c *gin.Context) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, c.Request.Body)
	if err != nil {
		panic(err)
	}
	str := buf.String()
	var shortener = &URLShortener{}
	url, err := shortener.shorten(str)
	if err != nil {
		panic(err)
	}
	scheme := "http://"
	if c.Request.TLS != nil {
			scheme = "https://"
	}
	url = scheme + c.Request.Host + "/" + url
	_, err = redisConnection.Do("HSET", "url", str, url)
	if err != nil {
		panic(err)
	}
	c.JSON(200, url)
}
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gomodule/redigo/redis"
)

var redisConnection redis.Conn

func main() {
	endpoint := os.Getenv("REDIS_ENDPOINT")
	password := os.Getenv("REDIS_PASSWORD")
	options := redis.DialPassword(password)
	var err error
	redisConnection, err = redis.Dial("tcp", endpoint, options)
	if err != nil {
		panic(err)
	}
	defer redisConnection.Close()
	r := gin.Default()
	r.POST("/shorten", shorten)
	r.GET("/:code", redirect)
  r.Run()
}

func shorten(c *gin.Context) {
	shortenDto := ShortenDTO{}
	if err := c.ShouldBindBodyWith(&shortenDto, binding.JSON); err != nil {
		panic(err)
	}
	var shortener = &URLShortener{}
	code, err := shortener.shorten(shortenDto.Url)
	if err != nil {
		panic(err)
	}
	scheme := "http://"
	if c.Request.TLS != nil {
			scheme = "https://"
	}
	url := scheme + c.Request.Host + "/" + code
	_, err = redisConnection.Do("HSET", "url", code, shortenDto.Url)
	if err != nil {
		panic(err)
	}
	c.JSON(200, url)
}

func redirect(c *gin.Context) {
	code := c.Param("code")
	fmt.Println(code)
	url, err := redis.String(redisConnection.Do("HGET", "url", code))
	if err != nil {
		c.JSON(400, err)
	}
	c.Redirect(http.StatusMovedPermanently, url)
}
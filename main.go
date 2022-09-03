package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

func main() {
	endpoint := os.Getenv("REDIS_ENDPOINT")
	password := os.Getenv("REDIS_PASSWORD")
	options := redis.DialPassword(password)
	redis, err := redis.Dial("tcp", endpoint, options)
	if err != nil {
		panic(err)
	}
	defer redis.Close()
  r := gin.Default()
  r.GET("/ping", func(c *gin.Context) {
		sec:= time.Now().Unix()
		reply, err := redis.Do("SET", "time", sec)
		if err != nil {
			panic(err)
		}
		fmt.Println(reply)
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
			"time": sec,
    })
  })
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
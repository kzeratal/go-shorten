package handler

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"goshorten/internal/redis"
	"net/http"
	"net/url"

	"github.com/catinello/base62"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var ctx = context.Background()

func Shorten(c *gin.Context) {
	URLs := []string{}
	if err := c.ShouldBindBodyWith(&URLs, binding.JSON); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// validate the URLs
	hasher := md5.New()
	shortenedURLs := map[string]string{}
	for _, URL := range URLs {
		if _, err := url.ParseRequestURI(URL); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		hasher.Write([]byte(URL))
		num := binary.BigEndian.Uint64(hasher.Sum(nil))
		data := base62.Encode(int(num))
		shortenedURLs[URL] = data
	}
	for key, value := range shortenedURLs {
		fmt.Println(key, ": ", value)
	}
	pipe := redis.Client.Pipeline()
	for key, value := range shortenedURLs {
		pipe.HSet(ctx, "url", map[string]interface{} { key: value })
	}
	if _, err := pipe.Exec(ctx); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}


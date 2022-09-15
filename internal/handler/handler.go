package handler

import (
	"context"
	"goshorten/internal/redis"
	"net/http"
	"net/url"

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
	for _, URL := range URLs {
		if _, err := url.ParseRequestURI(URL); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}
	if err := redis.Client.HSet(ctx, "url", "hashed").Err(); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}


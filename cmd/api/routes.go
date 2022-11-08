package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Config) routes() http.Handler {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, "worked")
	})

	r.POST("/message", app.SendMessage)
	r.GET("/message/list", app.ListMessage)

	return r
}

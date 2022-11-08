package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Config) SendMessage(c *gin.Context) {
	payload := JsonResponse{
		Sender:   "",
		Receiver: "",
		Message:  "Jello there,",
	}

	app.writeJSON(c, http.StatusOK, payload)
}

func (app *Config) ListMessage(c *gin.Context) {
	app.readJSON(c)

}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/norand94/jsonservice/core"
	"time"
)

var app *core.App

func main(){
	fmt.Println("Started")
	app = core.NewApp()
	r := gin.Default()
	r.POST("/", saveRequest)
	r.Run(":8080")
}

func saveRequest(c *gin.Context) {
	fmt.Println(c.Request.Header.Get("Content-Type"))
	if time.Now().Sub(app.LastReq()) > 5 * time.Second {
		app.IncrCurrPos()
	}
	app.SetLastReq(time.Now())
	c.JSON(200, gin.H{
		"pos" : app.CurrPos(),
	})
}
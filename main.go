package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/norand94/jsonservice/core"
	"github.com/go-redis/redis"
)

var app *core.App
var rcli *redis.Client

func main(){
	fmt.Println("Started")
	rcli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	app = core.NewApp(rcli)

	r := gin.Default()
	r.POST("/", app.SaveRequest)
	r.Run(":8080")
}


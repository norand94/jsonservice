package main

import (
	"fmt"
	"github.com/norand94/jsonservice/core"
	"github.com/go-redis/redis"
	"github.com/valyala/fasthttp"
)

var app *core.App

func main(){
	fmt.Println("Started")
	rcli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	app = core.NewApp(rcli)

	fasthttp.ListenAndServe(":8080", app.SaveRequest)
	//fasthttp.
	//r := gin.Default()
	//r.POST("/", app.SaveRequest)
	//r.GET("/stats", app.Stats)
	//r.Run(":8080")
}




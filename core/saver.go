package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/valyala/fasthttp"
	"time"
	"encoding/json"
)

type JData struct {
	App    JApp    `json:"app"`
	Device JDevice `json:"device"`
}

type JApp struct {
	Bundle string `json:"bundle"`
}

type JDevice struct {
	Geo JGeo   `json:"geo"`
	OS  string `json:"os"`
}

type JGeo struct {
	Country string `json:"country"`
}

func (d JData) Key() string {
	return "stat:" + d.Device.Geo.Country + ":" + d.Device.OS + ":" + d.App.Bundle
}

type StatsReport struct {
	Country string `json:"country"`
	App string `json:"app"`
	Platform string `json:"platform"`
	Count string `json:"count"`
}

func createReport(key string) StatsReport {
	sp := strings.Split(key, ":")
	if sp[0] != "stat" {
		panic("key is not stat! - " + key)
	}
	return StatsReport{
		Country: sp[1],
		App: sp[2],
		Platform: sp[3],
	}
}

func (app *App) SaveRequest(ctx *fasthttp.RequestCtx) {
	defer ctx.SetConnectionClose()
	if time.Now().Sub(app.LastReq()) > 5*time.Second {
		app.IncrCurrPos()
	}
	app.SetLastReq(time.Now())

	var data JData
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		fmt.Println("Decode:",err.Error())
		ctx.Response.SetStatusCode(400)
		fmt.Fprint(ctx,err)
		return
	}

	_, err = app.RClient.Incr(data.Key()).Result()
	if err != nil {
		fmt.Println(err)
		ctx.Response.SetStatusCode(500)
		fmt.Fprint(ctx,err)
		return
	}

	ctx.SetStatusCode(200)
	fmt.Fprint(ctx, "{pos:", app.currPos, "}")
	/*

	_, err = app.RClient.Incr(data.Key()).Result()
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"err": err,
		})
		return
	}

	c.JSON(200, gin.H{
		"pos": app.CurrPos(),
	})
	*/
}

func (a *App)Stats(c *gin.Context) {
	keys, err := a.RClient.Keys("stat:*").Result()

	if err != nil {
		c.JSON(500, gin.H{
			"err" : err,
		})
		return
	}

	reports := make([]StatsReport,0,0)
	for _, k := range keys {
		report := createReport(k)
		count, err := a.RClient.Get(k).Result()
		if err != nil {
			fmt.Println( "Не удалось получить статистику! " ,err)
			continue
		}
		report.Count = count
		reports = append(reports, report)
	}

	c.JSON(200, reports)
	/*c.JSON(200, gin.H{
		"stats" : keys,
	})*/
}
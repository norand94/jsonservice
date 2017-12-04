package core

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"strings"
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

func (app *App) SaveRequest(c *gin.Context) {
	fmt.Println(c.Request.Header.Get("Content-Type"))
	if time.Now().Sub(app.LastReq()) > 5*time.Second {
		app.IncrCurrPos()
	}
	app.SetLastReq(time.Now())

	dec := json.NewDecoder(c.Request.Body)
	var data JData
	err := dec.Decode(&data)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	fmt.Printf("%+v \n", data)
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
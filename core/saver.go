package core

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
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
	return d.Device.Geo.Country + ":" + d.Device.OS + ":" + d.App.Bundle
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
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	fmt.Printf("%+v \n", data)
	cmd := app.RClient.Incr(data.Key())
	if cmd.Err != nil {
		c.JSON(500, gin.H{
			"err": cmd.Err(),
		})
		return
	}

	c.JSON(200, gin.H{
		"pos": app.CurrPos(),
	})
}

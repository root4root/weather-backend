package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"weather/common"
	"weather/openweathermap"

	"github.com/labstack/echo/v4"
)

var adata = common.Apidata{
	Phenomena: "N/A",
	Main:      "N/A",
}

func main() {
	confPath := flag.String("c", "weather.yml", "path to config file")
	logPath := flag.String("l", "weather.log", "path to log file")
	apiPort := flag.String("p", "8000", "listen port")

	flag.Parse()

	if flag.NFlag() < 3 {
		flag.PrintDefaults()
	}

	loggerInit(*logPath)

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	rapi := openweathermap.New(*confPath)

	go cron(ticker, rapi)

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/json", adata.GetJson())
	})

	log.Println("=== Weather backend API successfully started")

	e.Logger.Fatal(e.Start(":" + *apiPort))
}

type RemoteApi interface {
	PrepareData() (common.Apidata, error)
}

func cron(tck *time.Ticker, api RemoteApi) {
	for {
		newData, err := api.PrepareData()

		if err != nil {
			log.Println(err)
		} else {
			adata.SetData(newData)
		}

		<-tck.C
	}
}

func loggerInit(logPath string) {
	logFile, logFileErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if logFileErr != nil {
		fmt.Fprintf(os.Stderr, "Log file error! %v\n", logFileErr)
		os.Exit(1)
	}

	log.SetOutput(logFile)
}

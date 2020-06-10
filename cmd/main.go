package main

import (
	"log"
	"os"

	exporter "github.com/siangyeh8818/ssl.effective.date.exporter/internal"
	"github.com/siangyeh8818/ssl.effective.date.exporter/internal/database"
)

var (
	appCache exporter.PNCache
)

func main() {

	// connect redis, which is sit inside the kubernetes cluster
	database.ConnectRedis()

	//exporter.VerifySSL("ch02.hnjump.cn")

	var config = exporter.BaseConfig{}
	path := os.Getenv("CONFIG_PATH")
	(&config).Initconfig(path)
	//(&config).Initconfig("/opt/gaia/gaiaDomains.json")

	log.Println((&config).Domain)
	//config.Initdoman("test.json")

	appCache := exporter.Initcache()

	metrics := exporter.AddMetrics()
	exp := exporter.Exporter{
		SSLMetrics: metrics,
		Config:     (&config),
		Cache:      appCache,
	}

	go func() {
		(&exp).HandlerGatherData()
	}()

	exporter.NewServer(exp).Start()
	//exporter.Run_Exporter_Server()
}

package main

import (
	"log"

	exporter "github.com/siangyeh8818/ssl.effective.date.exporter/internal"
)

func main() {

	//exporter.VerifySSL("ch02.hnjump.cn")

	var config = exporter.BaseConfig{}
	(&config).Initconfig("/opt/gaia/gaiaDomains.json")

	log.Println((&config).Domain)
	//config.Initdoman("test.json")

	metrics := exporter.AddMetrics()
	exp := exporter.Exporter{
		SSLMetrics: metrics,
		Config:     (&config),
	}

	exporter.NewServer(exp).Start()
	//exporter.Run_Exporter_Server()
}
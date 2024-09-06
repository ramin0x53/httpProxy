package main

import (
	"httpProxy/config"
	"httpProxy/handler"
	"httpProxy/logger"
	"log"
	"net/http"
)

func main() {
	config := config.NewConfig()
	showBanner(config.Server.ListenAddr)
	loggerInstance := logger.NewLogger(config.Logger)
	handler := handler.NewDefaultHandler(config.Proxy, loggerInstance)

	if config.Server.CertPath != "" && config.Server.KeyPath != "" {
		err := http.ListenAndServeTLS(config.Server.ListenAddr, config.Server.CertPath, config.Server.KeyPath, handler)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := http.ListenAndServe(config.Server.ListenAddr, handler)
		if err != nil {
			log.Fatal(err)
		}
	}
}

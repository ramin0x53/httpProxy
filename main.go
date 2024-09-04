package main

import (
	"httpProxy/config"
	"httpProxy/handler"
	"httpProxy/logger"
	"net/http"
)

func main() {
	config := config.NewConfig()
	loggerInstance := logger.NewLogger(config.Logger)
	handler := handler.NewDefaultHandler(config.Proxy, loggerInstance)
	http.ListenAndServe(config.Server.ListenAddr, handler)
}

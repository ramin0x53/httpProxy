package main

import (
	"httpProxy/config"
	"httpProxy/handler"
	"net/http"
)

func main() {
	config := config.NewConfig()
	handler := handler.NewDefaultHandler(config.Proxy)
	http.ListenAndServe(config.Server.ListenAddr, handler)
}

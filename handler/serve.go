package handler

import (
	"httpProxy/config"
	"httpProxy/logger"
	"httpProxy/proxy"
	"net/http"
)

type DefaultHandler struct {
	proxyConfig *config.ProxyConfig
	logger      *logger.Logger
}

func NewDefaultHandler(cfg *config.ProxyConfig, logger *logger.Logger) *DefaultHandler {
	return &DefaultHandler{proxyConfig: cfg, logger: logger}
}

func (h *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxy := proxy.NewHttpProxy(h.proxyConfig.RemoteProtocol, h.proxyConfig.RemoteHost, h.proxyConfig.RemotePort, h.proxyConfig.PreventHostReplace, w, r)
	proxy.Redirect()
	h.logger.LogHttpRequest(proxy.GetProcessData())
}

package config

import "strings"

type Config struct {
	Logger *LoggerConfig
	Server *ServerConfig
	Proxy  *ProxyConfig
}

func NewConfig() *Config {
	cfg := &Config{}
	cfg.Logger = &LoggerConfig{}
	cfg.Server = &ServerConfig{}
	cfg.Proxy = &ProxyConfig{}
	cfg.get()
	return cfg
}

func (c *Config) get() {
	c.Proxy.RemoteProtocol = "https"
	c.Proxy.RemoteHost = "api.hamsterkombat.io"
	c.Proxy.RemotePort = "443"
	c.Proxy.PreventHostReplace = false
	c.Server.ListenAddr = "2020"

	if !strings.Contains(c.Server.ListenAddr, ":") {
		c.Server.ListenAddr = "0.0.0.0:" + c.Server.ListenAddr
	}
}

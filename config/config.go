package config

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
)

type keyValueFlag map[string]string

func (m *keyValueFlag) Set(value string) error {
	parts := strings.Split(value, "=")
	if len(parts) != 2 {
		return fmt.Errorf("invalid map item: %s", value)
	}
	key := parts[0]
	val := parts[1]
	(*m)[key] = val
	return nil
}

func (m *keyValueFlag) String() string {
	var pairs []string
	for k, v := range *m {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(pairs, ",")
}

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
	listenAddr := flag.String("l", "0.0.0.0:8080", "Listen address (default: 0.0.0.0:8080)")
	proxyAddr := flag.String("p", "", "Proxy address (e.g. https://example.com:3000)")
	preventHostReplace := flag.Bool("r", false, "Prevent to replace Host header (default: false)")
	curlEnable := flag.Bool("c", false, "Show request curl (default: false)")
	pathInclude := flag.String("fpi", "", "Filter paths that contains following string")
	reqBodyInclude := flag.String("frqbi", "", "Filter request bodies that contains following string")
	resBodyInclude := flag.String("frsbi", "", "Filter response bodies that contains following string")

	var reqHeaderInclude keyValueFlag = make(map[string]string)
	var resHeaderInclude keyValueFlag = make(map[string]string)
	flag.Var(&reqHeaderInclude, "frqhi", "Filter request header that contains following string (e.g. Referrer=test.com)")
	flag.Var(&resHeaderInclude, "frshi", "Filter response header that contains following string (e.g. Origin=test.com)")

	flag.Parse()

	c.Logger.Curl = *curlEnable
	c.Logger.PathInclude = *pathInclude
	c.Logger.ResBodyInclude = *resBodyInclude
	c.Logger.ReqBodyInclude = *reqBodyInclude
	c.Logger.ReqHeaderInclude = reqHeaderInclude
	c.Logger.ResHeaderInclude = resHeaderInclude

	c.Server.ListenAddr = *listenAddr

	pAddr, err := url.Parse(*proxyAddr)
	if err != nil {
		log.Fatal(err)
	}

	c.Proxy.PreventHostReplace = *preventHostReplace
	c.Proxy.RemoteProtocol = pAddr.Scheme
	c.Proxy.RemoteHost = pAddr.Hostname()
	c.Proxy.RemotePort = pAddr.Port()

	if !strings.Contains(c.Server.ListenAddr, ":") {
		c.Server.ListenAddr = "0.0.0.0:" + c.Server.ListenAddr
	}
}

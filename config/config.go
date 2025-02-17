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
	certPath := flag.String("cert", "", "Cert path")
	keyPath := flag.String("key", "", "Key path")
	proxyAddr := flag.String("p", "", "Proxy address (e.g. https://example.com:3000)")
	preventHostReplace := flag.Bool("r", false, "Prevent to replace Host header (default: false)")
	curlEnable := flag.Bool("c", false, "Show request curl (default: false)")
	pathInclude := flag.String("fpi", "", "Filter paths that contains following string")
	reqBodyInclude := flag.String("frqbi", "", "Filter request bodies that contains following string")
	resBodyInclude := flag.String("frsbi", "", "Filter response bodies that contains following string")

	statuCodeInclude := flag.Int("fsi", 0, "Filter status code that contains following number")
	statuCodeExclude := flag.Int("fse", 0, "Filter status code that not contains following number")

	var reqHeaderInclude keyValueFlag = make(map[string]string)
	var resHeaderInclude keyValueFlag = make(map[string]string)
	flag.Var(&reqHeaderInclude, "frqhi", "Filter request header that contains following string (e.g. Referrer=test.com)")
	flag.Var(&resHeaderInclude, "frshi", "Filter response header that contains following string (e.g. Origin=test.com)")

	pathExclude := flag.String("fpe", "", "Filter paths that not contains following string")
	reqBodyExclude := flag.String("frqbe", "", "Filter request bodies that not contains following string")
	resBodyExclude := flag.String("frsbe", "", "Filter response bodies that not contains following string")

	var reqHeaderExclude keyValueFlag = make(map[string]string)
	var resHeaderExclude keyValueFlag = make(map[string]string)
	flag.Var(&reqHeaderExclude, "frqhe", "Filter request header that not contains following string (e.g. Referrer=test.com)")
	flag.Var(&resHeaderExclude, "frshe", "Filter response header that not contains following string (e.g. Origin=test.com)")

	flag.Parse()

	c.Logger.Curl = *curlEnable
	c.Logger.PathInclude = *pathInclude
	c.Logger.ResBodyInclude = *resBodyInclude
	c.Logger.ReqBodyInclude = *reqBodyInclude
	c.Logger.ReqHeaderInclude = reqHeaderInclude
	c.Logger.ResHeaderInclude = resHeaderInclude

	c.Logger.PathExclude = *pathExclude
	c.Logger.ResBodyExclude = *resBodyExclude
	c.Logger.ReqBodyExclude = *reqBodyExclude
	c.Logger.ReqHeaderExclude = reqHeaderExclude
	c.Logger.ResHeaderExclude = resHeaderExclude

	c.Server.ListenAddr = *listenAddr
	c.Server.CertPath = *certPath
	c.Server.KeyPath = *keyPath

	c.Logger.StatusCodeInclude = *statuCodeInclude
	c.Logger.StatusCodeExclude = *statuCodeExclude

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

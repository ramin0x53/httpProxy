package config

type ProxyConfig struct {
	RemoteProtocol     string
	RemoteHost         string
	RemotePort         string
	PreventHostReplace bool
}

package config

type LoggerConfig struct {
	Curl             bool
	PathInclude      string
	ReqBodyInclude   string
	ReqHeaderInclude map[string]string
	ResBodyInclude   string
	ResHeaderInclude map[string]string
}

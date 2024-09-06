package config

type LoggerConfig struct {
	Curl              bool
	PathInclude       string
	ReqBodyInclude    string
	ReqHeaderInclude  map[string]string
	ResBodyInclude    string
	ResHeaderInclude  map[string]string
	PathExclude       string
	ReqBodyExclude    string
	ReqHeaderExclude  map[string]string
	ResBodyExclude    string
	ResHeaderExclude  map[string]string
	StatusCodeInclude int
	StatusCodeExclude int
}

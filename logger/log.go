package logger

import (
	"httpProxy/config"
	"strings"
)

type httpLog struct {
	reqProtocol string
	resProtocol string
	method      string
	host        string
	statusCode  int
	StatusText  string
	reqHeader   map[string][]string
	resHeader   map[string][]string
	curl        string
	reqBodyStr  string
	resBodyStr  string
	path        string
	err         error
}

type HttpInfo interface {
	ReqProtocol() string
	ResProtocol() string
	Method() string
	Host() string
	StatusCode() (int, string)
	ReqHeader() map[string][]string
	ResHeader() map[string][]string
	CURL() (string, error)
	ReqBodyStr() string
	ResBodyStr() string
	Path() string
	GetError() error
}

type Logger struct {
	Config *config.LoggerConfig
	queue  chan *httpLog
}

func NewLogger(cfg *config.LoggerConfig) *Logger {
	logger := &Logger{Config: cfg}
	logger.queue = make(chan *httpLog, 1000)
	go logger.logWorker()
	return logger
}

func (l *Logger) logWorker() {
	for log := range l.queue {
		httpTerminalPrint(log)
	}
}

func (l *Logger) LogHttpRequest(data HttpInfo) {
	statusCode, statusTxt := data.StatusCode()

	log := &httpLog{reqProtocol: data.ReqProtocol(), resProtocol: data.ResProtocol(), method: data.Method(), host: data.Host(), statusCode: statusCode, StatusText: statusTxt, reqHeader: data.ReqHeader(), resHeader: data.ResHeader(), path: data.Path(), err: data.GetError(), reqBodyStr: data.ReqBodyStr(), resBodyStr: data.ResBodyStr()}

	if l.Config.Curl {
		var err error
		log.curl, err = data.CURL()
		if err != nil {
			log.curl = err.Error()
		}
	} else {
		log.curl = ""
	}

	go l.filter(log)
}

func (l *Logger) filter(log *httpLog) {
	if l.Config.ReqBodyInclude != "" && !strings.Contains(log.reqBodyStr, l.Config.ReqBodyInclude) {
		return
	}

	if l.Config.PathInclude != "" && !strings.Contains(log.path, l.Config.PathInclude) {
		return
	}

	if l.Config.ResBodyInclude != "" && !strings.Contains(log.resBodyStr, l.Config.ResBodyInclude) {
		return
	}

	for key, value := range l.Config.ReqHeaderInclude {
		for _, headerValue := range log.reqHeader[key] {
			if value != "" && !strings.Contains(headerValue, value) {
				return
			}
		}
	}

	for key, value := range l.Config.ResHeaderInclude {
		for _, headerValue := range log.resHeader[key] {
			if value != "" && !strings.Contains(headerValue, value) {
				return
			}
		}
	}

	l.queue <- log
}

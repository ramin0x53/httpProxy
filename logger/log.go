package logger

import (
	"httpProxy/config"
	"httpProxy/utility"
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
	if checkRes, exist := utility.IncludeCheck(log.reqBodyStr, l.Config.ReqBodyInclude); exist && !checkRes {
		return
	}

	if checkRes, exist := utility.IncludeCheck(log.path, l.Config.PathInclude); exist && !checkRes {
		return
	}

	if checkRes, exist := utility.IncludeCheck(log.resBodyStr, l.Config.ResBodyInclude); exist && !checkRes {
		return
	}

	for key, value := range l.Config.ReqHeaderInclude {
		for _, headerValue := range log.reqHeader[key] {
			if checkRes, exist := utility.IncludeCheck(headerValue, value); exist && !checkRes {
				return
			}
		}
	}

	for key, value := range l.Config.ResHeaderInclude {
		for _, headerValue := range log.resHeader[key] {
			if checkRes, exist := utility.IncludeCheck(headerValue, value); exist && !checkRes {
				return
			}
		}
	}

	if checkRes, exist := utility.ExcludeCheck(log.reqBodyStr, l.Config.ReqBodyExclude); exist && !checkRes {
		return
	}

	if checkRes, exist := utility.ExcludeCheck(log.path, l.Config.PathExclude); exist && !checkRes {
		return
	}

	if checkRes, exist := utility.ExcludeCheck(log.resBodyStr, l.Config.ResBodyExclude); exist && !checkRes {
		return
	}

	for key, value := range l.Config.ReqHeaderExclude {
		for _, headerValue := range log.reqHeader[key] {
			if checkRes, exist := utility.ExcludeCheck(headerValue, value); exist && !checkRes {
				return
			}
		}
	}

	for key, value := range l.Config.ResHeaderExclude {
		for _, headerValue := range log.resHeader[key] {
			if checkRes, exist := utility.ExcludeCheck(headerValue, value); exist && !checkRes {
				return
			}
		}
	}

	l.queue <- log
}

package proxy

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

type HttpData struct {
	ReqBody        *bytes.Buffer
	ResBody        *bytes.Buffer
	TargetRequest  *http.Request
	TargetResponse *http.Response
	Error          error
}

type HttpProxy struct {
	Protocol           string
	Host               string
	Port               string
	PreventHostReplace bool
	Request            *http.Request
	ResponseWriter     http.ResponseWriter
	processData        *HttpData
}

// TODO: check for query string implementation

func NewHttpProxy(protocol, host, port string, preventHostReplace bool, w http.ResponseWriter, r *http.Request) *HttpProxy {
	proxy := &HttpProxy{Protocol: protocol, Host: host, Port: port, PreventHostReplace: preventHostReplace, Request: r, ResponseWriter: w}
	return proxy
}

func (p *HttpProxy) Redirect() {
	res, err := p.proxyRequest()
	if err != nil {
		p.processData.Error = err
		http.Error(p.ResponseWriter, "Error while proxying the request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	err = p.proxyResponse(res)
	if err != nil {
		p.processData.Error = err
		http.Error(p.ResponseWriter, "Error while proxying the response", http.StatusInternalServerError)
		return
	}
}

func (p *HttpProxy) proxyResponse(res *http.Response) error {
	for header, values := range res.Header {
		for _, value := range values {
			p.ResponseWriter.Header().Add(header, value)
		}
	}

	p.ResponseWriter.WriteHeader(res.StatusCode)

	teeRes, resBody := p.TeeReader(res.Body)
	p.processData.ResBody = resBody

	if _, err := io.Copy(p.ResponseWriter, teeRes); err != nil {
		return err
	}
	return nil
}

func (p *HttpProxy) proxyRequest() (*http.Response, error) {
	remoteUrl := p.GetRemoteUrl()

	client := &http.Client{}

	teeReq, reqBody := p.TeeReader(p.Request.Body)
	p.processData.ReqBody = reqBody

	req, err := http.NewRequest(p.Request.Method, remoteUrl, teeReq)
	p.processData.TargetRequest = req
	if err != nil {
		return nil, err
	}

	req.Host = p.changeReqHost(req.Host)

	for header, values := range p.Request.Header {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}

	res, err := client.Do(req)
	p.processData.TargetResponse = res
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *HttpProxy) changeReqHost(host string) string {
	if p.PreventHostReplace {
		inputHost := ""
		if strings.Contains(p.Request.Host, ":") {
			inputHost = strings.Split(p.Request.Host, ":")[0]
		} else {
			inputHost = p.Request.Host
		}

		outputPort := ""
		if strings.Contains(host, ":") {
			outputPort = strings.Split(host, ":")[1]
		}

		return inputHost + outputPort
	} else {
		return host
	}
}

func (p *HttpProxy) GetRemoteUrl() string {
	url := p.Protocol + "://" + p.Host + ":" + p.Port + p.Request.URL.Path
	return url
}

func (p *HttpProxy) TeeReader(r io.Reader) (io.Reader, *bytes.Buffer) {
	var buffer bytes.Buffer
	return io.TeeReader(r, &buffer), &buffer
}

func (p *HttpProxy) GetProcessData() *HttpData {
	return p.processData
}

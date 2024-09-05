package proxy

import (
	"bytes"
	"errors"
	"httpProxy/utility"
	"io"
	"net/http"
	"strings"

	"moul.io/http2curl"
)

type ReadCloser struct{ *bytes.Buffer }

func (r ReadCloser) Close() error {
	return nil
}

type HttpData struct {
	ReqBody        *bytes.Buffer
	ResBody        *bytes.Buffer
	TargetRequest  *http.Request
	TargetResponse *http.Response
	Error          error
}

func (d *HttpData) ReqProtocol() string {
	if d.TargetRequest == nil {
		return ""
	}

	return d.TargetRequest.Proto
}

func (d *HttpData) ResProtocol() string {
	if d.TargetResponse == nil {
		return ""
	}

	return d.TargetResponse.Proto
}

func (d *HttpData) Method() string {
	if d.TargetRequest == nil {
		return ""
	}

	return d.TargetRequest.Method
}

func (d *HttpData) Host() string {
	if d.TargetRequest == nil {
		return ""
	}

	return d.TargetRequest.Host
}

func (d *HttpData) StatusCode() (int, string) {
	if d.TargetResponse == nil {
		return 0, ""
	}

	statusCode := d.TargetResponse.StatusCode
	return statusCode, http.StatusText(statusCode)
}

func (d *HttpData) ReqHeader() map[string][]string {
	if d.TargetRequest == nil {
		return nil
	}

	return d.TargetRequest.Header
}

func (d *HttpData) ResHeader() map[string][]string {
	if d.TargetResponse == nil {
		return nil
	}

	return d.TargetResponse.Header
}

func (d *HttpData) CURL() (string, error) {
	if d.TargetRequest == nil {
		return "", nil
	}

	readCloser := ReadCloser{d.ReqBody}
	d.TargetRequest.Body = readCloser

	command, err := http2curl.GetCurlCommand(d.TargetRequest)
	if err != nil {
		return "", err
	}
	return command.String(), nil
}

func (d *HttpData) ReqBodyStr() string {
	if d.ReqBody == nil {
		return ""
	}

	if values, exists := d.ReqHeader()["Content-Encoding"]; exists {
		body, err := utility.DecodeContent(strings.Trim(values[0], " "), d.ReqBody)
		if err != nil {
			return err.Error()
		}
		return string(*body)
	}

	return d.ReqBody.String()
}

func (d *HttpData) ResBodyStr() string {
	if d.ResBody == nil {
		return ""
	}

	if values, exists := d.ResHeader()["Content-Encoding"]; exists {
		body, err := utility.DecodeContent(strings.Trim(values[0], " "), d.ResBody)
		if err != nil {
			return err.Error()
		}
		return string(*body)
	}

	return d.ResBody.String()
}

func (d *HttpData) Path() string {
	if d.TargetRequest == nil {
		return ""
	}

	path := d.TargetRequest.URL.Path
	if d.TargetRequest.URL.RawQuery != "" {
		path = path + "?" + d.TargetRequest.URL.RawQuery
	}
	return path
}

func (d *HttpData) GetError() error {
	return d.Error
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

func NewHttpProxy(protocol, host, port string, preventHostReplace bool, w http.ResponseWriter, r *http.Request) *HttpProxy {
	proxy := &HttpProxy{Protocol: protocol, Host: host, Port: port, PreventHostReplace: preventHostReplace, Request: r, ResponseWriter: w, processData: &HttpData{}}
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
	if p.Host == "" {
		return p.handleEmptyHost()
	}

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

		return inputHost + ":" + outputPort
	} else {
		return host
	}
}

func (p *HttpProxy) handleEmptyHost() (*http.Response, error) {
	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, p.Request.Body)
	if err != nil {
		return nil, err
	}

	p.processData.ReqBody = &buffer
	p.Request.URL.Host = p.Request.Host
	p.Request.URL.Scheme = "http"
	p.processData.TargetRequest = p.Request
	return nil, errors.New("Host is empty")
}

func (p *HttpProxy) GetRemoteUrl() string {
	url := p.Protocol + "://" + p.Host + ":" + p.Port + p.Request.URL.Path
	if p.Request.URL.RawQuery != "" {
		url = url + "?" + p.Request.URL.RawQuery
	}
	return url
}

func (p *HttpProxy) TeeReader(r io.Reader) (io.Reader, *bytes.Buffer) {
	var buffer bytes.Buffer
	return io.TeeReader(r, &buffer), &buffer
}

func (p *HttpProxy) GetProcessData() *HttpData {
	return p.processData
}

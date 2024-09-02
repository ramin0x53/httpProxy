package handler

import (
	"bytes"
	"fmt"
	"httpProxy/config"
	"io"
	"log"
	"net/http"
	"strings"
)

type DefaultHandler struct {
	proxyConfig *config.ProxyConfig
}

func NewDefaultHandler(cfg *config.ProxyConfig) *DefaultHandler {
	return &DefaultHandler{proxyConfig: cfg}
}

func (h *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remoteUrl := h.proxyConfig.RemoteProtocol + "://" + h.proxyConfig.RemoteHost + ":" + h.proxyConfig.RemotePort
	remoteUrl = remoteUrl + r.URL.Path

	client := &http.Client{}

	var reqBody bytes.Buffer
	teeReq := io.TeeReader(r.Body, &reqBody)
	req, err := http.NewRequest(r.Method, remoteUrl, teeReq)
	if err != nil {
		log.Panicln(err)
		return
	}

	if h.proxyConfig.PreventHostReplace {
		inputHost := ""
		if strings.Contains(r.Host, ":") {
			inputHost = strings.Split(r.Host, ":")[0]
		} else {
			inputHost = r.Host
		}

		outputPort := ""
		if strings.Contains(req.Host, ":") {
			outputPort = strings.Split(req.Host, ":")[1]
		}

		req.Host = inputHost + outputPort
	}

	for header, values := range r.Header {
		for _, value := range values {
			req.Header.Add(header, value)
		}
		// req.Header.Add(header, value[0])
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	for header, values := range res.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
		// w.Header().Set(header, value[0])
	}

	w.WriteHeader(res.StatusCode)

	var resBody bytes.Buffer
	teeRes := io.TeeReader(res.Body, &resBody)

	if _, err := io.Copy(w, teeRes); err != nil {
		http.Error(w, "Error reading stream", http.StatusInternalServerError)
		return
	}

	err = logPrepare(&reqBody, &resBody, req.Header, w.Header(), r.Method, r.URL.Path, r.Proto, req.Host, res.StatusCode)
	if err != nil {
		log.Println("Error while logging request: ", err)
		return
	}
}

func logPrepare(reqBody, resBody *bytes.Buffer, reqHeader, resHeader map[string][]string, method, path, Proto, reqHost string, statusCode int) error {
	fmt.Println(reqHeader)
	fmt.Println("++++++++++++++++++++++++++")
	fmt.Println(resHeader)

	// a, err := io.ReadAll(reqBody)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(string(a))

	// b, err := io.ReadAll(resBody)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(string(b))
	return nil
}

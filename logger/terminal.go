package logger

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func httpTerminalPrint(log *httpLog) {
	color.HiWhite(fmt.Sprintf("-----------------------------| %s |-----------------------------", time.Now().Format("2006-01-02 15:04:05.000000000")))
	fmt.Printf("\n")

	color.White("Request:")
	color.Green(makeHttpRequestLog(log.method, log.path, log.reqProtocol, log.host, log.reqBodyStr, log.reqHeader))
	fmt.Printf("\n\n")

	color.White("Response:")
	color.Cyan(makeHttpResponseLog(log.statusCode, log.StatusText, log.resProtocol, log.resBodyStr, log.resHeader))
	fmt.Printf("\n\n")

	if log.curl != "" {
		color.White("CURL:")
		color.Yellow(log.curl)
		fmt.Printf("\n\n")
	}

	if log.err != nil {
		color.Red("Error:")
		color.Red(log.curl)
		fmt.Printf("\n\n")
	}
}

func makeHttpResponseLog(statusCode int, statusTxt, protocol, body string, header map[string][]string) string {
	headerStr := headerToString(header)
	return fmt.Sprintf("%s %d %s\n%s\n%s", protocol, statusCode, statusTxt, headerStr, body)
}

func makeHttpRequestLog(method, path, protocol, host, body string, header map[string][]string) string {
	headerStr := headerToString(header)
	return fmt.Sprintf("%s %s %s\nHost: %s\n%s\n%s", method, path, protocol, host, headerStr, body)
}

func headerToString(header map[string][]string) string {
	txt := ""
	for key, values := range header {
		for _, value := range values {
			txt = txt + key + ": " + value + "\n"
		}
	}
	return txt
}

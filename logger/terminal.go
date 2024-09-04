package logger

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func httpTerminalPrint(log *httpLog) {
	color.Magenta("-----------------------", time.Now().String(), "-----------------------", "\n")

	color.White("Request:")
	color.Green(makeHttpRequestLog(log.method, log.path, log.reqProtocol, log.host, log.reqBodyStr, log.reqHeader))
	fmt.Printf("\n\n\n")

	color.White("Response:")
	color.Cyan(makeHttpResponseLog(log.statusCode, log.StatusText, log.resProtocol, log.resBodyStr, log.resHeader))
	fmt.Printf("\n\n\n")

	if log.curl != "" {
		color.White("CURL:")
		color.Yellow(log.curl)
		fmt.Printf("\n\n\n")
	}

	if log.err != nil {
		color.Red("Error:")
		color.Red(log.curl)
		fmt.Printf("\n\n\n")
	}
}

func makeHttpResponseLog(statusCode int, statusTxt, protocol, body string, header map[string][]string) string {
	headerStr := headerToString(header)
	return fmt.Sprintf("%s %d %s\n%s\n\n%s", protocol, statusCode, statusTxt, headerStr, body)
}

func makeHttpRequestLog(method, path, protocol, host, body string, header map[string][]string) string {
	headerStr := headerToString(header)
	return fmt.Sprintf("%s %s %s\nHost: %s\n%s\n\n%s", method, path, protocol, host, headerStr, body)
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

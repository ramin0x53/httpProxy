package main

import (
	"fmt"

	"github.com/dimiro1/banner"
	"github.com/mattn/go-colorable"
)

func showBanner(listenAddr string) {
	templ := fmt.Sprintf(`{{ .AnsiColor.BrightGreen }}{{ .Title "httpProxy" "" 8 }}
   {{ .AnsiColor.BrightCyan }}Author: {{ .AnsiColor.Yellow }}Ramin Sardari{{ .AnsiColor.Default }}
   {{ .AnsiColor.BrightCyan }}ListenAddr: {{ .AnsiColor.Yellow }}%s{{ .AnsiColor.Default }}
   {{ .AnsiColor.BrightCyan }}Now: {{ .AnsiColor.Yellow }}{{ .Now "Mon, 02 Jan 2006 15:04:05 -0700" }}{{ .AnsiColor.Default }}`+"\n\n", listenAddr)

	banner.InitString(colorable.NewColorableStdout(), true, true, templ)
}

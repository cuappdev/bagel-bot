package main

import (
	"github.com/op/go-logging"
	"os"
)

var log = func() *logging.Logger {
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
	return logging.MustGetLogger("example")
}()

func main() {
	go SlackListenAndServe()
	BagelRepl()
}

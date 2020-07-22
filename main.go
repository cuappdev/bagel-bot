package main

import (
    "github.com/op/go-logging"
    "os"
)

var log = logging.MustGetLogger("example")

func initializeLog() {
    format := logging.MustStringFormatter(
        `%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
    )
    backend := logging.NewLogBackend(os.Stderr, "", 0)
    backendFormatter := logging.NewBackendFormatter(backend, format)
    logging.SetBackend(backendFormatter)
}

func main() {
    initializeLog()
    d := Database{Filename: "data.db"}
    d.initialize()

    slackApiKey := os.Getenv("BAGEL_SLACK_API_KEY")
    if len(slackApiKey) == 0 {
        log.Warning("Slack API key is missing; lots of errors ahead")
    }
}

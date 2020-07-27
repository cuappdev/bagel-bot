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
	db := OpenDB("data.db")
	MigrateDB(db)

	if SlackApiKey == "" {
		log.Warning("Slack API key missing; prepare for lots of errors")
	}
	s := &Slack{Token: SlackApiKey}

	go SlackListenAndServe(db, s)
	BagelRepl(db, s)

	log.Debug("Closing connection to DB")
	if err := db.Close(); err != nil {
		panic(err)
	}
}

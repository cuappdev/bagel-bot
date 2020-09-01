package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/op/go-logging"
	"gorm.io/gorm"
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
	var MainCLI struct {
		Repl    bool     `help:"Start repl"`
		LocalDB bool     `help:"Use local sqlite3 db instead of connecting with postgres"`
		Args    []string `help:"Arguments to bagel. Ignored if --repl is used"`
	}
	_ = kong.Parse(&MainCLI)

	var db *gorm.DB
	if MainCLI.LocalDB {
		db = OpenSqlite3DBLocal("data.db")
	} else {
		var err error
		db, err = OpenPostgresDB(PostgresHost, PostgresPort, PostgresUser, PostgresPassword, PostgresDbName)
		if err != nil {
			log.Error(err)
		}
	}

	err := MigrateDB(db)
	if err != nil {
		log.Error(err)
		return
	}

	if SlackApiKey == "" {
		log.Warning("Slack API key missing; prepare for lots of errors")
	}
	s := &Slack{Token: SlackApiKey}

	if MainCLI.Repl {
		go SlackListenAndServe(db, s)
		BagelRepl(db, s)
	} else {
		exited := false
		cli := &CLI{}
		parser, err := kong.New(cli, kong.Exit(func(i int) {
			exited = true
		}))
		if err != nil {
			log.Error(err)
			return
		}

		ctx, err := parser.Parse(MainCLI.Args)
		if exited {
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		if err = ctx.Run(); err != nil {
			fmt.Println(err)
		}
	}
}

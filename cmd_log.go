package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/jinzhu/gorm"
	"io"
	"strings"
	"time"
)

type CmdLog struct {
	Last  int    `help:"Print the n-last log entries" default:"10"`
	LogID string `arg optional help:"Print the bagels associated with a log. <logid> can either be \"last\" or the ID of a log."`
}

func (cmd *CmdLog) Run(ctx *kong.Context, db *gorm.DB) (err error) {
	switch ctx.Command() {
	case "log":
		if cmd.Last < 1 {
			_, err = io.WriteString(ctx.Stderr, "--last must be a positive integer")
			return err
		}

		msg := fmt.Sprintf("| %5s | %s: %s\n", "id", "date", "invocation")
		if _, err = io.WriteString(ctx.Stdout, msg); err != nil {
			return err
		}

		var logs []BagelLog
		db.Order("date desc").Limit(cmd.Last).Find(&logs)
		for _, log := range logs {
			t := time.Unix(log.Date, 0)
			msg := fmt.Sprintf("| %5d | %s: %s\n", log.ID, t.Format(time.UnixDate), log.Invocation)
			if _, err = io.WriteString(ctx.Stdout, msg); err != nil {
				return err
			}
		}

	case "log <log-id>":
		log, err := BagelLog_Fetch(db, cmd.LogID)
		if err != nil {
			return err
		}
		if log.ID == 0 {
			_, err = io.WriteString(ctx.Stderr, "no such log found\n")
			return err
		}

		msg := fmt.Sprintf("invoked with: %s\n", log.Invocation)
		if _, err = io.WriteString(ctx.Stdout, msg); err != nil {
			return err
		}

		msg = fmt.Sprintf("| %5s | %s: %s\n", "id", "slack id", "usernames")
		if _, err = io.WriteString(ctx.Stdout, msg); err != nil {
			return err
		}

		var bagels []Bagel
		db.Model(&log).Association("Bagels").Find(&bagels)
		for _, bagel := range bagels {
			var users []User
			db.Model(&bagel).Association("Users").Find(&users)

			var usernames []string
			for _, user := range users {
				usernames = append(usernames, user.Name)
			}

			msg := fmt.Sprintf("| %5d | %s: %s\n", bagel.ID, bagel.SlackConversationID, strings.Join(usernames, ", "))
			if _, err = io.WriteString(ctx.Stdout, msg); err != nil {
				return err
			}
		}
	default:
		panic(ctx.Command())

	}

	return nil
}

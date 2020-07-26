package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/jinzhu/gorm"
	"io"
)

type CmdMsg struct {
	Send CmdMsgSend `cmd`
	Read CmdMsgRead `cmd`
}

type CmdMsgSend struct {
	Text string `required help:"The text to send"`
	Log  struct {
		ID string `arg required help:"The id of the log to send a message to"`
	} `cmd help:"Send messages to a set of bagels previously created"`
	Channel struct {
		SlackID string `arg required help:"The id of the slack channel/group to send a message to"`
	} `cmd help:"Send messages to a specific slack channel"`
}

func (cmd *CmdMsgSend) Run(ctx *kong.Context, db *gorm.DB, s *Slack) (err error) {
	switch ctx.Command() {
	case "msg send log <id>":
		log, err := BagelLog_Fetch(db, cmd.Log.ID)
		if err != nil {
			return err
		}

		var bagels []Bagel
		db.Model(&log).Association("Bagels").Find(&bagels)
		for _, bagel := range bagels {
			conversation := bagel.SlackConversationID
			if err = s.ChatPostMessage(conversation, cmd.Text); err != nil {
				_, err = io.WriteString(ctx.Stderr, err.Error())
				return err
			}
		}
		return nil

	case "msg send channel <slack-id>":
		if err = s.ChatPostMessage(cmd.Channel.SlackID, cmd.Text); err != nil {
			_, err = io.WriteString(ctx.Stderr, err.Error())
			return err
		}
		return nil

	default:
		panic(ctx.Command())
	}
}

type CmdMsgRead struct {
	Limit   int `help:"The maximum number of messages to print" default:"10"`
	Channel struct {
		SlackID string `arg required help:"The id of the slack channel to read messages"`
	} `cmd help:"Read messages from a specific slack channel"`
}

func (cmd *CmdMsgRead) Run(ctx *kong.Context, db *gorm.DB, s *Slack) (err error) {
	switch ctx.Command() {
	case "msg read channel <slack-id>":
		messages, err := s.ConversationsHistory(cmd.Channel.SlackID, cmd.Limit)
		if err != nil {
			_, err = io.WriteString(ctx.Stderr, err.Error())
			return err
		}

		for _, slackMessage := range messages {
			var user User
			db.Where("slack_id = ?", slackMessage.User).FirstOrInit(&user)

			msg := fmt.Sprintf("%s: %s\n", user.Name, slackMessage.Text)
			if _, err = io.WriteString(ctx.Stdout, msg); err != nil {
				return err
			}
		}

		return nil

	default:
		panic(ctx.Command())
	}
}

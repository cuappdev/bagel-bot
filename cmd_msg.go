package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/jinzhu/gorm"
	"io"
	"strconv"
)

type CmdMsg struct {
	Send     CmdMsgSend     `cmd`
	Read     CmdMsgRead     `cmd`
	Feedback CmdMsgFeedback `cmd`
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
			if err = s.ChatPostMessage(conversation, cmd.Text, nil); err != nil {
				_, err = io.WriteString(ctx.Stderr, err.Error())
				return err
			}
		}
		return nil

	case "msg send channel <slack-id>":
		if err = s.ChatPostMessage(cmd.Channel.SlackID, cmd.Text, nil); err != nil {
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

type CmdMsgFeedback struct {
	Log struct {
		ID string `arg required help:"The id of the log to ask feedback from"`
	} `cmd help:"Send feedback request to a set of bagels previously created"`
	Bagel struct {
		ID int `arg required help:"The id of the bagel to ask feedback from"`
	} `cmd help:"Send feedback request to a specific bagel"`
}

func (cmd *CmdMsgFeedback) Run(ctx *kong.Context, db *gorm.DB, s *Slack) (err error) {
	text := "Could I get an update on the status of your bagel chat?"

	switch ctx.Command() {
	case "msg feedback log <id>":
		log, err := BagelLog_Fetch(db, cmd.Log.ID)
		if err != nil {
			return err
		}

		var bagels []Bagel
		db.Model(&log).Association("Bagels").Find(&bagels)
		for _, bagel := range bagels {
			feedbackMsg := FeedbackMsg_Create(db)
			db.Model(&bagel).Association("FeedbackMsgs").Append(&feedbackMsg)
			blocks := feedbackMsg.SlackBlocks("")

			conversation := bagel.SlackConversationID
			if err = s.ChatPostMessage(conversation, text, blocks); err != nil {
				_, err = io.WriteString(ctx.Stderr, err.Error())
				return err
			}
		}
		return nil

	case "msg feedback bagel <id>":
		var bagel Bagel
		db.Where("id = ?", cmd.Bagel.ID).First(&bagel)
		if bagel.ID == 0 {
			_, err = io.WriteString(ctx.Stderr, "no such bagel "+strconv.Itoa(cmd.Bagel.ID))
			return err
		}

		feedbackMsg := FeedbackMsg_Create(db)
		db.Model(&bagel).Association("FeedbackMsgs").Append(&feedbackMsg)
		blocks := feedbackMsg.SlackBlocks("")

		conversation := bagel.SlackConversationID
		if err = s.ChatPostMessage(conversation, text, blocks); err != nil {
			_, err = io.WriteString(ctx.Stderr, err.Error())
			return err
		}
		return nil

	default:
		panic(ctx.Command())
	}
}

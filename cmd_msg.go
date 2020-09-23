package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"gorm.io/gorm"
	"io"
	"strconv"
)

type CmdMsg struct {
	Send     CmdMsgSend     `cmd help:"Send a text message to a bagel"`
	Read     CmdMsgRead     `cmd help:"Read messages from bagels"`
	Feedback CmdMsgFeedback `cmd help:"Prompt bagels for their meeting status"`
	Stats    CmdMsgStats    `cmd help:"Post stats update of a bagel log"`
}

type CmdMsgSend struct {
	Text    string `help:"The text to send"`
	Default string `help:"The default message type to send"`
	Log     string `help:"The id of the log to send a message to"`
	Channel string `help:"The id of the slack channel/group to send a message to"`
}

func resolveDefault(def string) (message string) {
	switch def{
	case "reminder":
		message = "Hello! Just a reminder to find a time to catch up :upside_down_face:. We will be organizing new chats this Sunday, so better find a time asap!"
	}
	return message
}

func (cmd *CmdMsgSend) Run(ctx *kong.Context, db *gorm.DB, s *Slack) (err error) {
	var text string
	if cmd.Text != "" {
		text = cmd.Text
	} else {
		text = resolveDefault(cmd.Default)
		if text == "" {
			_, err = io.WriteString(ctx.Stderr, "No text or default message provided\n")
			return err
		}
	}

	if cmd.Log != "" {
		log, err := BagelLog_Fetch(db, cmd.Log)
		if err != nil {
			return err
		}

		var bagels []Bagel
		db.Model(&log).Association("Bagels").Find(&bagels)
		for _, bagel := range bagels {
			if err = s.ChatPostMessage(bagel.SlackConversationID, text, nil); err != nil {
				_, err = io.WriteString(ctx.Stderr, err.Error())
				return err
			}
			return nil
		}
	}

	if cmd.Channel != "" {
		if err = s.ChatPostMessage(cmd.Channel, text, nil); err != nil {
			_, err = io.WriteString(ctx.Stderr, err.Error())
			return err
		}
		return nil
	}

	return nil
}

type CmdMsgRead struct {
	Limit   int    `help:"The maximum number of messages to print" default:"10"`
	Channel string `required help:"The id of the slack channel to read messages"`
}

func (cmd *CmdMsgRead) Run(ctx *kong.Context, db *gorm.DB, s *Slack) (err error) {
	messages, err := s.ConversationsHistory(cmd.Channel, cmd.Limit)
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
}

type CmdMsgFeedback struct {
	Log   string `help:"The id of the log to ask feedback from"`
	Bagel int    `help:"The id of the bagel to ask feedback from"`
}

func (cmd *CmdMsgFeedback) Run(ctx *kong.Context, db *gorm.DB, s *Slack) (err error) {
	text := "I'm here to keep track of your bagel chat. When you've planned/completed your bagel chat, don't forget to mark that you've done so."

	if cmd.Log != "" {
		log, err := BagelLog_Fetch(db, cmd.Log)
		if err != nil {
			_, err = io.WriteString(ctx.Stderr, err.Error())
			return err
		}

		var bagels []Bagel
		db.Model(&log).Association("Bagels").Find(&bagels)
		for _, bagel := range bagels {
			feedbackMsg := FeedbackMsg_Create(db)
			db.Model(&bagel).Association("FeedbackMsgs").Append(&feedbackMsg)
			blocks := SlackBlocks_FeedbackMsg(feedbackMsg, text)

			conversation := bagel.SlackConversationID
			if err = s.ChatPostMessage(conversation, text, blocks); err != nil {
				_, err = io.WriteString(ctx.Stderr, err.Error())
				return err
			}
		}
	}

	if cmd.Bagel != 0 {
		var bagel Bagel
		db.Where("id = ?", cmd.Bagel).First(&bagel)
		if bagel.ID == 0 {
			_, err = io.WriteString(ctx.Stderr, "no such bagel "+strconv.Itoa(cmd.Bagel))
			return err
		}

		feedbackMsg := FeedbackMsg_Create(db)
		db.Model(&bagel).Association("FeedbackMsgs").Append(&feedbackMsg)
		blocks := SlackBlocks_FeedbackMsg(feedbackMsg, text)

		conversation := bagel.SlackConversationID
		if err = s.ChatPostMessage(conversation, text, blocks); err != nil {
			_, err = io.WriteString(ctx.Stderr, err.Error())
			return err
		}
	}

	return nil
}

type CmdMsgStats struct {
	Log     string `required help:"The bagel log to tally feedback for"`
	Channel string `required help:"The name/slackid of the slack channel to send the status update"`
}

func (cmd *CmdMsgStats) Run(ctx *kong.Context, db *gorm.DB, s *Slack) (err error) {
	log, err := BagelLog_Fetch(db, cmd.Log)
	if err != nil {
		return err
	}

	var bagels []Bagel
	db.Model(&log).Association("Bagels").Find(&bagels)
	var firstBagelCompleted Bagel
	var incomplete, planned, completed int
	for _, bagel := range bagels {
		var users []User
		db.Model(&bagel).Association("Users").Find(&users)

		userCount := 0
		for _, user := range users {
			if !user.IsBot {
				userCount++
			}
		}

		if bagel.IsCompleted {
			completed += userCount

			if firstBagelCompleted.ID == 0 {
				firstBagelCompleted = bagel
			} else if bagel.FeedbackDate < firstBagelCompleted.FeedbackDate {
				firstBagelCompleted = bagel
			}
		} else if bagel.IsPlanned {
			planned += userCount
		} else {
			incomplete += userCount
		}
	}

	var firstGroupCompleted []string
	if firstBagelCompleted.ID != 0 {
		var firstGroupUsers []User
		db.Model(&firstBagelCompleted).Association("Users").Find(&firstGroupUsers)

		for _, user := range firstGroupUsers {
			if !user.IsBot {
				firstGroupCompleted = append(firstGroupCompleted, user.Name)
			}
		}
	}

	channel, err := s.FindChannel(cmd.Channel, cmd.Channel)
	if err != nil {
		return err
	}

	if channel == nil {
		_, err = io.WriteString(ctx.Stderr, "no such channel "+cmd.Channel)
		return err
	}

	blocks := SlackBlocks_FeedbackStatistics(completed, firstGroupCompleted, planned)

	if err = s.ChatPostMessage(channel.ID, "Updated stats about bagel chats", blocks); err != nil {
		_, err = io.WriteString(ctx.Stderr, err.Error())
		return err
	}

	return nil
}

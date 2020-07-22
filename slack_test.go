package main

import (
	"os"
	"testing"
)

func NewTestingSlack(t *testing.T) Slack {
	apiKey := os.Getenv("BAGEL_SLACK_API_KEY")
	if len(apiKey) == 0 {
		t.Skip("Slack API key missing")
	}

	slack := Slack{Token: apiKey}

	err := slack.ApiTest()
	if err != nil {
		t.Skip("Unable to make trivial call to Slack API: " + err.Error())
	}

	return slack
}

func TestSlack_ConversationsList(t *testing.T) {
	slack := NewTestingSlack(t)
	_, err := slack.ConversationsList(true, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestSlack_ConversationsMembers(t *testing.T) {
	slack := NewTestingSlack(t)
	channels, err := slack.ConversationsList(true, nil)
	if err != nil {
		t.Error(err)
	}

	if len(channels) == 0 {
		t.Skip("No channel found to retrieve members for")
	}

	channel := channels[0]

	_, err = slack.ConversationsMembers(channel.ID)
	if err != nil {
		t.Error(err)
	}
}

func TestSlack_UsersConversations(t *testing.T) {
	slack := NewTestingSlack(t)
	_, err := slack.UsersConversations(true, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestSlack_UsersList(t *testing.T) {
	slack := NewTestingSlack(t)
	_, err := slack.UsersList()
	if err != nil {
		t.Error(err)
	}
}

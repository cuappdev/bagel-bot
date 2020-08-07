package main

import (
	"fmt"
	"strings"
)

func SlackBlocks_FeedbackMsg(msg FeedbackMsg, text string) []interface{} {
	if text == "" {
		text = "How's it going? I'm here to get an update on the status of your bagel chat 😁"
	}

	return []interface{}{
		map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": text,
			},
		},
		map[string]interface{}{
			"type": "actions",
			"elements": []map[string]interface{}{
				{
					"type":      "button",
					"action_id": msg.IncompleteActionID,
					"text": map[string]string{
						"type": "plain_text",
						"text": "We're not planning to meet",
					},
				},
				{
					"type":      "button",
					"action_id": msg.PlannedActionID,
					"text": map[string]string{
						"type": "plain_text",
						"text": "We've planned our bagel chat",
					},
				},
				{
					"type":      "button",
					"action_id": msg.CompletedActionID,
					"text": map[string]string{
						"type": "plain_text",
						"text": "We've had our bagel chat",
					},
				},
			},
		},
	}
}

func ToEnglish_JoinAnd(elements []string) string {
	if len(elements) == 0 {
		return ""
	} else if len(elements) == 1 {
		return elements[0]
	} else if len(elements) == 2 {
		return fmt.Sprintf("%s and %s", elements[0], elements[1])
	} else {
		return strings.Join(elements[:len(elements)-1], ", ") + ", and " + elements[len(elements)-1]
	}
}

func SlackBlocks_FeedbackStatistics(completed int, firstGroupCompleted []string, planned int) []interface{} {
	text := "📊 I've got some statistics about bagel chats. Thank you to all who have participated."

	var firstGroupShoutOut string
	if firstGroupCompleted == nil {
		firstGroupShoutOut = ""
	} else {
		firstGroupShoutOut = fmt.Sprintf("Shout out to %s for being the first group to chat it up.", ToEnglish_JoinAnd(firstGroupCompleted))
	}

	return []interface{}{
		map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": text,
			},
		},
		map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": fmt.Sprintf(" • *%d* members have *completed* their bagel chats. %s", completed, firstGroupShoutOut),
			},
		},
		map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": fmt.Sprintf(" • *%d* members have *planned* their bagel chats. Make sure to post a selfie (or a screenshot) to the #bagel-chats channel, and don't forget to mark that you've completed your bagel chat once you've met up.", planned),
			},
		},
	}
}

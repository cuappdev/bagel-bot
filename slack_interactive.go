package main

import (
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"strings"
)

type SlackInteractiveHandler struct {
	DB *gorm.DB
	Slack *Slack
}

func (params SlackInteractiveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestBody, err := verifyRequest(r)
	if err != nil {
		panic(err)
	}

	payload, err := url.QueryUnescape(string(requestBody))
	if err != nil {
		panic(err)
	}
	payload = strings.TrimPrefix(payload, "payload=")

	var jsonBody map[string]interface{}
	err = json.Unmarshal([]byte(payload), &jsonBody)
	if err != nil {
		panic(err)
	}

	var ok bool

	var opaqueActions []interface{}
	if opaqueActions, ok = jsonBody["actions"].([]interface{}); !ok {
		log.Warning("unable to extract actions", jsonBody)
		return
	}

	for _, opaqueAction := range opaqueActions {
		var action map[string]interface{}
		if action, ok = opaqueAction.(map[string]interface{}); !ok {
			log.Warning("unable to convert action to json dict: ", jsonBody)
			continue
		}

		var actionId string
		if actionId, ok = action["action_id"].(string); !ok {
			log.Warning("unable to convert actionId to string: ", jsonBody)
			continue
		}

		if strings.HasPrefix(actionId, "feedback_msg:") {
			if err = params.handleFeedbackMsg(jsonBody, actionId); err != nil {
				panic(err)
			}
		} else {
			log.Warning("unknown actionId: " + actionId)
			continue
		}
	}
}

func (params SlackInteractiveHandler) handleFeedbackMsg(reqBody map[string]interface{}, actionId string) error {
	var bagel Bagel
	var feedbackMsg FeedbackMsg
	var incomplete, planned, completed bool

	if params.DB.Where("incomplete_action_id = ?", actionId).Find(&feedbackMsg); feedbackMsg.ID != 0 {
		if err := params.DB.Model(&feedbackMsg).Association("bagels").Find(&bagel); err != nil {
			return err
		}
		incomplete = true
	} else if params.DB.Where("planned_action_id = ?", actionId).Find(&feedbackMsg); feedbackMsg.ID != 0 {
		if err := params.DB.Model(&feedbackMsg).Association("bagels").Find(&bagel); err != nil {
			return err
		}
		planned = true
	} else if params.DB.Where("completed_action_id = ?", actionId).Find(&feedbackMsg); feedbackMsg.ID != 0 {
		if err := params.DB.Model(&feedbackMsg).Association("bagels").Find(&bagel); err != nil {
			return err
		}
		completed = true
	} else {
		log.Warning("unable to find feedback msg with actionId ", actionId)
	}

	if bagel.ID == 0 {
		log.Warning("unable to find bagel from feedback msg")
		return nil
	}

	log.Infof("updating status for bagel %d: IsPlanned(%t) IsCompleted(%t)", bagel.ID, planned, completed)

	bagel.IsPlanned = planned
	bagel.IsCompleted = completed
	params.DB.Save(&bagel)

	var blocks []interface{}
	if incomplete {
		blocks = SlackBlocks_FeedbackMsg(feedbackMsg, "😥 I'm sad that you've decided not to meet up. If you change your mind, you can always push the buttons below to let me know.")
	} else if planned {
		blocks = SlackBlocks_FeedbackMsg(feedbackMsg, "😀 Great! I'm exited that you've planned to chat. Once you do meet up, go ahead and push the button below to let me know you've 🥯'd. And don't forget about posting a selfie to the #bagel-chats channel.")
	} else if completed {
		blocks = SlackBlocks_FeedbackMsg(feedbackMsg, "😁 Perfect! Thank you for contributing to AppDev's social culture.")
	}

	var ok bool

	var container map[string]interface{}
	if container, ok = reqBody["container"].(map[string]interface{}); !ok {
		log.Warning("unable to extract container", reqBody)
		return nil
	}

	var ts string
	if ts, ok = container["message_ts"].(string); !ok {
		log.Warning("unable to extract timestamp", reqBody)
		return nil
	}

	var channelID string
	if channelID, ok = container["channel_id"].(string); !ok {
		log.Warning("unable to extract channel id", reqBody)
		return nil
	}

	return params.Slack.ChatUpdate(channelID, ts, "Someone responded to bagel feedback", blocks)
}


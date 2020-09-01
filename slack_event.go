package main

import (
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

type SlackEventHandler struct {
	DB *gorm.DB
	Slack *Slack
}

func (params SlackEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestBody, err := verifyRequest(r)
	if err != nil {
		panic(err)
	}

	var jsonBody map[string]interface{}
	err = json.Unmarshal(requestBody, &jsonBody)
	if err != nil {
		panic(err)
	}

	requestType := jsonBody["type"].(string)
	log.Info("Handling event of type", requestType)
	switch requestType {
	case "url_verification":
		err = handleUrlVerification(w, r, jsonBody)
	default:
		log.Notice("Unhandled request type: " + requestType)
	}
	if err != nil {
		panic(err)
	}
}

func handleUrlVerification(w http.ResponseWriter, r *http.Request, requestBody map[string]interface{}) error {
	challenge := requestBody["challenge"].(string)
	_, err := w.Write([]byte(challenge))
	return err
}


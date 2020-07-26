package main

import (
	"encoding/json"
	"net/http"
)

func SlackInteractive_Handler(w http.ResponseWriter, r *http.Request) {
	requestBody, err := verifyRequest(r)
	if err != nil {
		panic(err)
	}

	var jsonBody map[string]interface{}
	err = json.Unmarshal(requestBody, &jsonBody)
	if err != nil {
		panic(err)
	}
}


package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = r.Body.Close()
	if err != nil {
		panic(err)
	}

	var jsonBody map[string]interface{}
	err = json.Unmarshal(bytes, &jsonBody)
	if err != nil {
		panic(err)
	}

	err = verifyRequest(r)
	if err != nil {
		panic(err)
	}

	requestType := jsonBody["type"].(string)
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

var didWarnUnverifiedRequests = false

func verifyRequest(r *http.Request) error {
	if SlackSigningSecret == "" {
		if didWarnUnverifiedRequests {
			return nil
		}

		log.Warning("The BAGEL_SLACK_SIGNING_SECRET environment variable is empty; requests will not be verified")
		didWarnUnverifiedRequests = true
	}

	timestamp := r.Header.Get("X-Slack-Request-Timestamp")
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	if timestampInt > now {
		return errors.New("request has a timestamp in the future")
	}

	var maxTimestampDifference int64 = 60
	if now - timestampInt > maxTimestampDifference {
		return errors.New("request is more than " + strconv.FormatInt(maxTimestampDifference, 10) + " seconds in the past")
	}

	signatureBasestring := "v0:" + timestamp + ":" + string(requestBody)

	h := hmac.New(sha256.New, []byte(SlackSigningSecret))
	h.Write([]byte(signatureBasestring))
	hexDigest := "v0" + hex.EncodeToString(h.Sum(nil))

	slackSignature := r.Header.Get("X-Slack-Signature")
	if !hmac.Equal([]byte(hexDigest), []byte(slackSignature)) {
		return errors.New("computed signature and slack signature are not the same; possible malicious slack event")
	}

	return nil
}

func handleUrlVerification(w http.ResponseWriter, r *http.Request, requestBody map[string]interface{}) error {
	err := verifyRequest(r)
	if err != nil {
		return err
	}

	challenge := requestBody["challenge"].(string)
	_, err = w.Write([]byte(challenge))
	return err
}

func SlackEventsListenAndServe() {
	http.HandleFunc("/slack/action-event", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

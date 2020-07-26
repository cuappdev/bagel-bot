package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var didWarnUnverifiedRequests = false

func verifyRequest(r *http.Request) ([]byte, error) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if SlackSigningSecret == "" {
		if !didWarnUnverifiedRequests {
			log.Warning("The BAGEL_SLACK_SIGNING_SECRET environment variable is empty; requests will not be verified")
			didWarnUnverifiedRequests = true
		}

		return requestBody, nil
	}

	timestamp := r.Header.Get("X-Slack-Request-Timestamp")

	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	if timestampInt > now {
		return nil, errors.New("request has a timestamp in the future")
	}

	var maxTimestampDifference int64 = 60
	if now-timestampInt > maxTimestampDifference {
		return nil, errors.New("request is more than " + strconv.FormatInt(maxTimestampDifference, 10) + " seconds in the past")
	}

	signatureBasestring := "v0:" + timestamp + ":" + string(requestBody)

	h := hmac.New(sha256.New, []byte(SlackSigningSecret))
	h.Write([]byte(signatureBasestring))
	hexDigest := "v0=" + hex.EncodeToString(h.Sum(nil))

	slackSignature := r.Header.Get("X-Slack-Signature")
	if !hmac.Equal([]byte(hexDigest), []byte(slackSignature)) {
		return nil, errors.New("computed signature and slack signature are not the same; possible malicious slack event")
	}

	return requestBody, nil
}

func SlackListenAndServe() {
	http.HandleFunc("/slack/action-event", SlackEvent_Handler)
	http.HandleFunc("/slack/interactive-endpoint", SlackInteractive_Handler)
	log.Fatal(http.ListenAndServe(":29138", nil))
}


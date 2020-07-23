package main

import "os"

var SlackSigningSecret = os.Getenv("BAGEL_SLACK_SIGNING_SECRET")
var SlackApiKey = os.Getenv("BAGEL_SLACK_API_KEY")

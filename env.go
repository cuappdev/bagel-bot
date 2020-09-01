package main

import "os"

var SlackSigningSecret = os.Getenv("BAGEL_SLACK_SIGNING_SECRET")
var SlackApiKey = os.Getenv("BAGEL_SLACK_API_KEY")
var PostgresHost = os.Getenv("BAGEL_POSTGRES_HOST")
var PostgresPort = os.Getenv("BAGEL_POSTGRES_PORT")
var PostgresUser = os.Getenv("BAGEL_POSTGRES_USER")
var PostgresPassword = os.Getenv("BAGEL_POSTGRES_PASSWORD")
var PostgresDbName = os.Getenv("BAGEL_POSTGRES_DB_NAME")

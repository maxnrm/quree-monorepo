package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load(".env")

// COMMON
var NATS_URL = getenvStr("NATS_URL", "nats://localhost:4222")
var TG_API_BASE_URL = getenvStr("TG_API_BASE_URL", "https://api.telegram.org")
var POSTGRES_CONN_STRING = getenvStr("POSTGRES_CONN_STRING", "")

var USER_PLACEHOLDER_BOT_TOKEN = getenvStr("USER_PLACEHOLDER_BOT_TOKEN", "")

var USER_BOT_TOKEN = getenvStr("USER_BOT_TOKEN", "")
var USER_WEBAPP_URL = getenvStr("USER_WEBAPP_URL", "")

var USER_WEBSERVER_PORT = getenvStr("USER_WEBSERVER_PORT", "3000")

var ADMIN_BOT_TOKEN = getenvStr("ADMIN_BOT_TOKEN", "")
var ADMIN_WEBAPP_URL = getenvStr("ADMIN_WEBAPP_URL", "")
var ADMIN_AUTH_CODE = getenvStr("ADMIN_AUTH_CODE", "")

// MESSAGES
var NATS_MESSAGES_STREAM = getenvStr("NATS_MESSAGES_STREAM", "tg-messages")
var NATS_RECEIVER_MESSAGES_SUBJECT = getenvStr("NATS_MESSAGES_SUBJECT", "tg.messages.*")

var NATS_USER_MESSAGES_CONSUMER = getenvStr("NATS_MESSAGES_CONSUMER", "tg-messages-user-consumer")
var NATS_ADMIN_MESSAGES_CONSUMER = getenvStr("NATS_MESSAGES_CONSUMER", "tg-messages-admin-consumer")

var NATS_USER_MESSAGES_SUBJECT = getenvStr("NATS_MESSAGES_SUBJECT", "tg.messages.user")
var NATS_ADMIN_MESSAGES_SUBJECT = getenvStr("NATS_MESSAGES_SUBJECT", "tg.messages.admin")

var RATE_LIMIT_GLOBAL = getenvInt("RATE_LIMIT_GLOBAL", 30)
var RATE_LIMIT_BURST_GLOBAL = getenvInt("RATE_LIMIT_BURST_GLOBAL", 30)

var RATE_LIMIT_USER = getenvInt("RATE_LIMIT_USER", 1)
var RATE_LIMIT_BURST_USER = getenvInt("RATE_LIMIT_BURST_USER", 1)

var EVENT_VISIT_DELAY_MINUTES = getenvInt("EVENT_VISIT_DELAY_MINUTES", 5)
var FINISH_PASS_DATE = getenvStr("FINISH_PASS_DATE", "2024-03-15")

// IMGPROXY
var IMGPROXY_INTERNAL_URL = getenvStr("IMGPROXY_INTERNAL_URL", "http://localhost:8080")
var IMGPROXY_PUBLIC_URL = getenvStr("IMGPROXY_PUBLIC_URL", "http://localhost:8080")

// S3
var S3_ENDPOINT = getenvStr("S3_ENDPOINT", "http://localhost:9000")
var S3_BUCKET = getenvStr("S3_BUCKET", "")
var S3_ACCESS_KEY = getenvStr("S3_ACCESS_KEY", "")
var S3_SECRET_KEY = getenvStr("S3_SECRET_KEY", "")

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

func getenvStr(key string, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func getenvInt(key string, defaultValue int) int {
	s := getenvStr(key, "")
	if s == "" {
		return defaultValue
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return v
}

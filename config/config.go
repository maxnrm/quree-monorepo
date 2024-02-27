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
var USER_BOT_TOKEN = getenvStr("USER_BOT_TOKEN", "")
var ADMIN_BOT_TOKEN = getenvStr("ADMIN_BOT_TOKEN", "")
var USER_WEBAPP_URL = getenvStr("USER_WEBAPP_URL", "")
var ADMIN_WEBAPP_URL = getenvStr("ADMIN_WEBAPP_URL", "")
var POSTGRES_CONN_STRING = getenvStr("POSTGRES_CONN_STRING", "")

// MESSAGES
var NATS_MESSAGES_STREAM = getenvStr("NATS_MESSAGES_STREAM", "tg-messages")
var NATS_MESSAGES_CONSUMER = getenvStr("NATS_MESSAGES_CONSUMER", "tg-messages-consumer")
var NATS_MESSAGES_SUBJECT = getenvStr("NATS_MESSAGES_SUBJECT", "messages.*")

var RATE_LIMIT_GLOBAL = getenvInt("RATE_LIMIT_GLOBAL", 30)
var RATE_LIMIT_PER_USER = getenvInt("RATE_LIMIT_PER_USER", 1)
var RATE_LIMIT_BURST = getenvInt("RATE_LIMIT_BURST", 30)

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

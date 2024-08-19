package utils

import (
	"log/slog"
	"os"
)

var USER_ID_KEY = "user_id"
var ACCESS_TOKEN_KEY = "access_token"

var TextLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

// or
var JsonLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

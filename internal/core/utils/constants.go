package utils

import (
	"log/slog"
	"os"
)

type ContextKey string

var USER_ID_KEY string = "user_id"
var ACCESS_TOKEN_KEY string = "access_token"

var TextLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))

// or
var JsonLogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

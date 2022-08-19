package main

import (
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

var logger logr.Logger

func init() {
	zl := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).With().Timestamp().Caller().Logger()
	logger = zerologr.New(&zl)
}

package main

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

var logger logr.Logger

func init() {
	zl := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	logger = zerologr.New(&zl)
}

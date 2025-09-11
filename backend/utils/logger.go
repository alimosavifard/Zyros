package utils

import (
	"github.com/rs/zerolog"
	"os"
)

// Logger از پکیج اصلی برای لاگینگ
var logger *zerolog.Logger

// InitLogger برای راه‌اندازی logger
func InitLogger() *zerolog.Logger {
	if logger == nil {
		l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
		logger = &l
	}
	return logger
}
package errorlog

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ErrorLog() *zap.Logger {
	cfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths: []string{"error.log"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			MessageKey:     "message",
			CallerKey:      "caller",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatal("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	file, err := os.Open("error.log")
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	return logger
}
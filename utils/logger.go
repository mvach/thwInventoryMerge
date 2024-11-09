package utils

import (
	"log"
	"os"
)

//counterfeiter:generate . Logger
type Logger interface {
	Info(message string)

	InfoIndented(message string)

	Warn(message string)

	WarnIndented(message string)

	Error(message string)

	Fatal(message string)
}

type logger struct {}

func NewLogger() Logger {
	return logger{}
}

func (l logger) Info(message string) {
	log.Println("[INFO] "+ message)
}

func (l logger) InfoIndented(message string) {
	log.Println("[INFO]    "+ message)
}

func (l logger) Warn(message string) {
	log.Println("[WARN] "+ message)
}

func (l logger) WarnIndented(message string) {
	log.Println("[WARN]    "+ message)
}

func (l logger) Error(message string) {
	log.Println("[ERROR] "+ message)
}


func (l logger) Fatal(message string) {
	log.Println("[FATAL] "+ message)
	os.Exit(1)
}


package logger

import (
	"log"
	"os"
)

var infoLogger *log.Logger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var warningLogger *log.Logger = log.New(os.Stdout, "WARNING\t", log.Ldate|log.Ltime)
var errorLogger *log.Logger = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

func Info(msg string) {
	infoLogger.Println(msg)
}

func Warning(msg string) {
	warningLogger.Println(msg)
}

func Error(msg string) {
	errorLogger.Println(msg)
}

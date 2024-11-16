package main

import (
	"log"
	"os"
)

type someLogger struct{}

var infoLogger, warningLogger, errorLogger *log.Logger

var logger someLogger

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
}

func setLogger() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l someLogger) Info(v ...any) {
	infoLogger.Println(v...)
}
func (l someLogger) Warn(v ...any) {
	warningLogger.Println(v...)
}
func (l someLogger) Error(v ...any) {
	errorLogger.Println(v...)
}

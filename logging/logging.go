package logging

import (
	"fmt"
	"log"
	"time"

	"gitlab.danawa.com/devops/file-searcher/model"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logging model.Logging
	Logger  = &lumberjack.Logger{}
	level   = 2
)

func LoadLogging(logging model.Logging) {
	Logging = logging
	if logging.Level == "trace" {
		level = 0
	} else if logging.Level == "debug" {
		level = 1
	} else if logging.Level == "info" {
		level = 2
	} else if logging.Level == "warn" {
		level = 3
	} else if logging.Level == "error" {
		level = 4
	} else if logging.Level == "fatal" {
		level = 5
	} else if logging.Level == "off" {
		level = 6
	}
	Logger = &lumberjack.Logger{
		Filename:   Logging.Filename,
		MaxSize:    Logging.MaxSize,
		MaxBackups: Logging.MaxBackups,
		MaxAge:     Logging.MaxAge,
		Compress:   Logging.Compress,
	}
	log.SetOutput(Logger)
	Info("Loaded Logging")
}

func Rotate() {
	Logger.Rotate()
}

func write(text string) {
	text = fmt.Sprintf("%v", time.Now()) + " " + text
	fmt.Println(text)
	_, _ = Logger.Write([]byte(text + "\n"))
}

func Trace(message string) {
	if level <= 0 {
		write("[TRACE] " + message)
	}
}
func Debug(message string) {
	if level <= 1 {
		write("[DEBUG] " + message)
	}
}
func Info(message string) {
	if level <= 2 {
		write("[INFO] " + message)
	}
}
func Warn(message string) {
	if level <= 3 {
		write("[WARN] " + message)
	}
}
func Error(message string) {
	if level <= 4 {
		write("[ERROR] " + message)
	}
}
func Fatal(message string) {
	if level <= 5 {
		write("[FATAL] " + message)
	}
}

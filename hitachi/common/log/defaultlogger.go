package common

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// LogLevel defines the log level type.
type LogLevel int

// Log level constants.
const (
	DEBUGT LogLevel = iota
	INFOT
	WARNT
	ERRORT
)

// currentLogLevel holds the global log level for the application.
var currentLogLevel = INFOT

// DefaultLogger is the default implementation of logger.
type DefaultLogger struct{}

// logWriterFile is used to write to our log file.
var logWriterFile *log.Logger

// init initializes the logger with log file and log level from environment.
func init() {
	path := "/var/log/hitachi/terraform/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating log directory:", err)
			return
		}
	}

	finalPath := fmt.Sprintf("%s%s", path, "hitachi-terraform.log")
	setNewLogFile(finalPath, 1, 1)

	// Set log level from environment variable.
	setLogLevelFromEnv()
}

// NewDefaultLogger creates a new default logger.
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}

// SetLogLevel sets the current global log level.
func SetLogLevel(level LogLevel) {
	currentLogLevel = level
}

// setLogLevelFromEnv sets the log level based on an environment variable.
func setLogLevelFromEnv() {
	envLogLevel := os.Getenv("TF_LOG_LEVEL")
	switch envLogLevel {
	case "DEBUG":
		SetLogLevel(DEBUGT)
	case "INFO":
		SetLogLevel(INFOT)
	case "WARN":
		SetLogLevel(WARNT)
	case "ERROR":
		SetLogLevel(ERRORT)
	default:
		// fmt.Println("Invalid or no LOG_LEVEL environment variable set. Defaulting to INFO.")
		SetLogLevel(INFOT)
	}
}

// shouldLog determines if a log should be written based on the current log level.
func shouldLog(level LogLevel) bool {
	return level >= currentLogLevel
}

// setNewLogFile sets a new log file on initialization.
func setNewLogFile(fname string, maxSize int, maxBackups int) error {
	newLogFile, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Unable to open log file: %+v\n", err)
		return err
	}

	logWriterFile = log.New(newLogFile, "", 0)
	logWriterFile.SetOutput(&lumberjack.Logger{
		Filename:   fname,
		MaxSize:    maxSize,    // megabytes
		MaxBackups: maxBackups, // number of backups
	})
	logWriterFile.Print("Starting new log file")
	return nil
}

// formatLog formats the log statement with severity, time, and message.
func formatLog(severity string, message string) string {
	_, funcname, filename, lineno := getSourceFileInfo(3)
	filesource := fmt.Sprintf("%s:%s:%d", funcname, filename, lineno)
	logStatement := fmt.Sprintf("%s\t%s\t%s\t%s", time.Now().Format("15:04:05"), severity, filesource, message)
	return logStatement
}

// Logging functions:

func (l *DefaultLogger) WriteDebug(message string, a ...interface{}) {
	if shouldLog(DEBUGT) {
		log := formatLog("DEBUG", fmt.Sprintf(message, a...))
		logWriterFile.Println(log)
	}
}

func (l *DefaultLogger) WriteInfo(message interface{}, a ...interface{}) {
	if shouldLog(INFOT) {
		msg := fmt.Sprintf("%v", message)
		log := formatLog("INFO", fmt.Sprintf(msg, a...))
		logWriterFile.Println(log)
	}
}

func (l *DefaultLogger) WriteWarn(message interface{}, a ...interface{}) {
	if shouldLog(WARNT) {
		msg := fmt.Sprintf("%v", message)
		log := formatLog("WARN", fmt.Sprintf(msg, a...))
		logWriterFile.Println(log)
	}
}

func (l *DefaultLogger) WriteError(message interface{}, a ...interface{}) {
	if shouldLog(ERRORT) {
		msg := fmt.Sprintf("%v", message)
		log := formatLog("ERROR", fmt.Sprintf(msg, a...))
		logWriterFile.Println(log)
	}
}

func (l *DefaultLogger) WriteEnter(a ...interface{}) {
	log := formatLog(string(ENTER), "")
	logWriterFile.Println(log)
}

func (l *DefaultLogger) WriteParam(format string, value interface{}) {
	if shouldLog(INFOT) {

		log := formatLog(string(PARAM), fmt.Sprintf("%v", value))
		logWriterFile.Println(log)
	}
}

func (l *DefaultLogger) WriteExit() {

	if shouldLog(INFOT) {

		log := formatLog(string(EXIT), "")
		logWriterFile.Println(log)
	}
}

package common

import (
	"fmt"
	"log"
	"os"
	"reflect"
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

const (
	LOGDIR        = "/var/log/hitachi/terraform/"
	LOGFILENAME   = "hitachi-terraform.log"
	LOGMAXSIZE    = 10 // MB
	LOGMAXBACKUPS = 10
)

// currentLogLevel holds the global log level for the application.
var currentLogLevel = INFOT

// DefaultLogger is the default implementation of logger.
type DefaultLogger struct{}

// logWriterFile is used to write to our log file.
var logWriterFile *log.Logger

// init initializes the logger with log file and log level from environment.
func init() {
	// Ensure log directory exists
	if _, err := os.Stat(LOGDIR); os.IsNotExist(err) {
		if err := os.MkdirAll(LOGDIR, os.ModePerm); err != nil {
			fmt.Println("Failed to create log directory:", err)
			return
		}
	}

	finalPath := LOGDIR + LOGFILENAME
	if err := setNewLogFile(finalPath, LOGMAXSIZE, LOGMAXBACKUPS); err != nil {
		fmt.Println("Log setup failed:", err)
	}

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
	envLogLevel := os.Getenv("TF_LOG")
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
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		// MaxAge:     28, // Optionally uncomment to limit log age
	})
	logWriterFile.Printf("Starting new log file. MaxSize:%dMB MaxBackups:%d", maxSize, maxBackups)
	return nil
}

// formatLog formats the log statement with severity, time, and message.
func formatLog(severity string, message string) string {
	_, funcname, filename, lineno := getSourceFileInfo(3)
	filesource := fmt.Sprintf("%s:%s:%d", funcname, filename, lineno)
	logStatement := fmt.Sprintf("%s\t%s\t%s\t%s", time.Now().Format("2006-01-02 15:04:05 MST"), severity, filesource, message)
	return logStatement
}

// Logging functions:

func (l *DefaultLogger) WriteDebug(message string, a ...interface{}) {
	if shouldLog(DEBUGT) {
		safe := safeArgs(a)
		log := formatLog("DEBUG", fmt.Sprintf(message, safe...))
		logWriterFile.Println(log)
	}
}

func (l *DefaultLogger) WriteInfo(message interface{}, a ...interface{}) {
	if shouldLog(INFOT) {
		msg := fmt.Sprintf("%v", message)
		safe := safeArgs(a)
		log := formatLog("INFO", fmt.Sprintf(msg, safe...))
		logWriterFile.Println(log)
	}
}

func (l *DefaultLogger) WriteWarn(message interface{}, a ...interface{}) {
	if shouldLog(WARNT) {
		msg := fmt.Sprintf("%v", message)
		safe := safeArgs(a)
		log := formatLog("WARN", fmt.Sprintf(msg, safe...))
		logWriterFile.Println(log)
	}
}

func (l *DefaultLogger) WriteError(message interface{}, a ...interface{}) {
	if shouldLog(ERRORT) {
		msg := fmt.Sprintf("%v", message)
		safe := safeArgs(a)
		log := formatLog("ERROR", fmt.Sprintf(msg, safe...))
		logWriterFile.Println(log)
	}
}

func (l *DefaultLogger) WriteEnter(a ...interface{}) {
	log := formatLog(string(ENTER), "")
	logWriterFile.Println(log)
}

func (l *DefaultLogger) WriteParam(format string, value interface{}) {
	if shouldLog(INFOT) {
		safe := safeValue(value)
		log := formatLog(string(PARAM), fmt.Sprintf("%v", safe))
		logWriterFile.Println(log)
	}
}

func (l *DefaultLogger) WriteExit() {

	if shouldLog(INFOT) {

		log := formatLog(string(EXIT), "")
		logWriterFile.Println(log)
	}
}

// safeValue converts a possibly-nil pointer into its zero value.
// Non-pointer values are returned untouched.
func safeValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	rv := reflect.ValueOf(v)

	// If not a pointer, return as is
	if rv.Kind() != reflect.Ptr {
		return v
	}

	// Pointer but nil → return zero value of the element type
	if rv.IsNil() {
		elemType := rv.Type().Elem()
		zero := reflect.Zero(elemType)
		return zero.Interface()
	}

	// Pointer with value → return the dereferenced value
	return rv.Elem().Interface()
}

// Apply safeValue to all variadic args
func safeArgs(args []interface{}) []interface{} {
	out := make([]interface{}, len(args))
	for i, a := range args {
		out[i] = safeValue(a)
	}
	return out
}

package common

import (
	"fmt"
	logger "log"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	LOGDIR        = "/var/log/hitachi/terraform/"
	LOGFILENAME   = "hitachi-terraform.log"
	LOGMAXSIZE    = 10 // MB
	LOGMAXBACKUPS = 10
)

// Default Log Provider - the default implementation of hitachi log interface ILogger
type DefaultLogger struct {
}

// use this to write to our log file
var logWriterFile *logger.Logger

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
}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}

// ─────────────────────────────────────────────────────────────
// Set new log file with rotation support
// ─────────────────────────────────────────────────────────────
func setNewLogFile(fname string, maxsize, maxbackups int) error {
	newLogFile, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("OpenLogFile: unable to open log file: %+v\n", err)
		return err
	}

	logWriterFile = logger.New(newLogFile, "", 0)
	logWriterFile.SetOutput(&lumberjack.Logger{
		Filename:   fname,
		MaxSize:    maxsize,
		MaxBackups: maxbackups,
		// MaxAge:     28, // Optionally uncomment to limit log age
	})

	logWriterFile.Print("Starting new log file")
	return nil
}

func formatLog(severity string, message string) string {
	_, funcname, filename, lineno := getSourceFileInfo(3)
	filesource := fmt.Sprintf("%s:%s:%d", funcname, filename, lineno)
	logStatement := fmt.Sprintf("%s\t%s\t%s\t%s", time.Now().Format("15:04:05"), severity, filesource, message)
	return logStatement
}

func (l *DefaultLogger) WriteEnter(a ...interface{}) {
	log := formatLog(string(ENTER), "")
	logWriterFile.Println(log)
}
func (l *DefaultLogger) WriteParam(format string, value interface{}) {
	log := formatLog(string(PARAM), fmt.Sprintf("%v", value))
	logWriterFile.Println(log)
}

func (l *DefaultLogger) WriteInfo(message interface{}, a ...interface{}) {
	msg := fmt.Sprintf("%v", message)
	log := formatLog(string(INFO), fmt.Sprintf(msg, a...))
	logWriterFile.Println(log)

}

func (l *DefaultLogger) WriteWarn(message interface{}, a ...interface{}) {
	msg := fmt.Sprintf("%v", message)
	log := formatLog(string(WARN), fmt.Sprintf(msg, a...))
	logWriterFile.Println(log)
}

func (l *DefaultLogger) WriteError(message interface{}, a ...interface{}) {
	msg := fmt.Sprintf("%v", message)
	log := formatLog(string(ERROR), fmt.Sprintf(msg, a...))
	logWriterFile.Println(log)
}

func (l *DefaultLogger) WriteDebug(message string, a ...interface{}) {
	log := formatLog(string(DEBUG), fmt.Sprintf(message, a...))
	logWriterFile.Println(log)
}

func (l *DefaultLogger) WriteExit() {
	log := formatLog(string(EXIT), "")
	logWriterFile.Println(log)
}

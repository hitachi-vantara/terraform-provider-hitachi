package common

import (
	"fmt"
	logger "log"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Default Log Provider - the default implementation of hitachi log interface ILogger
type DefaultLogger struct {
}

// use this to write to our log file
var logWriterFile *logger.Logger

// init for logger
func init() {

	path := "/var/log/hitachi/terraform/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if os.IsNotExist(err) {
            err = os.MkdirAll(path, os.ModePerm)
        }
		if err != nil {
			fmt.Println(err)
		}
	}

	finalpath := fmt.Sprintf("%s%s", path, "hitachi-terraform.log")
	setNewLogFile(finalpath, 1, 1)

}

func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}

func fmtPrintf(a string, b interface{}) {

}

// setNewLogFile set new log file on init
func setNewLogFile(fname string,
	maxsize int,
	maxbackups int) error {

	fmtPrintf("OpenLogFile:180  Log Provider setNewLogFile, fname: %s\n", fname)
	fmtPrintf("OpenLogFile: maxsize: %d\n", maxsize)
	fmtPrintf("OpenLogFile: maxbackups: %d\n", maxbackups)
	newLogFile, err := os.OpenFile(fname,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		fmtPrintf("OpenLogFile newLogFile: %v\n", newLogFile)
		logWriterFile = logger.New(newLogFile, "", 0)
		logWriterFile.SetOutput(&lumberjack.Logger{
			Filename:   fname,
			MaxSize:    maxsize,
			MaxBackups: maxbackups,
			// MaxAge:     28, //days
		})

		if logWriterFile != nil {
			// write to rotated log files
			fmtPrintf("OpenLogFile: file opened: %v\n", fname)
			logWriterFile.Print("starting new log file")
		}

	} else {
		// rlogIssue("Unable to open log file: %s", err)
		fmtPrintf("OpenLogFile  Log Provider unable to open log file: %+v\n", err)
		return err
	}

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

package common



type LogSeverity string

const (
	INFO  LogSeverity = "INFO"
	WARN  LogSeverity = "WARN"
	ERROR LogSeverity = "ERROR"
	ENTER LogSeverity = "ENTER"
	EXIT  LogSeverity = "EXIT"
	PARAM LogSeverity = "PARAM"
	DEBUG LogSeverity = "DEBUG"
)

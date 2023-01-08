package logger

// Level - Severity levels.
type Level int

const (
	Debug Level = 0
	Info  Level = 1
	Warn  Level = 2
	Error Level = 3
	Fatal Level = 4
)

func (d Level) String() string {
	switch d {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARNING"
	case Error:
		return "ERROR"
	case Fatal:
		return "CRITICAL"
	}
	return "DEBUG"
}

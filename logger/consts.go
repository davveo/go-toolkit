package logger

type (
	LogLevel int8
	Format   int8
)

// log level
const (
	LevelDebug LogLevel = iota - 1
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

const (
	callerSkip = 2
	// perm: rwx
	permRWX = 0777
	// perm: rw-
	permRW = 0666
)

const (
	FormatText Format = iota + 1
	FormatJSON
)

const (
	KeyTime       = "timestamp"
	KeyLevel      = "level"
	KeyCaller     = "location"
	KeyMessage    = "msg"
	KeyStacktrace = "estack"
	KeyBizName    = "bo"
	KeyName       = "logger"
)

const (
	DefaultColor      = false
	DefaultStack      = true
	DefaultCallerSkip = 0
	DefaultLevel      = LevelDebug
	DefaultFormat     = FormatText
	DefaultTimeLayout = "2006/01/02 15:04:05"
)

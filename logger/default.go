package logger

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"runtime"
	"time"
)

type logger struct {
	color      bool
	stack      bool
	ts         string
	lvl        LogLevel
	format     Format
	callerSkip int
	fs         map[string]string

	clock       zapcore.Clock
	internal    zapcore.Core
	errorOutput zapcore.WriteSyncer
	w           io.WriteCloser
}

// 约束输出级别
type levelEnablerFunc func(LogLevel) bool

func (f levelEnablerFunc) Enabled(lvl zapcore.Level) bool {
	return f(LogLevel(lvl))
}

func NewLogger(w io.WriteCloser, opts ...Option) *logger {
	// 按照新的go规范，var变量放括号里~
	var (
		enc zapcore.Encoder
		ws  zapcore.WriteSyncer
		fs  []zapcore.Field
	)
	// 初始化，将日志的错误输出到stderr
	l := &logger{
		color:       DefaultColor,
		stack:       DefaultStack,
		ts:          DefaultTimeLayout,
		lvl:         DefaultLevel,
		format:      DefaultFormat,
		callerSkip:  DefaultCallerSkip,
		fs:          make(map[string]string),
		clock:       zapcore.DefaultClock,
		errorOutput: zapcore.Lock(os.Stderr),
		w:           w,
	}
	for _, f := range opts {
		f(l)
	}
	// 编码各种键和时长，调用者等的格式
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       KeyTime,
		LevelKey:      KeyLevel,
		NameKey:       KeyName,
		CallerKey:     KeyCaller,
		MessageKey:    KeyMessage,
		StacktraceKey: KeyStacktrace,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(l.ts))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	// 指定编码格式
	if l.format == FormatJSON {
		enc = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		if l.color {
			encoderCfg.EncodeLevel = zapcore.LowercaseColorLevelEncoder
			encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(l.ts))
			}
			encoderCfg.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(Yellow.Add(caller.TrimmedPath()))
			}
			enc = zapcore.NewConsoleEncoder(encoderCfg)
		}
	}

	if w == os.Stdout {
		// 因为os.Stdout自带了落盘，所以是个syncer~
		// 锁定是为了 concurrent safe
		ws = zapcore.Lock(os.Stdout)
	} else {
		ws = zapcore.AddSync(w)
	}
	// 约束导出的等级，排除低级别的日志
	l.internal = zapcore.NewCore(enc, ws, levelEnablerFunc(func(lvl LogLevel) bool {
		return lvl >= l.lvl
	}))

	for k, v := range l.fs {
		fs = append(fs, zapcore.Field{
			Key:    k,
			Type:   zapcore.StringType,
			String: v,
		})
	}

	l.internal = l.internal.With(fs)
	return l
}

// 检测entry是否符合
func (l *logger) check(lvl zapcore.Level, msg string, callerSkipOffset ...int) *zapcore.CheckedEntry {
	skip := 2

	if len(callerSkipOffset) > 0 {
		skip = callerSkipOffset[0]
	}

	ent := zapcore.Entry{
		Level:   lvl,
		Time:    l.clock.Now(),
		Message: msg,
	}
	ce := l.internal.Check(ent, nil)
	if ce == nil {
		return nil
	}

	frame, defined := getCallerFrame(l.callerSkip + skip)
	if !defined {
		fmt.Fprintf(l.errorOutput, "%v Logger.check error: failed to get caller\n", ent.Time.UTC())
		l.errorOutput.Sync()
	}

	ce.Entry.Caller = zapcore.EntryCaller{
		Defined:  defined,
		PC:       frame.PC,
		File:     frame.File,
		Line:     frame.Line,
		Function: frame.Function,
	}

	return ce
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if ce := l.check(zapcore.DebugLevel, fmt.Sprintf(format, v...)); ce != nil {
		ce.Write()
	}
}

func (l *logger) Infof(format string, v ...interface{}) {
	if ce := l.check(zapcore.InfoLevel, fmt.Sprintf(format, v...)); ce != nil {
		ce.Write()
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	if ce := l.check(zapcore.WarnLevel, fmt.Sprintf(format, v...)); ce != nil {
		ce.Write()
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	if ce := l.check(zapcore.ErrorLevel, fmt.Sprintf(format, v...)); ce != nil {
		ce.Write()
	}
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	if ce := l.check(zapcore.FatalLevel, fmt.Sprintf(format, v...)); ce != nil {
		ce.Write()
	}
}

func (l *logger) DebugKV(msg string, kvs ...Entry) {
	fields := WrapKV(nil, false, kvs)
	if ce := l.check(zapcore.DebugLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func (l *logger) InfoKV(msg string, kvs ...Entry) {
	fields := WrapKV(nil, false, kvs)
	if ce := l.check(zapcore.InfoLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func (l *logger) WarnKV(msg string, kvs ...Entry) {
	fields := WrapKV(nil, false, kvs)
	if ce := l.check(zapcore.WarnLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func (l *logger) ErrorKV(msg string, kvs ...Entry) {
	fields := WrapKV(nil, false, kvs)
	if ce := l.check(zapcore.ErrorLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func (l *logger) FatalKV(msg string, kvs ...Entry) {
	fields := WrapKV(nil, false, kvs)
	if ce := l.check(zapcore.FatalLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

func (l *logger) DebugErr(msg string, err error, kvs ...Entry) {
	if l.format == FormatText && l.stack {
		fields := WrapKV(nil, l.stack, kvs)
		if ce := l.check(zapcore.DebugLevel, msg); ce != nil {
			ce.Write(fields...)
		}
		fmt.Println(err)
	} else {
		fields := WrapKV(err, l.stack, kvs)
		if ce := l.check(zapcore.DebugLevel, msg); ce != nil {
			ce.Write(fields...)
		}
	}
}

func (l *logger) InfoErr(msg string, err error, kvs ...Entry) {
	if l.format == FormatText && l.stack {
		fields := WrapKV(nil, l.stack, kvs)
		if ce := l.check(zapcore.InfoLevel, msg); ce != nil {
			ce.Write(fields...)
		}
		fmt.Println(err)
	} else {
		fields := WrapKV(err, l.stack, kvs)
		if ce := l.check(zapcore.InfoLevel, msg); ce != nil {
			ce.Write(fields...)
		}
	}
}

func (l *logger) WarnErr(msg string, err error, kvs ...Entry) {
	if l.format == FormatText && l.stack {
		fields := WrapKV(nil, l.stack, kvs)
		if ce := l.check(zapcore.WarnLevel, msg); ce != nil {
			ce.Write(fields...)
		}
		fmt.Println(err)
	} else {
		fields := WrapKV(err, l.stack, kvs)
		if ce := l.check(zapcore.WarnLevel, msg); ce != nil {
			ce.Write(fields...)
		}
	}
}

func (l *logger) ErrorErr(msg string, err error, kvs ...Entry) {
	if l.format == FormatText && l.stack {
		fields := WrapKV(nil, l.stack, kvs)
		if ce := l.check(zapcore.ErrorLevel, msg); ce != nil {
			ce.Write(fields...)
		}
		fmt.Println(err)
	} else {
		fields := WrapKV(err, l.stack, kvs)
		if ce := l.check(zapcore.ErrorLevel, msg); ce != nil {
			ce.Write(fields...)
		}
	}
}

func (l *logger) FatalErr(msg string, err error, kvs ...Entry) {
	if l.format == FormatText && l.stack {
		fields := WrapKV(nil, l.stack, kvs)
		if ce := l.check(zapcore.FatalLevel, msg); ce != nil {
			ce.Write(fields...)
		}
		fmt.Println(err)
	} else {
		fields := WrapKV(err, l.stack, kvs)
		if ce := l.check(zapcore.FatalLevel, msg); ce != nil {
			ce.Write(fields...)
		}
	}
}

func (l *logger) Close() error {
	return l.w.Close()
}

func getCallerFrame(skip int) (frame runtime.Frame, ok bool) {
	const skipOffset = 2 // skip getCallerFrame and Callers

	pc := make([]uintptr, 1)
	numFrames := runtime.Callers(skip+skipOffset, pc)
	if numFrames < 1 {
		return
	}

	frame, _ = runtime.CallersFrames(pc).Next()
	return frame, frame.PC != 0
}

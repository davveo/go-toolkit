package logger

import (
	"io"
)

type Logger interface {
	FormatLogger
	ErrorLogger
	ExtendedLogger
	io.Closer
}

type FormatLogger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

type ExtendedLogger interface {
	DebugKV(msg string, kvs ...Entry)
	InfoKV(msg string, kvs ...Entry)
	WarnKV(msg string, kvs ...Entry)
	ErrorKV(msg string, kvs ...Entry)
	FatalKV(msg string, kvs ...Entry)
}

type ErrorLogger interface {
	DebugErr(msg string, err error, kvs ...Entry)
	InfoErr(msg string, err error, kvs ...Entry)
	WarnErr(msg string, err error, kvs ...Entry)
	ErrorErr(msg string, err error, kvs ...Entry)
	FatalErr(msg string, err error, kvs ...Entry)
}

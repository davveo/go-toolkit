package logger

type Option func(*logger)

func WithLevel(lvl LogLevel) Option {
	return func(l *logger) {
		l.lvl = lvl
	}
}

func WithColor(color bool) Option {
	return func(l *logger) {
		l.color = color
	}
}

func WithTimestampLayout(layout string) Option {
	return func(l *logger) {
		l.ts = layout
	}
}

func WithForamt(format Format) Option {
	return func(l *logger) {
		l.format = format
	}
}

func WithStack(stack bool) Option {
	return func(l *logger) {
		l.stack = stack
	}
}

func WithField(key, value string) Option {
	return func(l *logger) {
		l.fs[key] = value
	}
}

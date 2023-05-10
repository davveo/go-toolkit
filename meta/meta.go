package meta

import (
	"github.com/davveo/go-toolkit/env"
	"github.com/google/uuid"
)

const (
	defaultPlatform = "test"
	defaultService  = "test"
	defaultEnv      = env.EnvDev
	defaultVersion  = "1.0.0"
	defaultLogPath  = "/tmp/logs/"
)

type (
	Meta interface {
		// Platform 业务线
		Platform() string
		// Service 服务名
		Service() string
		// Env 运行环境: dev/test/uat/prod
		Env() env.AppEnv
		// Version 版本号
		Version() string
		// LogPath 日志路径
		LogPath() string
		// ID id
		ID() string
	}
	// MetaEnv 应用元信息
	MetaEnv struct {
		// 应用所属业务平台
		platform string
		// 应用所属服务
		service string
		// 运行环境: dev/test/uat/prod
		env env.AppEnv
		// 版本号
		version string
		// 日志地址
		logPath string
		//id
		id string
	}
	Option func(metaEnv *MetaEnv)
)

func NewMetaEnv(opts ...Option) *MetaEnv {
	m := &MetaEnv{
		service:  defaultService,
		env:      defaultEnv,
		platform: defaultPlatform,
		id:       uuid.New().String(),
		version:  defaultVersion,
		logPath:  defaultLogPath,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *MetaEnv) Platform() string {
	return m.platform
}

func (m *MetaEnv) Service() string {
	return m.service
}

func (m *MetaEnv) Env() env.AppEnv {
	return m.env
}

func (m *MetaEnv) Version() string {
	return m.version
}

func (m *MetaEnv) LogPath() string {
	return m.logPath
}

func (m *MetaEnv) ID() string {
	return m.id
}

func Platform(platform string) Option {
	return func(metaEnv *MetaEnv) {
		metaEnv.platform = platform
	}
}

func ID(id string) Option {
	return func(metaEnv *MetaEnv) {
		metaEnv.id = id
	}
}

func LogPath(logPath string) Option {
	return func(metaEnv *MetaEnv) {
		metaEnv.logPath = logPath
	}
}

func Version(version string) Option {
	return func(metaEnv *MetaEnv) {
		metaEnv.version = version
	}
}

func Env(env env.AppEnv) Option {
	return func(metaEnv *MetaEnv) {
		metaEnv.env = env
	}
}

func Service(service string) Option {
	return func(metaEnv *MetaEnv) {
		metaEnv.service = service
	}
}

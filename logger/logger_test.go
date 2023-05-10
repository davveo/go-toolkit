package logger

import (
	"fmt"
	"github.com/davveo/go-toolkit/env"
	"github.com/davveo/go-toolkit/meta"
	"strconv"
	"testing"
)

func TestLogger(t *testing.T) {
	service := "lemon"
	metaEnv := meta.NewMetaEnv(
		meta.Env(env.EnvTest),
		meta.Service(service),
		meta.Platform("platform"),
		meta.LogPath(fmt.Sprintf("/tmp/%s/", service)))

	if err := InitLogger(metaEnv); err != nil {
		panic(err)
	}
	defer Close()
	Infof("服务器运行中...  端口: " + strconv.Itoa(111))
}

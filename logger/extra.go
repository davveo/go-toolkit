package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 暴露额外kv对，方便搜索，参考trace的经验
// key必须要是string，方便查找，值因为是业务传入的，可以不一样。
// Expose extra KV pair for search which we have been learn from trace design experience
// Keys must be in form of string so that we can search with the specific keys, but value can be any form.
type Entry interface {
	Key() string
	Value() interface{}
	Entry()
}

type kv struct {
	k string
	v interface{}
}

func (m *kv) Key() string {
	return m.k
}

func (m *kv) Value() interface{} {
	return m.v
}

func (m *kv) Entry() {}

func KV(key string, value interface{}) Entry {
	return &kv{k: key, v: value}
}

func WrapKV(err error, stack bool, kvs []Entry) []zapcore.Field {
	var fields []zapcore.Field

	fieldsSz := len(kvs) + 1
	if err != nil {
		fieldsSz++
		if stack {
			fieldsSz++
		}
	}

	fields = make([]zap.Field, 0, fieldsSz)

	// Inject error to our log
	if err != nil {
		fields = append(fields, zap.Error(err))
		if stack {
			fields = append(fields, zap.StackSkip(KeyStacktrace, 2))
		}
	}

	if len(kvs) > 0 {
		// we can directly use namespace from zap.
		fields = append(fields, zap.Namespace(KeyBizName))
		for _, ent := range kvs {
			fields = append(fields, zap.Any(ent.Key(), ent.Value()))
		}
	}

	return fields
}

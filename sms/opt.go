package sms

type InitOptions struct {
	AccessId  string
	AccessKey string
	Sign      string
	Template  string
	Extra     []string
}

type InitOption func(options *InitOptions)

func WithAccessId(accessId string) InitOption {
	return func(options *InitOptions) {
		options.AccessId = accessId
	}
}

func WithAccessKey(accessKey string) InitOption {
	return func(options *InitOptions) {
		options.AccessKey = accessKey
	}
}

func WithSign(sign string) InitOption {
	return func(options *InitOptions) {
		options.Sign = sign
	}
}

func WithTemplate(template string) InitOption {
	return func(options *InitOptions) {
		options.Template = template
	}
}

func WithExtra(extra ...string) InitOption {
	return func(options *InitOptions) {
		options.Extra = extra
	}
}

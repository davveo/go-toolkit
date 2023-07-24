package sms

import (
	"fmt"
	"github.com/davveo/go-toolkit/sms/aliyun"
	"github.com/davveo/go-toolkit/sms/huawei"
	"github.com/davveo/go-toolkit/sms/huyi"
	"github.com/davveo/go-toolkit/sms/smsbao"
	"github.com/davveo/go-toolkit/sms/submail"
	"github.com/davveo/go-toolkit/sms/twilio"
	"github.com/davveo/go-toolkit/sms/tx"
	"github.com/davveo/go-toolkit/sms/volcengine"
)

type SmsClient interface {
	SendMessage(param map[string]string, targetPhoneNumber ...string) error
}

func NewSms(provider phoneChannel, options ...InitOption) (SmsClient, error) {
	switch provider {
	case AliYun:
		return aliyun.NewAliYunClient(options...)
	case TencentCloud:
		return tx.NewTencentClient(options...)
	case VolcEngine:
		return volcengine.NewVolcClient(options...)
	case HuYi:
		return huyi.NewHuYiClient(options...)
	case HuaweiCloud:
		return huawei.NewHuaweiClient(options...)
	case Twilio:
		return twilio.NewTwilioClient(options...)
	case SmsBao:
		return smsbao.NewSmsBaoClient(options...)
	case SubMail:
		return submail.NewSubMailClient(options...)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

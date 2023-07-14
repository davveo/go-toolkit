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

func NewSmsClient(provider phoneChannel, accessId string,
	accessKey string, sign string, template string, other ...string) (SmsClient, error) {
	switch provider {
	case AliYun:
		return aliyun.GetAliyunClient(accessId, accessKey, sign, template)
	case TencentCloud:
		return tx.GetTencentClient(accessId, accessKey, sign, template, other)
	case VolcEngine:
		return volcengine.GetVolcClient(accessId, accessKey, sign, template, other)
	case HuYi:
		return huyi.GetHuyiClient(accessId, accessKey, template)
	case HuaweiCloud:
		return huawei.GetHuaweiClient(accessId, accessKey, sign, template, other)
	case Twilio:
		return twilio.GetTwilioClient(accessId, accessKey, template)
	case SmsBao:
		return smsbao.GetSmsbaoClient(accessId, accessKey, sign, template, other)
	case SubMail:
		return submail.GetSubmailClient(accessId, accessKey, template)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

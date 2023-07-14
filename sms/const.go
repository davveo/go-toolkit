package sms

type phoneChannel string

const (
	AliYun       phoneChannel = "Aliyun SMS"        // 阿里云短信
	TencentCloud phoneChannel = "Tencent Cloud SMS" // 腾讯短信
	VolcEngine   phoneChannel = "Volc Engine SMS"   // 火山短信
	HuYi         phoneChannel = "Huyi SMS"          // 互亿短信
	HuaweiCloud  phoneChannel = "Huawei Cloud SMS"  // 华为短信
	Twilio       phoneChannel = "Twilio SMS"
	SmsBao       phoneChannel = "SmsBao SMS" // 短信宝
	SubMail      phoneChannel = "SUBMAIL SMS"
)

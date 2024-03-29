package aliyun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davveo/go-toolkit/sms"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

const region = "cn-hangzhou"

type AliyunClient struct {
	template string
	sign     string
	core     *dysmsapi.Client
}

type AliyunResult struct {
	RequestId string
	Message   string
}

/*
	NewAliYunClient(
		WithAccessId(""),
		WithAccessKey(),
		WithSign(),
		WithTemplate()
	)
*/

func NewAliYunClient(options ...sms.InitOption) (*AliyunClient, error) {
	opts := &sms.InitOptions{}
	for _, option := range options {
		option(opts)
	}
	if opts.AccessKey == "" || opts.AccessId == "" ||
		opts.Template == "" || opts.Sign == "" {
		return nil, fmt.Errorf("参数有误")
	}
	client, err := dysmsapi.NewClientWithAccessKey(
		region, opts.AccessId, opts.AccessKey)
	if err != nil {
		return nil, err
	}
	aliYunClient := &AliyunClient{
		template: opts.Template,
		core:     client,
		sign:     opts.Sign,
	}

	return aliYunClient, nil
}

func (c *AliyunClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	requestParam, err := json.Marshal(param)
	if err != nil {
		return err
	}

	if len(targetPhoneNumber) < 1 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	phoneNumbers := bytes.Buffer{}
	phoneNumbers.WriteString(targetPhoneNumber[0])
	for _, s := range targetPhoneNumber[1:] {
		phoneNumbers.WriteString(",")
		phoneNumbers.WriteString(s)
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phoneNumbers.String()
	request.TemplateCode = c.template
	request.TemplateParam = string(requestParam)
	request.SignName = c.sign

	response, err := c.core.SendSms(request)
	if response.Code != "OK" {
		aliyunResult := AliyunResult{}
		json.Unmarshal(response.GetHttpContentBytes(), &aliyunResult)
		return fmt.Errorf(aliyunResult.Message)
	}
	return err
}

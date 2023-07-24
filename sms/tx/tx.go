package tx

import (
	"fmt"
	"github.com/davveo/go-toolkit/sms"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	_sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"strconv"
)

type TencentClient struct {
	core     *_sms.Client
	appId    string
	sign     string
	template string
}

func NewTencentClient(options ...sms.InitOption) (*TencentClient, error) {
	opts := &sms.InitOptions{}
	for _, option := range options {
		option(opts)
	}

	if len(opts.Extra) < 1 {
		return nil, fmt.Errorf("missing parameter: appId")
	}

	credential := common.NewCredential(opts.AccessId, opts.AccessKey)
	config := profile.NewClientProfile()
	config.HttpProfile.ReqMethod = "POST"

	region := "ap-guangzhou"
	client, err := _sms.NewClient(credential, region, config)
	if err != nil {
		return nil, err
	}

	tencentClient := &TencentClient{
		core:     client,
		sign:     opts.Sign,
		template: opts.Template,
		appId:    opts.Extra[0],
	}

	return tencentClient, nil
}

func (c *TencentClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	var paramArray []string
	index := 0
	for {
		value := param[strconv.Itoa(index)]
		if len(value) == 0 {
			break
		}
		paramArray = append(paramArray, value)
		index++
	}

	request := _sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(c.appId)
	request.SignName = common.StringPtr(c.sign)
	request.TemplateParamSet = common.StringPtrs(paramArray)
	request.TemplateId = common.StringPtr(c.template)
	request.PhoneNumberSet = common.StringPtrs(targetPhoneNumber)

	_, err := c.core.SendSms(request)
	return err
}

package volcengine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davveo/go-toolkit/sms"
	_sms "github.com/volcengine/volc-sdk-golang/service/sms"
)

type VolcClient struct {
	core       *_sms.SMS
	sign       string
	template   string
	smsAccount string
}

func NewVolcClient(options ...sms.InitOption) (*VolcClient, error) {
	opts := &sms.InitOptions{}
	for _, option := range options {
		option(opts)
	}

	if len(opts.Extra) < 1 {
		return nil, fmt.Errorf("missing parameter: smsAccount")
	}

	client := _sms.NewInstance()
	client.Client.SetAccessKey(opts.AccessId)
	client.Client.SetSecretKey(opts.AccessKey)

	volcClient := &VolcClient{
		core:       client,
		sign:       opts.Sign,
		template:   opts.Template,
		smsAccount: opts.Extra[0],
	}

	return volcClient, nil
}

func (c *VolcClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
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

	req := &_sms.SmsRequest{
		SmsAccount:    c.smsAccount,
		Sign:          c.sign,
		TemplateID:    c.template,
		TemplateParam: string(requestParam),
		PhoneNumbers:  phoneNumbers.String(),
	}
	_, _, err = c.core.Send(req)
	return err
}

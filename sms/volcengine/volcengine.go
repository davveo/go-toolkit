package volcengine

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/volcengine/volc-sdk-golang/service/sms"
)

type VolcClient struct {
	core       *sms.SMS
	sign       string
	template   string
	smsAccount string
}

func GetVolcClient(accessId, accessKey, sign, templateId string, smsAccount []string) (*VolcClient, error) {
	if len(smsAccount) < 1 {
		return nil, fmt.Errorf("missing parameter: smsAccount")
	}

	client := sms.NewInstance()
	client.Client.SetAccessKey(accessId)
	client.Client.SetSecretKey(accessKey)

	volcClient := &VolcClient{
		core:       client,
		sign:       sign,
		template:   templateId,
		smsAccount: smsAccount[0],
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

	req := &sms.SmsRequest{
		SmsAccount:    c.smsAccount,
		Sign:          c.sign,
		TemplateID:    c.template,
		TemplateParam: string(requestParam),
		PhoneNumbers:  phoneNumbers.String(),
	}
	_, _, err = c.core.Send(req)
	return err
}

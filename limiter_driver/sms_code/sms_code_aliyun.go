package sms_code

import (
	"context"
	"fmt"

	dysmsclient "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	"github.com/thinkgos/proc-extra/limiter/limit_verified"
)

var _ limit_verified.LimitVerifiedProvider = &SmsAliyun{}

type SmsAliyun struct {
	Client       *dysmsclient.Client // 阿里云短信服务客户端
	SignName     string              // 短信签名名称
	TemplateCode string              // 短信模板代号, 如 SMS_XXXX
	Template     string              // 短信模块, 如 {"code":"%s"}
}

func (*SmsAliyun) Name() string { return "sms-aliyun" }

func (s *SmsAliyun) SendCode(ctx context.Context, target, code string) error {
	request := &dysmsclient.SendSmsRequest{
		PhoneNumbers:  new(target),
		SignName:      new(s.SignName),
		TemplateCode:  new(s.TemplateCode),
		TemplateParam: new(fmt.Sprintf(s.Template, code)),
	}
	response, err := s.Client.SendSms(request)
	if err != nil {
		return err
	}

	// 默认流控：使用同一个签名，对同一个手机号码发送短信验证码，
	// 支持1条/分钟，5条/小时 ，累计10条/天
	if *response.Body.Code == "isv.BUSINESS_LIMIT_CONTROL" {
		return limit_verified.ErrReachMaximumQuota
	}
	if *response.Body.Code != "OK" {
		return fmt.Errorf("发送失败, %s", *response.Body.Code)
	}
	return nil
}

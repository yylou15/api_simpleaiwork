package mail

import (
	"encoding/json"
	"fmt"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm20151123 "github.com/alibabacloud-go/dm-20151123/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
)

var client *dm20151123.Client

// CreateEmailClient 初始化账号Client
func CreateEmailClient() (*dm20151123.Client, error) {
	// 工程代码建议使用更安全的无AK方式，凭据配置方式请参见： `https://help.aliyun.com/document_detail/378661.html。`
	cred, err := credential.NewCredential(nil)
	if err != nil {
		return nil, err
	}

	config := &openapi.Config{
		Credential: cred,
	}
	// Endpoint 请参考 `https://api.aliyun.com/product/Dm`
	config.Endpoint = tea.String("dm.aliyuncs.com")

	client, err = dm20151123.NewClient(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// SendEmailCode 发送邮箱验证码
// toAddress: 接收邮箱
// code: 验证码
// subject: 邮件主题 (可选, 默认为 "SayRight Verify Code")
func SendEmailCode(toAddress, code string, subject string) error {
	if subject == "" {
		subject = "SayRight Verify Code"
	}

	textBody := fmt.Sprintf("Your Verify Code Is %s", code)

	singleSendMailRequest := &dm20151123.SingleSendMailRequest{
		AccountName:    tea.String("no-reply@mail.simpleaiwork.com"),
		AddressType:    tea.Int32(1),
		ToAddress:      tea.String(toAddress),
		Subject:        tea.String(subject),
		TextBody:       tea.String(textBody),
		ReplyToAddress: tea.Bool(false),
	}

	runtime := &util.RuntimeOptions{}

	// 使用 recover 捕获可能的 panic
	var tryErr error
	func() {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				tryErr = r
			}
		}()
		_, err := client.SingleSendMailWithOptions(singleSendMailRequest, runtime)
		if err != nil {
			tryErr = err
		}
	}()

	if tryErr != nil {
		var sdkError *tea.SDKError
		if t, ok := tryErr.(*tea.SDKError); ok {
			sdkError = t
		} else {
			sdkError = &tea.SDKError{
				Message: tea.String(tryErr.Error()),
			}
		}

		// 尝试解析错误详情
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(sdkError.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			if recommend, ok := m["Recommend"]; ok {
				return fmt.Errorf("sdk error: %s, recommend: %v", tea.StringValue(sdkError.Message), recommend)
			}
		}
		return fmt.Errorf("sdk error: %s", tea.StringValue(sdkError.Message))
	}

	return nil
}

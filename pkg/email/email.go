package email

import (
	"bytes"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

// Email 结构体用于发送邮件
type Email struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
}

// NewEmail 创建一个新的Email实例
func NewEmail(smtpHost string, smtpPort int, username, password, from string) *Email {
	return &Email{
		SMTPHost: smtpHost,
		SMTPPort: smtpPort,
		Username: username,
		Password: password,
		From:     from,
	}
}

// Send 方法用于发送邮件
func (e *Email) Send(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(e.SMTPHost, e.SMTPPort, e.Username, e.Password)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	return nil
}

func (e *Email) SendVerificationCode(to string, code string, emailTemplate string) error {
	// 解析 HTML 模板
	tmpl, err := template.New("verificationEmail").Parse(emailTemplate)
	if err != nil {
		return fmt.Errorf("解析模板失败: %v", err)
	}

	// 准备模板数据
	data := struct {
		VerificationCode string
	}{
		VerificationCode: code,
	}

	// 执行模板
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("执行模板失败: %v", err)
	}

	// 发送邮件
	subject := "Email Verification"
	err = e.Send([]string{to}, subject, body.String())
	if err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	return nil
}

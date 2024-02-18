package mail

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

type mail struct {
	host     string
	port     string
	username string
	password string
	subject  string
	from     string
	to       []string
	cc       string
}

type Option func(mail *mail)

func WithHost(host string) Option {
	return func(mail *mail) {
		mail.host = host
	}
}

func WithPort(port string) Option {
	return func(mail *mail) {
		mail.port = port
	}
}

func WithUsername(username string) Option {
	return func(mail *mail) {
		mail.username = username
	}
}

func WithPassword(password string) Option {
	return func(mail *mail) {
		mail.password = password
	}
}

func WithSubject(subject string) Option {
	return func(mail *mail) {
		mail.subject = subject
	}
}

func WithFrom(from string) Option {
	return func(mail *mail) {
		mail.from = from
	}
}

// New 实例化
func New(options ...Option) *mail {
	var (
		host     = "smtp.qq.com"
		port     = "25"
		username = "735273025@qq.com"
		password = "gtxknrzhmctlbbfb"
		from     = "735273025@qq.com"
		mail     = &mail{
			host:     host,
			port:     port,
			username: username,
			password: password,
			from:     from,
		}
	)
	for _, option := range options {
		option(mail)
	}
	return mail
}

// Subject 设置邮件标题
func (mail *mail) Subject(subject string) *mail {
	mail.subject = subject
	return mail
}

// To 设置邮件接收人
func (mail *mail) To(to string) *mail {
	mail.to = []string{to}
	return mail
}

// Cc 设置抄送人
func (mail *mail) Cc(cc string) *mail {
	mail.cc = cc
	return mail
}

// Send 发送
func (mail *mail) Send(content string) (err error) {
	if content == "" {
		err = errors.New("please confirm the email content")
		return
	}
	if len(mail.to) < 1 {
		err = errors.New("please confirm who the email is sent to")
		return
	}

	header := make(map[string]string)
	if mail.subject != "" {
		header["Subject"] = mail.subject
	}
	header["To"] = mail.to[0]
	header["From"] = mail.from
	if mail.cc != "" {
		header["Cc"] = mail.cc
	}
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"
	var buf bytes.Buffer
	for k, v := range header {
		buf.WriteString(k)
		buf.WriteString(":")
		buf.WriteString(v)
		buf.WriteString("\r\n")
	}
	buf.WriteString(base64.StdEncoding.EncodeToString([]byte(content)))

	message := buf.Bytes()

	auth := smtp.PlainAuth("", mail.username, mail.password, mail.host)

	host := fmt.Sprintf("%s:%s", mail.host, mail.port)

	err = smtp.SendMail(host, auth, mail.from, mail.to, message)
	if err != nil {
		err = fmt.Errorf("failed to send mail to %s, error: %s", strings.Join(mail.to, ","), err)
	}
	return
}

package mail

import (
	"bytes"
	"encoding/json"
	"html/template"
	"strconv"
	"testing"
)

func TestSendEmail(t *testing.T) {
	var (
		buf  bytes.Buffer
		jbuf = bytes.NewBuffer([]byte{})
	)
	j := map[string]string{
		"status":  "success",
		"message": "失败&",
	}
	tpl, _ := template.New("email").Parse(`
	<div>任务：{{.Name}}</div>
	<div>ID：{{.ID}}</div>
	<div>返回：<pre>{{.Response}}</pre></div>
	<div>状态：{{.Status}}</div>
	<div>提示：{{.Error}}</div>
	`)
	encoder := json.NewEncoder(jbuf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	encoder.Encode(j)
	tpl.Execute(&buf, map[string]string{
		"ID":       "avaddsr13231",
		"Name":     "任务1",
		"Response": jbuf.String(),
		"Status":   strconv.Itoa(int(1)),
		"Error":    "执行失败",
	})
	err := New().Subject("测试邮件").To("735273025@qq.com").Send(buf.String())
	if err != nil {
		t.Fatal("发送失败", err)
	}
	t.Log("发送成功")
}

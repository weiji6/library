package tool

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
)

func sendEmailBYQQEmailAndFormat(to string) error {
	from := "2954585173@qq.com"
	password := "ecqrpsjppazedgbf" // 邮箱授权码
	smtpServer := "smtp.qq.com:465"
	// 读取图片
	imgPath := "./image1.png"
	imgData, err := ioutil.ReadFile(imgPath)
	if err != nil {
		log.Fatalf("无法读取图片: %v", err)
	}
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)
	// 邮件内容
	body := `
                <h1>这是一级标题</h1>
                <h2>这是二级标题</h2>
                <p>这是 <strong>` + `加粗` + `</strong></p>
                <p>这是 <em>` + `斜体` + `</em></p>
                <p>这是 <u>` + `下划线` + `</u></p>
                <p>这是 <s>` + `删除线` + `</s></p>
                <p>下面是一张图片</p>
                <img src="cid:image001" alt="chart" width="180" height="180">         `
	// 邮件头部
	header := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      "这是一封测试邮件",
		"MIME-Version": "1.0",
		"Content-Type": `multipart/related; boundary="BOUNDARY"`,
	}
	var message bytes.Buffer
	// 添加头部
	for k, v := range header {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	// 添加正文部分
	message.WriteString("--BOUNDARY\r\n")
	message.WriteString(
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
	message.WriteString(body + "\r\n")
	// 添加图片部分
	message.WriteString("--BOUNDARY\r\n")
	message.WriteString("Content-Type: chart/png\r\n")
	message.WriteString("Content-Transfer-Encoding: base64\r\n")
	message.WriteString("Content-ID: <image001>\r\n\r\n")
	message.WriteString(imgBase64 + "\r\n")
	// MIME 结束
	message.WriteString("--BOUNDARY--")
	// 设置 PlainAuth
	// 第一个 "" 可以看作一个可选参数，多数情况下不需要设置，传空即可。
	// 它的存在是为了满足 SMTP 标准协议中的扩展需求，但实际应用中很少需要自定义。
	auth := smtp.PlainAuth("", from, password, "smtp.qq.com")
	// 创建 tls 配置
	// InsecureSkipVerify: true：表示跳过对服务器证书的验证。这在生产环境中是不安全的，通常只在开发或测试环境中使用。
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "smtp.qq.com",
	}
	// 连接到 SMTP 服务器
	conn, err := tls.Dial("tcp", smtpServer, tlsconfig)
	if err != nil {
		return fmt.Errorf("TLS 连接失败: %v", err)
	}
	defer conn.Close()
	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, "smtp.qq.com")
	if err != nil {
		return fmt.Errorf("SMTP 客户端创建失败: %v", err)
	}
	defer client.Quit()
	// 使用 auth 进行认证
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("认证失败: %v", err)
	}
	// 设置发件人和收件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("发件人设置失败: %v", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("收件人设置失败: %v", err)
	}
	// 写入邮件内容
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("数据写入失败: %v", err)
	}
	defer wc.Close()
	// 发送邮件
	_, err = wc.Write(message.Bytes())
	if err != nil {
		return fmt.Errorf("消息发送失败: %v", err)
	}
	return nil
}

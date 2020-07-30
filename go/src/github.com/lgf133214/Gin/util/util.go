package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/lgf133214/Gin/config"
	"gopkg.in/gomail.v2"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	pattern = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg     = regexp.MustCompile(pattern)
)

func GetMD5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

func GetSession(s string) string {
	now := time.Now().Unix()
	return GetMD5(s + strconv.FormatInt(now, 10))
}

func GetAddr(remoteAddr string) string {
	return strings.Split(remoteAddr, ":")[0]
}

func SendMail(mailTo []string, subject string, body string) error {
	mailConn := map[string]string{
		"user": config.From,
		"pass": config.Password,
		"host": config.Host,
		"port": config.Port,
	}

	port, _ := strconv.Atoi(mailConn["port"])

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(mailConn["user"], "contact"))
	m.SetHeader("To", mailTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err
}

func SendVerifyCode(code string, register bool, mail string) bool {
	subject := "验证链接"
	body := "感谢您的使用，请点击链接完成操作："
	if register {
		url := "http://" + config.Domain + "/activation/" + code
		err := SendMail([]string{mail}, subject, body+url)
		if err != nil {
			log.Println(err)
			return false
		}
	} else {
		url := "http://" + config.Domain + "/reset/pw/" + code
		err := SendMail([]string{mail}, subject, body+url)
		if err != nil {
			log.Println(err)
			return false
		}
	}
	return true
}

func SendMessage(name, msg, mail string) bool {
	subject := "博客反馈信息"
	m := mail + "\n" + name + "\n" + msg
	err := SendMail([]string{config.From}, subject, m)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func VerifyEmailFormat(email string) bool {
	return reg.MatchString(email)
}

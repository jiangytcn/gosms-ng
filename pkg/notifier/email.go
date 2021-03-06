/**
 * File              : email.go
 * Author            : Jiang Yitao <jiangyt.cn#gmail.com>
 * Date              : 17.08.2019
 * Last Modified Date: 20.09.2019
 * Last Modified By  : Jiang Yitao <jiangyt.cn#gmail.com>
 */
package notifier

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	gomail "gopkg.in/mail.v2"

	"github.com/jiangytcn/gosms-ng/config"
)

type Message struct {
	To           []string
	From         string
	Subject      string
	Body         string
	Info         string
	EmbededFiles []string
	SMTPCfg      config.SMTPConfig
}

func (m Message) WebRequestSend() error {
	return fmt.Errorf("not implemented")
}

func (msg Message) GenericSend() error {

	var failedReceiver []string

	dialer, err := createDialer(msg.SMTPCfg)
	if err != nil {
		return err
	}

	var num int
	for _, address := range msg.To {
		m := gomail.NewMessage()
		m.SetHeader("From", fmt.Sprintf("%s <%s>", msg.SMTPCfg.FromName, msg.SMTPCfg.FromAddress))
		m.SetHeader("To", address)
		m.SetHeader("Subject", msg.Subject)
		for _, file := range msg.EmbededFiles {
			m.Embed(file)
		}

		m.SetBody("text/html", msg.Body)

		e := dialer.DialAndSend(m)
		if e != nil {

			err = e
			failedReceiver = append(failedReceiver, address)
			continue
		}

		time.Sleep(10 * time.Second)

		num++
	}

	if len(msg.To) != num {
		return fmt.Errorf("faile to delver to %s", strings.Join(failedReceiver, ","))
	}

	return nil
}

func createDialer(cfg config.SMTPConfig) (*gomail.Dialer, error) {

	tlsconfig := &tls.Config{
		ServerName: cfg.Server,
	}

	d := gomail.NewDialer(cfg.Server, cfg.Port, cfg.User, cfg.Password)
	d.TLSConfig = tlsconfig

	return d, nil
}

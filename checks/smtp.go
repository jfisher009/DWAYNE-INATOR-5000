package checks

import (
	"bytes"
	"net/smtp"
	"strconv"
)

type Smtp struct {
	checkBase
	Sender    string
	Receiver  string
	Body      string
	Encrypted bool
}

type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

func (c Smtp) Run(teamID uint, boxIp string, res chan Result) {
	m, err := smtp.Dial(boxIp + ":" + int64(c.Port),10)

	if err != nil {
		res <- Result{
			Error: "error sending email",
			Debug: err.Error(),
		}
		return
	}

	defer m.Close()
	m.Mail(c.Sender)
	m.Rcpt(c.Receiver)
	wc, err := m.Data()
	if err != nil {
		res <- Result{
			Error: "error sending emaiL",
			Debug: err.Error(),
		}
		return
	}
	defer wc.Close()
	buf := bytes.NewBufferString(c.Body)
	if _, err = buf.WriteTo(wc); err != nil {
		res <- Result {
			Error: "error sending email",
			Debug: err.Error(),
		}
	}

	res <- Result {
		Status: true,
		Debug:  "successfully sent email",
	}
	return
}

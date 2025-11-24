package mailer

import (
	"bytes"
	"embed"
	"time"

	ht "html/template"
	tt "text/template"

	"github.com/wneessen/go-mail"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	client *mail.Client
	sender string
}

func New(host string, port int, username, password, sender string) (*Mailer, error) {
	// creating a new client instance
	client, err := mail.NewClient(
		host,
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	mailer := &Mailer{
		client: client,
		sender: sender,
	}

	return mailer, nil
}

func (m *Mailer) Send(recipient string, templateFile string, data any) error {
	// Parsing templateFile from the in memory templates/ folder
	textTmpl, err := tt.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// A buffer for storing subject
	subject := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	// fetching the html template from templates/
	htmlTmpl, err := ht.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// A buffer for storing htmlBody template
	htmlBody := new(bytes.Buffer)
	err = htmlTmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	// Creating a new message instance
	msg := mail.NewMsg()

	err = msg.To(recipient)
	if err != nil {
		return err
	}

	err = msg.From(m.sender)
	if err != nil {
		return err
	}

	msg.Subject(subject.String())
	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	msg.AddAlternativeString(mail.TypeTextHTML, htmlBody.String())

	for i := 1; i <= 3; i++ {
		err = m.client.DialAndSend(msg)
		if err == nil {
			return nil
		}
		// If it didn't work, sleep for a short time and retry.
		if i != 3 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return err
}

package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

// Mailer contains a mail.Dialer instance (used to connect to an
// SMTP server) and the sender information for emails.
type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	// Initialize a new mail.Dialer instance with the given SMTP server settings.
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second
	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

// Send takes the recipient email address as the first parameter, the name of the
// file containing the templates, and any dynamic data for the templates.
func (m Mailer) Send(recipient, templateFile string, data interface{}) error {
	// Parse the required template file from the embedded file system.
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Execute the named template "subject".
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// Execute the named template "plainBody".
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	// Execute the named template "htmlBody".
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	// Create a mail message.
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	// Open a connection to the SMTP server, send the message, abd close the connection.
	// If there is a timeout, it will return a "dial tcp: i/o timeout" error.
	for i := 1; i <= 3; i++ {
		err = m.dialer.DialAndSend(msg)
		// If everything worked, return nil.
		if nil == err {
			return nil
		}
		// If it didn't work, sleep for a short time and retry.
		time.Sleep(500 * time.Millisecond)
	}

	return err
}

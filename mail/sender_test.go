package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xiaowuzai/simplebank/util"
)

func TestGmailSender(t *testing.T) {
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1> Hello world</h1>
	<p> This is a test message</p>
	`

	to := []string{"z158834522@163.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}

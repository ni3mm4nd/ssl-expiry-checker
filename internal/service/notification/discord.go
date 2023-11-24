package notification

import (
	"fmt"
	"log"
	"strings"

	"github.com/ni3mm4nd/ssl-expiry-checker/internal/domain/sslcheck"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/service/notification/discord"
)

type discordNotification struct {
	service DiscordConfig
}

type webhook string

type DiscordConfig struct {
	Enabled  bool    `yaml:"enabled"`
	Webhook  webhook `yaml:"webhook"`
	Username string  `yaml:"username"`
}

func (webhook) MarshalYAML() (interface{}, error) {
	return "******", nil
}

func (w webhook) String() string {
	return string(w)
}

func (d discordNotification) Notify(checks []sslcheck.SSLCheck) {
	log.Println("Discord notification in progress...")
	var s strings.Builder

	s.WriteString("WARNING: SSL certificates are about to expire!\n")
	for _, c := range checks {
		s.WriteString(fmt.Sprintf("\n\nurl: %s\ndays left: %v days\n", c.TargetURL, c.DaysLeft))
		if c.Error != "" {
			s.WriteString(fmt.Sprintf("error: %s\n", c.Error))
		}
	}

	content := s.String()
	err := discord.SendMessage(d.service.Webhook.String(), discord.Message{
		Username: &d.service.Username,
		Content:  &content,
	})
	if err != nil {
		log.Println(err)
	}

	log.Println("Discord notification sent")
}

func NewDiscordNotification(config DiscordConfig) INotification {
	return &discordNotification{
		service: config,
	}
}

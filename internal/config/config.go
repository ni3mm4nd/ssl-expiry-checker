package config

import (
	"log"
	"os"

	"github.com/ni3mm4nd/ssl-expiry-checker/internal/repository"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/repository/sslcheckrepo"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/service"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/service/notification"
	"gopkg.in/yaml.v2"
)

const (
	CONFIG_FILE_NAME = "ssl-expiry-checker.conf"
)

var conf *Config

type Config struct {
	NotificationsConfig  NotificationsConfig `yaml:"notifications"`
	ALERTDaysLeft        int                 `yaml:"alert.daysleft"`
	RescanCronExpression string              `yaml:"rescan.cron"`
}

type NotificationsConfig struct {
	SMTPConfig    notification.SMTPConfig    `yaml:"smtp"`
	DiscordConfig notification.DiscordConfig `yaml:"discord"`
}

func Get() *Config {
	if conf == nil {
		conf = &Config{}
	}
	return conf
}

func (c *Config) loadConf() *Config {

	yamlFile, err := os.ReadFile(CONFIG_FILE_NAME)
	if err != nil {
		log.Printf("yamlFile.Get err: %v\n", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func Init() {
	var config Config
	config.loadConf()

	repository.New()
	repository.GetRepos().SSLCheckRepo = sslcheckrepo.NewFileStore()

	if config.NotificationsConfig.SMTPConfig.Enabled {
		notification.Add(
			notification.NewEmailNotification(config.NotificationsConfig.SMTPConfig),
		)
	}

	if config.NotificationsConfig.DiscordConfig.Enabled {
		notification.Add(
			notification.NewDiscordNotification(config.NotificationsConfig.DiscordConfig),
		)
	}

	if config.ALERTDaysLeft == 0 {
		config.ALERTDaysLeft = 10
	}

	err := service.NewScheduler(config.RescanCronExpression, config.ALERTDaysLeft)
	if err != nil {
		log.Fatal(err)
	}

	service.PrintScheduledJobs()

	conf = &config
}

func (c *Config) GetFullYamlConfig() ([]byte, error) {
	return yaml.Marshal(c)
}

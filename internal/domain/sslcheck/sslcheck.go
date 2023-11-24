package sslcheck

import "time"

type SSLCheck struct {
	TargetURL string    `json:"target_url"`
	LastCheck time.Time `json:"last_check"`
	Expiry    time.Time `json:"expiry"`
	DaysLeft  int       `json:"days_left"`
	Error     string    `json:"error"`
}

type SSLCheckRepository interface {
	ReadAll() ([]SSLCheck, error)
	WriteAll([]SSLCheck) error
	Read(hostname string) (SSLCheck, error)
	Write(val SSLCheck) error
	Delete(hostname string) error
}

func NewSSLCheck(targetURL string) SSLCheck {
	return SSLCheck{
		TargetURL: targetURL,
	}
}

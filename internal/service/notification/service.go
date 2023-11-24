package notification

import (
	"log"

	"github.com/ni3mm4nd/ssl-expiry-checker/internal/domain/sslcheck"
)

var ser *service

type INotification interface {
	Notify(check []sslcheck.SSLCheck)
}

type service struct {
	notifications []INotification
}

func Get() *service {
	if ser == nil {
		ser = &service{}
	}
	return ser
}

func NewServices(services ...INotification) *service {
	ser = &service{
		notifications: services,
	}
	return ser
}

func New(serv INotification) *service {
	return NewServices(serv)
}

func Add(serv INotification) {
	if ser == nil {
		ser = &service{}
	}
	ser.notifications = append(ser.notifications, serv)
}

func (s *service) NotifyAll(checks []sslcheck.SSLCheck) {
	log.Printf("Notifying all (%d) services: %#v", len(s.notifications), s.notifications)
	for _, n := range s.notifications {
		log.Printf("Notifying %T", n)
		go n.Notify(checks)
	}
}

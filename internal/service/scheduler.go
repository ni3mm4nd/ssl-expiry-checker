package service

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/domain/sslcheck"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/repository"
	"github.com/ni3mm4nd/ssl-expiry-checker/internal/service/notification"
)

func GetAllTargets() []string {
	targets, err := repository.GetRepos().SSLCheckRepo.ReadAll()
	if err != nil {
		panic(err)
	}

	res := make([]string, 0)

	for _, target := range targets {
		res = append(res, target.TargetURL)
	}

	return res
}

// NetworkConnector is an interface to abstract network connections.
type NetworkConnector interface {
	Dial(network, address string) (net.Conn, error)
}

// RealNetworkConnector implements NetworkConnector using the real net.Dialer.
type RealNetworkConnector struct{}

func (rnc RealNetworkConnector) Dial(network, address string) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	return tls.DialWithDialer(dialer, network, address, tlsConfig)
}

// FormatURL formats the given URL to the desired format.
// If the URL starts with "https://", this function removes the prefix,
// trims any trailing "/" character, and appends ":443" to indicate the
// default HTTPS port. If the URL does not start with schema, Path is used and appends":443".
// The formatted URL is returned as the result.
// for example https://example.com/ is going to return example.com:443
// for example example.com is going to return example.com:443
func formatURL(targetURL string) (string, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		log.Println(parsedURL)
		return "", errors.New("cannot parse URL")
	}
	if parsedURL.Scheme == "https" {
		targetURL = parsedURL.Host + ":443"
	}
	if parsedURL.Scheme == "" {
		targetURL = parsedURL.Path + ":443"
	}
	return targetURL, nil
}

// DialNetwork uses the provided NetworkConnector to establish a network connection.
func dialNetwork(targetURL string, connector NetworkConnector) (net.Conn, error) {
	return connector.Dial("tcp", targetURL)
}

// GetCertificateExpiryDate retrieves the expiry date of the peer certificate.
func getCertificateExpiryDate(conn net.Conn) (time.Time, error) {
	certChain := conn.(*tls.Conn).ConnectionState().PeerCertificates

	if len(certChain) == 0 {
		return time.Time{}, fmt.Errorf("no certificate found")
	}

	return certChain[0].NotAfter, nil
}

// CalculateRemainingDays calculates the remaining days until the given date.
func calculateRemainingDays(expiryDate time.Time) int {
	return int(expiryDate.Sub(time.Now()).Hours() / 24)
}

var schedulerInstance *gocron.Scheduler

func NewScheduler(cronExpression string, alertDaysLeft int) error {
	if cronExpression == "" {
		cronExpression = "@daily"
	}

	var oldSchedulerP *gocron.Scheduler

	if schedulerInstance != nil {
		log.Println("Scheduler already initialized")
		oldSchedulerP = schedulerInstance
	}

	schedulerInstance = gocron.NewScheduler(time.UTC)
	_, err := schedulerInstance.Cron(cronExpression).Do(func() {
		CheckTargets(GetAllTargets())
		allChecks, err := repository.GetRepos().SSLCheckRepo.ReadAll()
		if err != nil {
			return
		}
		expiringCeritficates := identifyExpiringCertificates(allChecks, alertDaysLeft)
		log.Printf("Found %d expiring certificates: %#v\n", len(expiringCeritficates), expiringCeritficates)
		if len(expiringCeritficates) > 0 {
			notification.Get().NotifyAll(expiringCeritficates)
		}
	})

	if err != nil {
		schedulerInstance = oldSchedulerP
		return err
	}

	if oldSchedulerP != nil {
		oldSchedulerP.Stop()
	}
	schedulerInstance.StartAsync()
	return nil
}

func identifyExpiringCertificates(sslChecks []sslcheck.SSLCheck, alertDaysLeft int) []sslcheck.SSLCheck {
	log.Printf("Checking for expiring certificates with alert days left: %d\n", alertDaysLeft)
	results := make([]sslcheck.SSLCheck, 0)
	for _, sslCheck := range sslChecks {
		if (sslCheck.DaysLeft <= alertDaysLeft) || (sslCheck.Error != "") {
			results = append(results, sslCheck)
		}
	}

	return results
}

func PrintScheduledJobs() {
	for i, job := range schedulerInstance.Jobs() {
		log.Printf("%d. Job %s next run: %s\n", i+1, job.GetName(), job.NextRun().String())
	}
}

func GetScheduler() *gocron.Scheduler {
	return schedulerInstance
}

func CheckTarget(target string) {
	CheckTargets([]string{target})
}

func CheckTargets(targets []string) {
	log.Printf("Checking SSL certificates: %v\n", targets)
	results := make([]sslcheck.SSLCheck, 0)
	var wg sync.WaitGroup

	// Check SSL certificates for each target
	for _, target := range targets {
		wg.Add(1)
		go func(target string) {
			defer wg.Done()
			sslCheck := sslcheck.SSLCheck{
				TargetURL: target,
				LastCheck: time.Now(),
				Error:     "",
			}
			formattedURL, err := formatURL(target)
			if err != nil {
				log.Println("Error:", err)
				sslCheck.Error = err.Error()
				results = append(results, sslCheck)
				return
			}
			conn, err := dialNetwork(formattedURL, RealNetworkConnector{})
			if err != nil {
				log.Println("Error:", err)
				sslCheck.Error = err.Error()
				results = append(results, sslCheck)
				return
			}
			defer conn.Close()
			expiryDate, err := getCertificateExpiryDate(conn)
			if err != nil {
				log.Println("Error:", err)
				sslCheck.Error = err.Error()
				results = append(results, sslCheck)
				return
			}
			remainingDays := calculateRemainingDays(expiryDate)
			sslCheck.Expiry = expiryDate
			sslCheck.DaysLeft = remainingDays
			results = append(results, sslCheck)
		}(target)
	}

	wg.Wait()

	for _, result := range results {
		err := repository.GetRepos().SSLCheckRepo.Write(result)
		if err != nil {
			log.Println(err)
		}
	}

}

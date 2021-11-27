package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	mffp "github.com/icemarkom/mffprober"
)

// // RebootFan executes HTTP POST to reboot the fan.
func RebootFan(cfg *mffp.Config) (int, error) {

	u := fmt.Sprintf("http://%s/%s", cfg.Host, mffp.MFFTarget)

	req, err := http.NewRequest("POST", u, bytes.NewBuffer([]byte(mffp.MFFReboot)))
	if err != nil {
		return 1, fmt.Errorf("error formatting request: %w", err)
	}

	hc := &http.Client{
		Timeout: cfg.Timeout,
	}
	r, err := hc.Do(req)
	if err != nil {
		// Fan does not respond to a reboot request, it just times-out.
		// We treat that as a potential success. Otherwise, it's an error.
		if err, ok := err.(*url.Error); ok && !err.Timeout() {
			return 1, fmt.Errorf("error rebooting the fan: %w", err)
		}
		r = new(http.Response)
		r.StatusCode = http.StatusRequestTimeout
	}
	// Just in a case the fans get API fix, treat StatusOK andd StatusRequestTimeout as success.
	if r.StatusCode != http.StatusOK && r.StatusCode != http.StatusRequestTimeout {
		return 1, fmt.Errorf("fan reported HTTP error: %s", r.Status)
	}
	if !cfg.Quiet {
		log.Printf("Reboot request to %q sent.", cfg.Host)
	}
	time.Sleep(cfg.Interval)
	fd, err := PollFan(cfg)
	if err != nil {
		return 1, fmt.Errorf("probably failed to reboot the fan: %w", err)
	}
	if !cfg.Quiet {
		log.Printf("Fan is alive: %s", fd)
	}
	return 0, nil
}

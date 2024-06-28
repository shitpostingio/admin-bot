package structs

import (
	"fmt"
	"strconv"
)

// WebhookConfiguration represents the webhook configuration
type WebhookConfiguration struct {
	Domain           string `type:"webhook"`
	IP               string `type:"webhook"`
	Port             int    `type:"webhook"`
	ReverseProxy     bool   `type:"webhook"`
	ReverseProxyPort int    `type:"webhook"`
	TLS              bool   `type:"webhook"`
	TLSCertPath      string `type:"webhook"`
	TLSKeyPath       string `type:"webhook"`
	MaxConnections   int    `type:"webhook"`
}

// BindString returns IP+Port, in a suitable syntax for http.ListenAndServe
func (c *WebhookConfiguration) BindString() string {
	return c.IP + ":" + strconv.Itoa(c.Port)
}

// WebHookURL returns the URL to listen on for WebHooks
func (c *WebhookConfiguration) WebHookURL(botToken string) string {

	port := c.Port
	if c.ReverseProxy {
		port = c.ReverseProxyPort
	}

	return fmt.Sprintf("https://%s:%d/%s/updates", c.Domain, port, botToken)
}

// IsStandardPort returns if the specified port is suitable
// for a webhook connection without reverse proxy
func IsStandardPort(port int) bool {
	switch port {
	case 443, 80, 88, 8443:
		return true
	default:
		return false
	}
}

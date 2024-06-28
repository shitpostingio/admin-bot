package utility

import (
	"fmt"
	"net/url"
	"strings"
)

// GetHostNameFromURL returns the host name given a url
func GetHostNameFromURL(inputHostname string) (string, error) {

	if !strings.HasPrefix(inputHostname, "http") {
		inputHostname = fmt.Sprintf("http://%s", inputHostname)
	}

	fullTextURL, err := UnshortenURL(inputHostname)
	if err != nil {
		return "", err
	}

	parsedURL, err := url.Parse(fullTextURL)
	if err != nil {
		return "", err
	}

	return strings.ToLower(parsedURL.Host), nil
}

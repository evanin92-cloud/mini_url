package utils

import (
	"net/url"
	"regexp"
	"strings"
)

var urlRegex = regexp.MustCompile(`^(https?://)[^\s/$.?#].[^\s]*$`)

func ValidateURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}

	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	if !urlRegex.MatchString(rawURL) {
		return false
	}

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	return parsed.Host != ""
}

func NormalizeURL(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}
	return rawURL
}
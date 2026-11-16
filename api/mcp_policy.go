package api

import (
	"net/url"
	"os"
	"strings"
)

func allowedMCPOrigin(origin, host string) bool {
	u, err := url.Parse(origin)
	if err != nil || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "" {
		return false
	}
	if (u.Scheme == "https" || u.Scheme == "http") && u.Host == host {
		return true
	}
	if origin == "https://chatgpt.com" || origin == "https://reminder.dev" {
		return true
	}
	for _, allowed := range strings.Split(os.Getenv("MCP_ALLOWED_ORIGINS"), ",") {
		if origin == strings.TrimSpace(allowed) {
			return true
		}
	}
	return false
}

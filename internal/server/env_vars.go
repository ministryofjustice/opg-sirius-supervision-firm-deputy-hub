package server

import (
	"os"
)

type EnvironmentVars struct {
	Port             string
	WebDir           string
	SiriusURL        string
	SiriusPublicURL  string
	ProHubURL        string
	Prefix           string
	FinanceAdminLink string
}

func NewEnvironmentVars() (EnvironmentVars, error) {
	return EnvironmentVars{
		Port:             getEnv("PORT", "1234"),
		WebDir:           getEnv("WEB_DIR", "web"),
		SiriusURL:        getEnv("SIRIUS_URL", "http://localhost:8080"),
		SiriusPublicURL:  getEnv("SIRIUS_PUBLIC_URL", "http://localhost:8080"),
		ProHubURL:        getEnv("PRO_HUB_HOST", "") + "/supervision/deputies",
		Prefix:           getEnv("PREFIX", ""),
		FinanceAdminLink: getEnv("FINANCE_ADMIN_LINK", "1"),
	}, nil
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}

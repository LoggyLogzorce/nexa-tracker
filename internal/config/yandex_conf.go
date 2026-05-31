package config

import (
	"golang.org/x/oauth2"
	"os"
)

func NewYandexConfig(cfg YandexOAuthConfig) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("YANDEX_CLIENT_ID"),
		ClientSecret: os.Getenv("YANDEX_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("YANDEX_REDIRECT_URL"),
		Scopes:       []string{"login:email", "login:info", "login:avatar"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.yandex.ru/authorize",
			TokenURL: "https://oauth.yandex.ru/token",
		},
	}
}

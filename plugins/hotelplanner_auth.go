package hotelplanner_auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func GetHotelPlannerAuthToken(apiKey string, secretKey string, accountId string, unixEpoch string) string {
	encodedAPIKey := base64.URLEncoding.EncodeToString([]byte(apiKey))
	signatureKey := fmt.Sprintf("%s|%s|%s", encodedAPIKey, accountId, unixEpoch)
	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write([]byte(signatureKey))
	hashValue := hash.Sum(nil)
	encodedHashValue := base64.URLEncoding.EncodeToString(hashValue)
	authToken := fmt.Sprintf("%s.%s", encodedAPIKey, encodedHashValue)
	return authToken
}

// Config the plugin configuration.
type Config struct {
	Headers  map[string]string
	HpConfig map[string]string
}

// CreateConfig creates plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers:  make(map[string]string),
		HpConfig: make(map[string]string),
	}
}

// Plugin.
type Plugin struct {
	next     http.Handler
	headers  map[string]string
	hpconfig map[string]string
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Headers) == 0 {
		return nil, fmt.Errorf("Headers cannot be empty!")
	}

	return &Plugin{
		headers:  config.Headers,
		next:     next,
		hpconfig: config.HpConfig,
	}, nil
}

func (a *Plugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Add headers passed from config
	for key, value := range a.headers {
		if value == "" {
			req.Header.Del(key)
		} else {
			req.Header.Set(key, value)
		}
	}
	// Getting config from dynamic config file
	apiKey := a.hpconfig["apiKey"]
	secretKey := a.hpconfig["secretKey"]
	accountId := a.hpconfig["accountId"]

	epoch := req.URL.Query().Get("epoch")
	// Add authorization header
	token := GetHotelPlannerAuthToken(apiKey, secretKey, accountId, epoch)
	req.Header.Set("Authorization", strings.TrimSuffix(token, "="))

	a.next.ServeHTTP(rw, req)
}

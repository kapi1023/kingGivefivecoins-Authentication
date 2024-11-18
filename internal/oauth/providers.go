package oauth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/linkedin"
)

type OAuthProvider string

const (
	Google   OAuthProvider = "google"
	LinkedIn OAuthProvider = "linkedin"
	Apple    OAuthProvider = "apple"
)

func GetProviderConfig(provider OAuthProvider, clientID, clientSecret, redirectURL string) (*oauth2.Config, error) {
	switch provider {
	case Google:
		return &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     google.Endpoint,
		}, nil
	case LinkedIn:
		return &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"r_liteprofile", "r_emailaddress"},
			Endpoint:     linkedin.Endpoint,
		}, nil
	// case Apple:
	// return &oauth2.Config{
	// ClientID:     clientID,
	// ClientSecret: clientSecret,
	// RedirectURL:  redirectURL,
	// Scopes:       []string{"name", "email"},
	// Endpoint:     apple.Endpoint,
	// }
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func ExchangeCodeForToken(ctx context.Context, provider OAuthProvider, code string, clientID, clientSecret, redirectURL string) (*oauth2.Token, error) {
	config, err := GetProviderConfig(provider, clientID, clientSecret, redirectURL)
	if err != nil {
		return nil, err
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, errors.New("failed to exchange authorization code for token: " + err.Error())
	}
	return token, nil
}

func ExtractEmailFromToken(token *oauth2.Token, provider OAuthProvider) (string, error) {
	switch provider {
	case Google:
		idToken, ok := token.Extra("id_token").(string)
		if !ok {
			return "", errors.New("id_token not found in token response")
		}

		parts := strings.Split(idToken, ".")
		if len(parts) != 3 {
			return "", errors.New("invalid id_token format")
		}

		payload, err := decodeBase64(parts[1])
		if err != nil {
			return "", err
		}

		var claims struct {
			Email string `json:"email"`
		}
		if err := json.Unmarshal(payload, &claims); err != nil {
			return "", errors.New("failed to parse id_token payload")
		}

		return claims.Email, nil
	case LinkedIn:
		return fetchLinkedInEmail(token)
	// case Apple:
	// return fetchAppleEmail(token)
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
}

func decodeBase64(encoded string) ([]byte, error) {
	for len(encoded)%4 != 0 {
		encoded += "="
	}
	return base64.URLEncoding.DecodeString(encoded)
}

func fetchLinkedInEmail(token *oauth2.Token) (string, error) {
	req, err := http.NewRequest("GET", "https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch email from LinkedIn API")
	}

	var result struct {
		Elements []struct {
			Handle struct {
				EmailAddress string `json:"emailAddress"`
			} `json:"handle~"`
		} `json:"elements"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Elements) == 0 || result.Elements[0].Handle.EmailAddress == "" {
		return "", errors.New("email not found in LinkedIn API response")
	}

	return result.Elements[0].Handle.EmailAddress, nil
}

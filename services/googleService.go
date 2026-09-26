package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleService struct {
	config *oauth2.Config
}

type GoogleUser struct {
	ID      string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func NewGoogleService(
	clientID string,
	clientSecret string,
	redirectURL string,
) *GoogleService {
	return &GoogleService{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"openid",
				"email",
				"profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (s *GoogleService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
	)
}

func (s *GoogleService) GetUser(
	ctx context.Context,
	code string,
) (*GoogleUser, error) {

	token, err := s.config.Exchange(ctx, code)

	if err != nil {
		return nil, fmt.Errorf("failed to exchange google code: %w", err)
	}

	client := s.config.Client(ctx, token)

	response, err := client.Get(
		"https://www.googleapis.com/oauth2/v3/userinfo",
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get google user: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"google userinfo returned status %d",
			response.StatusCode,
		)
	}

	var user GoogleUser

	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf(
			"failed to decode google user: %w",
			err,
		)
	}

	return &user, nil
}

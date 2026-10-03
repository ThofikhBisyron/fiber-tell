package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type FacebookService struct {
	config *oauth2.Config
}

type FacebookUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

func NewFacebookService(
	clientID string,
	clientSecret string,
	redirectURL string,
) *FacebookService {
	return &FacebookService{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"email",
				"public_profile",
			},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.facebook.com/v24.0/dialog/oauth",
				TokenURL: "https://graph.facebook.com/v24.0/oauth/access_token",
			},
		},
	}
}

func (s *FacebookService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

func (s *FacebookService) GetUser(
	ctx context.Context,
	code string,
) (*FacebookUser, error) {
	token, err := s.config.Exchange(ctx, code)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to exchange facebook code: %w",
			err,
		)
	}

	params := url.Values{}
	params.Set("fields", "id,name,email,picture")
	params.Set("access_token", token.AccessToken)

	response, err := http.Get(
		"https://graph.facebook.com/v24.0/me?" + params.Encode(),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get facebook user: %w",
			err,
		)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"facebook userinfo returned status %d",
			response.StatusCode,
		)
	}

	var user FacebookUser

	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode facebook user: %w",
			err,
		)
	}

	return &user, nil
}

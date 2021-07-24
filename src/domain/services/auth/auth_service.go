package authservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"oauth2_example/src/config"
	inmemorycache "oauth2_example/src/infra/in_memory_cache"
)

const oAuth2RedirectURI = "http://localhost:3000/auth/oauth2_callback"

var cache = inmemorycache.New()

func GenerateOAuth2URI() string {
	requestIdentifier := uuid.NewString()

	cache.SetWithExpiration(requestIdentifier, struct{}{}, 1*time.Minute)

	return fmt.Sprintf(
		"https://accounts.spotify.com/authorize?response_type=code&client_id=%s&scope=%s&redirect_uri=%s&state=%s",
		config.Configs.OAuth2.ClientID,
		url.QueryEscape("user-read-private user-read-email"),
		url.QueryEscape(oAuth2RedirectURI),
		url.QueryEscape(requestIdentifier),
	)
}

type requestSpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func exchangeCodeForToken(code string) (*requestSpotifyTokenResponse, error) {
	form := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {oAuth2RedirectURI},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	encodedForm := form.Encode()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://accounts.spotify.com/api/token",
		strings.NewReader(encodedForm),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", strconv.Itoa(len(encodedForm)))

	request.SetBasicAuth(
		url.QueryEscape(config.Configs.OAuth2.ClientID),
		url.QueryEscape(config.Configs.OAuth2.ClientSecret),
	)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer response.Body.Close()

	responseBodyAsBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var tokenResponse requestSpotifyTokenResponse

	err = json.Unmarshal(responseBodyAsBytes, &tokenResponse)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &tokenResponse, nil
}

var ErrInvalidRequestIdentifier = errors.New("Invalid request identifier")

func HandleOAuth2Callback(code, requestIdentifier string) (*UserInfo, error) {
	if !cache.Has(requestIdentifier) {
		return nil, errors.WithStack(ErrInvalidRequestIdentifier)
	}
	cache.Delete(requestIdentifier)

	tokenResponse, err := exchangeCodeForToken(code)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userInfo, err := fetchUserInfo(tokenResponse.AccessToken)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return userInfo, nil
}

type UserInfo struct {
	ID          string `json:"id"`
	Country     string `json:"country"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func fetchUserInfo(token string) (*UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.spotify.com/v1/me",
		nil,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer response.Body.Close()

	responseBodyAsBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var info UserInfo
	err = json.Unmarshal(responseBodyAsBytes, &info)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &info, nil
}

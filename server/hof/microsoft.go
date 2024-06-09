package hof

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type Header struct {
	Kid string `json:"kid"`
}

var (
	httpClient    = &http.Client{}
	metadataCache map[string]interface{}
	jwksCache     map[string]interface{}
	cacheMutex    sync.Mutex
)

func fetchJSON(
	ctx context.Context,
	url string,
) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func getCachedJSON(
	ctx context.Context,
	url string,
	cache *map[string]interface{},
) (map[string]interface{}, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if *cache == nil {
		data, err := fetchJSON(ctx, url)
		if err != nil {
			return nil, err
		}
		*cache = data
	}
	return *cache, nil
}

func getJWKSKey(
	jwks map[string]interface{},
	kid string,
) (map[string]interface{}, error) {
	for _, data := range jwks["keys"].([]interface{}) {
		if k, ok := data.(map[string]interface{}); ok {
			if k["kid"] == kid {
				return k, nil
			}
		}
	}
	return nil, errors.New("key not found")
}

func parseRSAKey(key map[string]interface{}) (*rsa.PublicKey, error) {
	nStr := key["n"].(string)
	eStr := key["e"].(string)
	nStr = strings.ReplaceAll(nStr, "_", "/")
	nStr = strings.ReplaceAll(nStr, "-", "+")
	nStr += "=="
	decodedN, err := base64.StdEncoding.DecodeString(nStr)
	if err != nil {
		return nil, err
	}
	decodedE, err := base64.StdEncoding.DecodeString(eStr)
	if err != nil {
		return nil, err
	}
	publicKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(decodedN),
		E: int(new(big.Int).SetBytes(decodedE).Int64()),
	}
	return publicKey, nil
}

func GetMicrosoftUserProfileFromToken(
	ctx context.Context,
	token *oauth2.Token,
	tenantID string,
) (map[string]string, error) {
	url := "https://login.microsoftonline.com/%s/v2.0/.well-known/openid-configuration"
	confURL := fmt.Sprintf(url, tenantID)
	metadata, err := getCachedJSON(ctx, confURL, &metadataCache)
	if err != nil {
		return nil, err
	}
	jwksURL := metadata["jwks_uri"].(string)
	jwks, err := getCachedJSON(ctx, jwksURL, &jwksCache)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(token.AccessToken, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}
	headerPart := parts[0]
	headerJson, err := base64.RawURLEncoding.DecodeString(headerPart)
	if err != nil {
		return nil, err
	}
	var header Header
	if err := json.Unmarshal(headerJson, &header); err != nil {
		return nil, err
	}
	key, err := getJWKSKey(jwks, header.Kid)
	if err != nil {
		return nil, err
	}
	publicKey, err := parseRSAKey(key)
	if err != nil {
		return nil, err
	}
	tkn, _ := jwt.Parse(token.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	claims := tkn.Claims.(jwt.MapClaims)
	return map[string]string{
		"name":  claims["name"].(string),
		"email": claims["email"].(string),
	}, nil
}

func GetMicrosoftOAuthConfig() (token *oauth2.Config, TenantID string) {
	b, err := os.ReadFile("microsoft.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	type microsoftCredentials struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		TenantID     string `json:"tenant_id"`
		RedirectURI  string `json:"redirect_uri"`
	}
	var j struct {
		Web *microsoftCredentials `json:"web"`
	}
	if err := json.Unmarshal(b, &j); err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return &oauth2.Config{
		ClientID:     j.Web.ClientID,
		ClientSecret: j.Web.ClientSecret,
		RedirectURL:  j.Web.RedirectURI,
		Endpoint:     microsoft.AzureADEndpoint(j.Web.TenantID),
		Scopes: []string{
			"User.Read", "email",
			"CallEvents.Read",
			"Calendars.ReadWrite",
		},
	}, j.Web.TenantID
}

func GetMicrosoftOAuthTokenFromWeb(config *oauth2.Config) (string, *oauth2.Token) {
	authURL := config.AuthCodeURL(
		"state", oauth2.AccessTypeOffline)
	if authURL != "" {
		return authURL, nil
	}
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return "", tok
}

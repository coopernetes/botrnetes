package twitch

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

const (
	idEnvVar       = "TWITCH_CLIENT_ID"
	secretEnvVar   = "TWITCH_CLIENT_SECRET"
	redirectEnvVar = "TWITCH_REDIRECT_URL"
)

var (
	oauth2Config *oauth2.Config
	httpClient   *http.Client
	ctx          *context.Context
)

func Init() string {
	log.Printf("Initializing")
	ctx := context.Background()
	id, secret := lookup(idEnvVar), lookup(secretEnvVar)
	chatScopes := []string{"chat:read", "chat:edit"}

	redirectUrl := os.Getenv(redirectEnvVar)
	if redirectUrl == "" {
		redirectUrl = "http://localhost:8080/login"
	}

	oauth2Config = &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint:     twitch.Endpoint,
		RedirectURL:  redirectUrl,
		Scopes:       chatScopes,
	}

	randomState := make([]string, 32)
	alphanum := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	for i, _ := range randomState {
		ranChar := alphanum[rand.Intn(len(alphanum))]
		randomState[i] = string(ranChar)
	}
	url := oauth2Config.AuthCodeURL(strings.Join(randomState, ""))
	log.Printf("Authorize this app here: %s", url)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := oauth2Config.Client(ctx, tok)
	log.Printf("Sending test request on new httpClient")
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/users?login=twitchdev", nil)
	req.Header.Add("Client-Id", id)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	var body strings.Builder
	for scanner.Scan() {
		body.WriteString(scanner.Text())
	}
	log.Printf("Response (%d): %s", resp.StatusCode, body.String())

	return tok.AccessToken // return token once auth'd
}

func lookup(envVar string) string {
	o, set := os.LookupEnv(envVar)
	if !set {
		log.Fatalf("%s is unset!", envVar)
	}
	if o == "" {
		log.Fatalf("%s is empty!", envVar)
	}
	return o
}

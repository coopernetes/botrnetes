package twitch

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

const (
	idEnvVar       = "TWITCH_CLIENT_ID"
	secretEnvVar   = "TWITCH_CLIENT_SECRET"
	redirectEnvVar = "TWITCH_REDIRECT_URL"
)

var (
	state        string
	oauth2Config *oauth2.Config
	httpClient   *http.Client
	ctx          *context.Context
)

func Init() string {
	log.Printf("Initializing")
	ctx := context.Background()
	id, secret := lookup(idEnvVar), lookup(secretEnvVar)
	// initialize OAuth2 state var for this running instance
	randomState := make([]string, 32)
	alphanum := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	for i, _ := range randomState {
		ranChar := alphanum[rand.Intn(len(alphanum))]
		randomState[i] = string(ranChar)
	}
	state := strings.Join(randomState, "")

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

	url := oauth2Config.AuthCodeURL(state)
	log.Printf("Authorize this app here: %s", url)

	c := make(chan *oauth2.Token, 1)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("state") != state {
			log.Print("Received non-matching state")
			http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
			return
		}
		if r.Form.Has("error") {
			log.Printf("Received error during login, details: %s", fmt.Sprintf("error=%s, description=%s", r.FormValue("error"), r.FormValue("error_description")))
			http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
			return
		}
		defer close(c)
		received := r.FormValue("code")
		log.Printf("Received code %s", received)
		tok, err := oauth2Config.Exchange(ctx, received)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Got token, expires in %s", tok.Expiry.String())
		c <- tok
	})
	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Login failed! Check the logs.")
	})
	log.Printf("Starting http server")
	go http.ListenAndServe(":8080", nil)

	log.Printf("Waiting for token exchange")
	select {
	case token := <-c:
		httpClient = oauth2Config.Client(ctx, token)
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

		return token.AccessToken // return token once auth'd
	case <-time.After(time.Minute * 2):
		log.Fatal("2m timeout reached, exiting...")
	}
	return ""
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

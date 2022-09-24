package twitch

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

const (
	idEnvVar     = "TWITCH_CLIENT_ID"
	secretEnvVar = "TWITCH_CLIENT_SECRET"
)

var (
	oauth2Config *clientcredentials.Config
	httpClient   *http.Client
)

func Init() {
	log.Printf("Initializing")
	ctx := context.Background()
	id, secret := lookup(idEnvVar), lookup(secretEnvVar)
	oauth2Config = &clientcredentials.Config{
		ClientID:     id,
		ClientSecret: secret,
		TokenURL:     twitch.Endpoint.TokenURL,
	}

	tSource := oauth2Config.TokenSource(ctx)
	//if err != nil {
	//		log.Fatal(err)
	//}

	httpClient := oauth2.NewClient(ctx, tSource)
	log.Printf("Done, httpClient=%p", &httpClient)

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users?login=twitchdev", nil)
	req.Header.Add("Client-Id", id)
	if err != nil {
		log.Fatal(err)
	}

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
	log.Printf("%s", body.String())
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

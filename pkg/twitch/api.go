package twitch

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

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

	httpC := oauth2.NewClient(ctx, tSource)
	log.Printf("Done, httpC=%p", &httpC)

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users?login=twitchdev", nil)
	req.Header.Add("Client-Id", id)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := httpC.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body := make([]byte, 0, 0)
	_, berr := resp.Body.Read(body)
	if berr != nil {
		log.Fatal(berr)
	}
	j, err := json.Marshal(body)
	log.Printf("%s", j)
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

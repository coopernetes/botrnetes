package main

import (
	"log"
	"os"

	t "github.com/coopernetes/botrnetes/pkg/twitch"
)

const (
	twitchId     = "TWITCH_CLIENT_ID"
	twitchSecret = "TWITCH_CLIENT_SECRET"
)

func main() {
	id := lookup(twitchId)
	secret := lookup(twitchSecret)
	t.Init(id, secret)
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

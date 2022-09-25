package twitch

import (
	"fmt"
	"log"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"
)

var (
	client *twitch.Client
)

func StartChat(token string) {
	client := twitch.NewClient("botrnetes", fmt.Sprintf("oauth:%s", token))
	client.Join("roastedfunction")
	log.Printf("Connecting to %s", client.IrcAddress)
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	client.Say("roastedfunction", "I'm here!")
}

func isCommand(msg twitch.PrivateMessage) bool {
	return strings.HasPrefix(strings.ToLower(msg.Message), "!")
}

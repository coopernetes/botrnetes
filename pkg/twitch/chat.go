package twitch

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
)

const channelEnvVar = "TWITCH_CHANNEL"

func StartChat(token string) {
	channel := os.Getenv(channelEnvVar)
	if channel == "" {
		channel = "roastedfunction"
	}

	client := twitch.NewClient("botrnetes", fmt.Sprintf("oauth:%s", token))
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Printf(message.Message)
	})
	client.OnSelfJoinMessage(func(message twitch.UserJoinMessage) {
		client.Say(channel, "I joined!")
	})
	client.Join(channel)
	go client.Connect()
	log.Printf("Before say")
	time.Sleep(time.Second * 60)
	client.Say("roastedfunction", "Hello! :wave:")
}

func isCommand(msg twitch.PrivateMessage) bool {
	return strings.HasPrefix(strings.ToLower(msg.Message), "!")
}

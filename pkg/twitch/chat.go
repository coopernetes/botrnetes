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

var (
	client  twitch.Client
	channel string
)

func StartChat(token string) {
	channel = os.Getenv(channelEnvVar)
	if channel == "" {
		channel = "roastedfunction"
	}

	client := twitch.NewClient("botrnetes", fmt.Sprintf("oauth:%s", token))
	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		log.Print(msg.Message)
		if !strings.HasPrefix(msg.Message, "!") {
			return
		}
		log.Printf("Received command %s", msg.Message)
		s := strings.TrimSpace(strings.ToLower(msg.Message))
		switch s {
		case "!github":
			client.Say(channel, "https://github.com/coopernetes for random bits & blops")
		default:
			log.Printf("No command for %s, ignoring.", msg.Message)
		}
	})
	client.OnSelfJoinMessage(func(message twitch.UserJoinMessage) {
		client.Say(channel, fmt.Sprintf("I joined at %s!", time.Now()))
	})
	client.Join(channel)
	go client.Connect()
	// Loop infinitely, sending a whisper every 10min to check in
	for {
		time.Sleep(time.Second * 60)
		client.Whisper(channel, "Just checking in!")
	}
}

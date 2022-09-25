package main

import (
	t "github.com/coopernetes/botrnetes/pkg/twitch"
)

func main() {
	token := t.Init()
	t.StartChat(token)
}

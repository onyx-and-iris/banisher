package main

import (
	"fmt"
	"log"

	"github.com/gtuk/discordwebhook"
)

type notifier interface {
	send(string, string)
}

type identifiers struct {
	Name string
	Url  string
}

type discordNotifier struct {
	identifiers
}

func newDiscordNotifier(url string) discordNotifier {
	dn := discordNotifier{}
	dn.Name = "discord"
	dn.Url = url
	return dn
}

func (dn discordNotifier) send(ruleName, ip string) {
	username := "Banisher"
	content := fmt.Sprintf("%s violation for %s", ruleName, ip)

	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}

	err := discordwebhook.SendMessage(dn.Url, message)
	if err != nil {
		log.Printf(err.Error())
	}
}

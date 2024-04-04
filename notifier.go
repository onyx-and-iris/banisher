package main

import (
	"fmt"
	"log"

	"github.com/gtuk/discordwebhook"
)

type notifier struct {
	Name string
	Url  string
}

func (n notifier) send(ruleName, ip string) {
	switch n.Name {
	case "discord":
		n.discordWebhook(ruleName, ip)
	}
}

func (n notifier) discordWebhook(ruleName, ip string) {
	var username = "Banisher"
	var content = fmt.Sprintf("%s violation for %s", ruleName, ip)
	var url = n.Url

	message := discordwebhook.Message{
		Username: &username,
		Content:  &content,
	}

	err := discordwebhook.SendMessage(url, message)
	if err != nil {
		log.Printf(err.Error())
	}
}

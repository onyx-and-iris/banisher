package main

import (
	"fmt"
	"log"
	"os/user"
)

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

func getRules(ruleName string) (rule, error) {
	for _, r := range config.Rules {
		if r.Name == ruleName {
			return r, nil
		}
	}
	return rule{}, fmt.Errorf("no rule set named %s", ruleName)
}

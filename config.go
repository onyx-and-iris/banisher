package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DefaultBanishmentDuration uint
	Whitelist                 []string
	Rules                     []rule
	Notifiers                 []notifier
}

func loadConfig(path string) (conf Config, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return conf, err
	}

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		return
	}
	//fmt.Printf("--- m:\n%v\n\n", m)

	// DefaultBanishmentDuration
	conf.DefaultBanishmentDuration = uint(m["defaultBanishmentDuration"].(int))

	// whitelist
	if m["whitelist"] != nil {
		for _, ip := range m["whitelist"].([]interface{}) {
			if net.ParseIP(ip.(string)) == nil {
				return conf, fmt.Errorf("%s is not a valid IP for whitelist", ip.(string))
			}
			if conf.isIPWhitelisted(ip.(string)) {
				return conf, fmt.Errorf("%s appears multiple time in your whitelist", ip.(string))

			}
			conf.Whitelist = append(conf.Whitelist, ip.(string))
		}
	}

	// rules
	for _, r := range m["rules"].([]interface{}) {
		rule2add := rule{}

		rs := r.(map[string]interface{})

		// name
		if rs["name"] == nil {
			return conf, errors.New("required field 'name' is missing in a rule")
		}
		rule2add.Name = rs["name"].(string)

		// ippos
		if rs["IPpos"] != nil {
			rule2add.IPpos = uint(rs["IPpos"].(int))
		} else {
			rule2add.IPpos = 0
		}

		// match
		if rs["match"] == nil {
			return conf, errors.New("required field 'match' is missing in a rule")
		}
		// to regex
		rule2add.Match, err = regexp.Compile(rs["match"].(string))
		if err != nil {
			return conf, fmt.Errorf("failed to compile regex: %s, %v", rs["match"].(string), err)
		}

		// notify?
		if rs["notify"] != nil {
			rule2add.Notify = rs["notify"].(bool)
		}

		// append rule
		conf.Rules = append(conf.Rules, rule2add)
	}

	// notifiers
	if notifiers, ok := m["notifiers"]; ok {
		for _, n := range notifiers.([]interface{}) {
			notifier2add := notifier{}

			ns := n.(map[string]interface{})

			//discord
			if ns["discord"] != nil {
				notifier2add.Name = "discord"
				notifier2add.Url = ns["discord"].(string)
			}

			conf.Notifiers = append(conf.Notifiers, notifier2add)
		}
	}

	return
}

// check if ip is whitelisted
func (c Config) isIPWhitelisted(ip string) bool {
	for _, ipw := range c.Whitelist {
		if ip == ipw {
			return true
		}
	}
	return false
}

func (c Config) rulesByName(ruleName string) (rule, error) {
	for _, r := range c.Rules {
		if r.Name == ruleName {
			return r, nil
		}
	}
	return rule{}, fmt.Errorf("no rule set named %s", ruleName)
}

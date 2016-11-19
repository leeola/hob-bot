package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	hobbot "github.com/leeola/hob-bot"
	"github.com/leeola/hob-bot/acls/simple"
	"github.com/leeola/hob-bot/messengers/telegram"
)

func main() {
	// parse the config path option
	var configPath string
	flag.StringVar(&configPath, "config", "./config.toml", "path to hob-bot config")
	flag.Parse()

	// open the config and decode the toml
	f, err := os.Open(configPath)
	if os.IsNotExist(err) {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var conf Config
	if _, err := toml.DecodeReader(f, &conf); err != nil {
		panic(err)
	}

	var messenger hobbot.Messenger
	if !conf.Telegram.IsZero() {
		m, err := telegram.New(telegram.Config{
			ApiToken: conf.Telegram.ApiToken,
		})
		if err != nil {
			panic(err)
		}
		messenger = m
	}

	// Because we're using a SimpleAcl, combine our two Acl lists into the SimpleAcl.
	acl := simple.SimpleAcl(map[string][]string{})
	if conf.Acl.Users != nil {
		for user, perms := range conf.Acl.Users {
			acl[user] = perms
		}
	}

	if conf.Acl.Channels != nil {
		for channel, perms := range conf.Acl.Channels {
			if _, ok := acl[channel]; ok {
				panic(fmt.Sprintf("channel %q overwrites username", channel))
			}
			acl[channel] = perms
		}
	}

	b, err := hobbot.New(hobbot.Config{
		Acl:       acl,
		HobAddr:   conf.HobAddr,
		Messenger: messenger,
	})
	if err != nil {
		panic(err)
	}

	if err := b.ListenAndServe(); err != nil {
		panic(err)
	}
}

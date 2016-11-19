package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/inconshreveable/log15"
	"github.com/leeola/errors"
	"github.com/leeola/hob-bot"
)

type Config struct {
	// The Telegram Bot API Token, as given by "BotFather"
	ApiToken string
	Log      log15.Logger
}

type Messenger struct {
	botApi *tgbotapi.BotAPI
	log    log15.Logger
}

func New(c Config) (*Messenger, error) {
	if c.Log == nil {
		c.Log = log15.New()
	}

	b, err := tgbotapi.NewBotAPI(c.ApiToken)
	if err != nil {
		log.Panic(err)
	}

	return &Messenger{
		botApi: b,
		log:    c.Log,
	}, nil
}

func (m *Messenger) Messages() (<-chan hobbot.Message, error) {
	m.log.Debug("starting update listener")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := m.botApi.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}

	// sync chan for now, may make it buffered in the future.
	botChan := make(chan hobbot.Message)
	go m.sendMessages(updates, botChan)

	return botChan, nil
}

func (m *Messenger) sendMessages(from <-chan tgbotapi.Update, to chan<- hobbot.Message) {
	for u := range from {
		log := m.log.New("updateId", u.UpdateID)

		// ignore updates without a message, as hob-bot currently doesn't handle
		// any of that.
		if u.Message == nil {
			log.Debug("dropping non-message")
			continue
		}

		from := u.Message.From

		// ignore updates without a From (*User), as hob-bot only cares about messages
		// from users.
		if from == nil {
			log.Debug("dropping message with no sender")
			continue
		}

		// construct the username from the id + username.
		// Why? Telegram has optional usernames, so the id must be used as the unique
		// identifier for each user. If however they have a username, we can use both
		// for extra specificity.
		//
		// The ids are the hard/ugly part, the part that impacts user experience when
		// adding an ACL/etc, so we're not losing any UX by adding the username... it's
		// already bad.. Telegram doesn't give us a choice.
		var username string
		if from.UserName != "" {
			username = fmt.Sprintf("%d:%s", from.ID, from.UserName)
		} else {
			username = strconv.Itoa(from.ID)
		}

		chat := u.Message.Chat
		if chat == nil {
			log.Warn("missing chat object", "message", u.Message)
			continue
		}

		if chat.Type != "private" {
			log.Debug("received chat update", "chat", chat)

			// drop chat updates for now, i want to see exactly how Telegram behaves.
			log.Crit("dropping chat update, not implemented")
		}

		log.Debug("chat info..", "chat", chat)

		to <- hobbot.Message{
			Sender:  username,
			Message: u.Message.Text,
			Channel: strconv.Itoa(int(chat.ID)),
		}
	}
}

func (m *Messenger) Send(recipient, message string) error {
	// since we combine the id & username for telegram, split on the first colon
	// if it exists.
	//
	// Note that the ID form Telegram is an Int, so it will never contain a colon.
	// This is safe.
	split := strings.SplitN(recipient, ":", 2)
	parsedInt, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse id from recipient")
	}

	msg := tgbotapi.NewMessage(parsedInt, message)
	m.botApi.Send(msg)

	return nil
}

package hobbot

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/inconshreveable/log15"
	"github.com/leeola/errors"
)

type Config struct {
	Acl       Acl
	HobAddr   string
	Log       log15.Logger
	Messenger Messenger
}

type Bot struct {
	acl       Acl
	hobAddr   *url.URL
	messenger Messenger
	log       log15.Logger
}

func New(c Config) (*Bot, error) {
	if c.Log == nil {
		c.Log = log15.New()
	}

	if c.Messenger == nil {
		return nil, errors.New("missing required argument: Messenger")
	}

	hobAddr, err := url.Parse(c.HobAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}

	return &Bot{
		acl:       c.Acl,
		hobAddr:   hobAddr,
		log:       c.Log,
		messenger: c.Messenger,
	}, nil
}

func (b *Bot) Send(r, m string) error {
	return b.messenger.Send(r, m)
}

func (b *Bot) ListenAndServe() error {
	messages, err := b.messenger.Messages()
	if err != nil {
		return err
	}

	for m := range messages {
		// TODO(leeola): Make the event splitting code configurable. Likely just
		// allow passing in a eventsplit func into the config. Eg:
		// 		Config.CommandSplitFunc(string) (string,string)

		switch {
		// If it doesn't start with ! it's not a event
		case !strings.HasPrefix(m.Message, "!"):
			b.log.Debug("dropping message, no event prefix")
			continue
			// "!" is not a event, ignore it.
		case m.Message == "!":
			b.log.Debug("dropping message, no event after prefix")
			continue
		}

		// Note that we're trimming the first char
		split := strings.SplitN(m.Message[1:], " ", 2)

		event := split[0]
		var payload string
		if len(split) > 1 {
			payload = split[1]
		}

		if b.acl != nil && !b.acl.Allowed(m.Sender, m.Channel, event) {
			b.log.Info("acl rejected command",
				"from", m.Sender, "in", m.Channel, "command", event)
			continue
		}

		b.log.Debug("sending Hob event", "event", event, "payload", payload)

		u := *b.hobAddr
		// join incase the given path is behind a embedded router
		u.Path = path.Join(u.Path, "events", event)
		uStr := u.String()

		b.log.Debug("sending Hob event", "event", event, "url", uStr, "payload", payload)

		res, err := http.Post(uStr, "text/plain", strings.NewReader(m.Message))
		if err != nil {
			b.log.Error("failed to post to Hob", "err", err)
			b.messenger.Send(m.Channel, "error encountered from Hob, see log for details")
			continue
		}
		defer res.Body.Close()

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(res.Body)
		if err != nil {
			b.log.Error("failed to read from response body", "err", err)
			b.messenger.Send(m.Channel, "error encountered in Hob-bot, see log for details")
			continue
		}
		body := buf.String()

		if res.StatusCode != http.StatusOK {
			b.log.Warn("Hob returned non-okay status", "code", res.StatusCode, "status", res.Status, "body", body)

			if body != "" {
				b.messenger.Send(m.Channel, fmt.Sprintf(
					"Hob event returned non-okay status. code=%d, status=%s",
					res.StatusCode, res.Status,
				))
			} else {
				b.messenger.Send(m.Channel, "Hob event returned non-okay status, see log for details")
			}
		}

		if body != "" {
			b.messenger.Send(m.Channel, body)
		}
	}

	return nil
}

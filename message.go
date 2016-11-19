package hobbot

type Message struct {
	Sender  string
	Message string
	Channel string

	// TODO(leeola) Support arbitrary data, such as image attachments that
	// are to be sent to Hob.
	// Data io.ReadCloser
}

type Messenger interface {
	Messages() (<-chan Message, error)

	// TODO(leeola) Support arbitrary data, such as image attachments that
	// are to be sent to Hob.
	Send(recipient, message string) error
}

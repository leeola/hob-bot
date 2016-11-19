package webhooks

type Message struct {
	Message   string `json:"message"`
	Recipient string `json:"recipient"`
}

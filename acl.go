package hobbot

type Acl interface {
	Allowed(sender, channel, command string) bool
}

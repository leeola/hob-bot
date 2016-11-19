package simple

// SimpleAcl implements a very simple Acl for hob-bot.
//
// It works by first checking if the username is in the approved list,
// and then checking if the command is allowed for the given user.
//
// The channel is checked in the same way. If the channel exists,
// the command is checked if it is allowed to be used in the
// given channel.
//
// A map entry *must* exist for both the user calling the command
// and the channel the command is being run in.
//
// In both user and channel if the []string is empty the caller has access to all
// commands. Ie, if user John has no commands specified he can run any command,
// assuming that command is also allowed in the given channel.
//
// Note that a channel and username can ofcourse conflict in this model. If this
// is a problem, use another Acl.
type SimpleAcl map[string][]string

func (a SimpleAcl) Allowed(username, channel, command string) bool {
	// first check to ensure the user can run the given command
	allowedUserCommands, ok := a[username]
	if !ok {
		return false
	}

	var commandIsAllowed bool
	for _, cmd := range allowedUserCommands {
		if cmd == command {
			commandIsAllowed = true
		}
	}

	if !commandIsAllowed {
		return false
	}

	// then check to ensure the channel can see the given command (useful for
	// private/sensitive commands)
	allowedChanCommands, ok := a[channel]
	if !ok {
		return false
	}

	commandIsAllowed = false
	for _, cmd := range allowedChanCommands {
		if cmd == command {
			commandIsAllowed = true
		}
	}

	if !commandIsAllowed {
		return false
	}

	return true
}

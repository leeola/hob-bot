
# hob-bot (unstable)

A bot plugin for Hob which is agnostic to the backend chat platform.

hob-bot serves to be the translation layer *(with a bit of sugar, like ACL)*
between Hob and the chat backend in use. ie Slack, Telegram, etc. Like Hob,
hob-bot aims to be simple and low of features.

## Usage

hob-bot works by simply translating and then proxying chat commands into
Hob events. It has a basic webhook API for communicating action results to
slack

## License

MIT

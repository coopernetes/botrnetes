# botrnetes, my Twitch/YT bot
A simple chat bot to post content based on commands.

The following environment variables must be set:

* `TWITCH_CLIENT_ID`
* `TWITCH_CLIENT_SECRET`
* `TWITCH_REDIRECT_URL`

## TODO
- [x] irc join
- [x] irc listen
- [x] receive auth code via HTTP server
- [ ] irc respond to a command
- [ ] deploy the server

## Install
```shell
go install github.com/coopernetes/botrnetes/cmd/bot@latest
```

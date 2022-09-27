# botrnetes, my Twitch/YT bot
A simple chat bot to post content based on commands.

The following environment variables must be set:

* `TWITCH_CLIENT_ID`
* `TWITCH_CLIENT_SECRET`
* `TWITCH_REDIRECT_URL`

## Build
```shell
go build cmd/bot/main.go
```

Build the image.
```shell
docker build -t coopernetes/botrnetes .
```

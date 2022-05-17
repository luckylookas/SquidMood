# Squidmood Discord bot


## v1.0

default prefix: `!squid `

### building squidmood

`go build -o squidmood`

tested with go 1.17 or later

### running squidmood

create a file named `config.json` at `[squidmood-binary-location]/private` containing:
```json
{
  "app_id": "[YOUR APP ID]",
  "bot_token": "[YOUR BOT TOKEN]"
}

```

then just run the binary

### commands

* `!squid ask` -  posts the available Squidwards
* reacting to `!squid ask` with the `:one` to `:nine:` emoji - selects and stores your mood
* `!squid [1-9]` - selects and stores your mood (no need for `!squid ask` if you know the number by heart)
* `!squid @Mention` - tells you `@Mention`'s mood

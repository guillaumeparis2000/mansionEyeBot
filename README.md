Before executing the bot remember to export yout telegram token like this:

```sh
export TELEGRAM_TOKEN=your:private_token
```

## Compilation for Raspberry Pi 1

```sh
GOOS=linux GOARCH=arm GOARM=6 go build  -ldflags="-s -w" -v main.go
```

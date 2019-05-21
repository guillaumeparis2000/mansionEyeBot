Before executing the bot remember to export yout telegram token like this:

```sh
export TELEGRAM_TOKEN=your:private_token
```

## Compilation for Raspberry Pi 1
This command will build for `ARM-6`. The `ldflags` `-s` and `-w` strip the debugging information and reduce the size of the executable.

```sh
GOOS=linux GOARCH=arm GOARM=6 go build  -ldflags="-s -w" -v main.go
```

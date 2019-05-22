Before executing the bot remember to export yout telegram token and your allowed users like this:

```sh
export TELEGRAM_TOKEN=your:private_token
export TELEGRAM_VALID_USERS=user1,user2,user9
```

## Compilation for Raspberry Pi 1
This command will build for `ARM-6`. The `ldflags` `-s` and `-w` strip the debugging information and reduce the size of the executable.

```sh
GOOS=linux GOARCH=arm GOARM=6 go build  -ldflags="-s -w" -v main.go
```

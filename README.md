Before executing the bot remember to export yout telegram token and your allowed users like this:

```sh
export TELEGRAM_TOKEN=your:private_token
export TELEGRAM_VALID_USERS=user1,user2,user9
export TELEGRAM_CHAT_IDS=123,1234545,123343
```

## Compilation for Raspberry Pi 1

This command will build for `ARM-6`. The `ldflags` `-s` and `-w` strip the debugging information and reduce the size of the executable.

```sh
GOOS=linux GOARCH=arm GOARM=6 go build  -ldflags="-s -w" -v main.go
```

## Deploy as a service

Copy the executable in `/opt/goldorackBot/mansionEyeBot`
Copy the `mansioneyebot.service` file in `/etc/systemd/system/`
Edit `mansioneyebot.service` and replace `TELEGRAM_TOKEN`, `TELEGRAM_VALID_USERS` and `TELEGRAM_CHAT_IDS` by the real values.

Make the file executable:

```sh
sudo chmod 755 /etc/systemd/system/mansioneyebot.service
```

Set auto start to the service:

```sh
sudo systemctl enable mansioneyebot.service
```

Start the serivce with:

```sh
sudo systemctl start mansioneyebot
```

Stop the service with:

```sh
sudo systemctl stop mansioneyebot
```

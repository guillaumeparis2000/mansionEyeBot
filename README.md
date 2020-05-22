Before executing the bot remember to export yout telegram token and your allowed users like this:

```sh
export TELEGRAM_TOKEN=your:private_token
export TELEGRAM_VALID_USERS=user1,user2,user9
export TELEGRAM_CHAT_IDS=123,1234545,123343
```

## Compilation

To compile the binary, you can use the make file like this:

- `make release`: Create a release
- `make release-pi`: Create a release for raspberry pi
- `make debug`: Create a release for debugging
- `make clean`: Clean previous builds

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

package main

import (
	"flag"
	"os"
	"strings"

	"github.com/guillaumeparis2000/mansionEyeBot/pkg/telegrambot"
)

func main() {
	actionPtr := flag.String("action", "service", "action name: service, send-picture")
	picturePtr := flag.String("picture", "picture", "picture path")

	token := *flag.String("token", os.Getenv("TELEGRAM_TOKEN"), "Set bot token, if empty get from env")
	validUsers := strings.Split(*flag.String("valid-users", os.Getenv("TELEGRAM_VALID_USERS"), "Set bot valid users, if empty get from env"), ",")
	chatIds := strings.Split(*flag.String("chat-ids", os.Getenv("TELEGRAM_CHAT_IDS"), "Set bot chat ids, if empty get from env"), ",")
	flag.Parse()

	bot := telegrambot.NewTelegramBot(token, validUsers, chatIds)

	if *actionPtr == "service" {
		bot.HandleService()
	} else if *actionPtr == "send-picture" && *picturePtr != "" {
		picture := *picturePtr
		bot.HandleSendPicture(picture)
	}
}


package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/guillaumeparis2000/mansionEyeBot/internal/pkg/telegrambot"
	"github.com/guillaumeparis2000/mansionEyeBot/internal/version"
)

func main() {
	versionPtr := flag.Bool("version", false, "MansionEyeBot version")
	actionPtr := flag.String("action", "service", "action name: service, send-picture")
	picturePtr := flag.String("picture", "picture", "picture path")

	token := *flag.String("token", os.Getenv("TELEGRAM_TOKEN"), "Set bot token, if empty get from env")
	validUsers := strings.Split(*flag.String("valid-users", os.Getenv("TELEGRAM_VALID_USERS"), "Set bot valid users, if empty get from env"), ",")
	chatIds := strings.Split(*flag.String("chat-ids", os.Getenv("TELEGRAM_CHAT_IDS"), "Set bot chat ids, if empty get from env"), ",")
	flag.Parse()

	bot := telegrambot.NewTelegramBot(token, validUsers, chatIds)

	if *versionPtr == true {
		buildData := version.Get()
		fmt.Printf("version: %s\n", buildData.Version)
		fmt.Printf("Git commit: %s\n", buildData.GitCommit)
		fmt.Printf("Go version: %s\n", buildData.GoVersion)
	} else if *actionPtr == "service" {
		bot.HandleService()
	} else if *actionPtr == "send-picture" && *picturePtr != "" {
		picture := *picturePtr
		bot.HandleSendPicture(picture)
	}
}


package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guillaumeparis2000/mansionEyeBot/internal/pkg/api"
	"github.com/guillaumeparis2000/mansionEyeBot/internal/pkg/telegrambot"
	"github.com/guillaumeparis2000/mansionEyeBot/internal/version"
)

func main() {
	versionPtr := flag.Bool("version", false, "MansionEyeBot version")

	token := os.Getenv("TELEGRAM_TOKEN")
	validUsers := strings.Split(os.Getenv("TELEGRAM_VALID_USERS"), ",")
	chatIds := strings.Split(os.Getenv("TELEGRAM_CHAT_IDS"), ",")

	flag.Parse()

	if *versionPtr == true {
		buildData := version.Get()
		fmt.Printf("version: %s\n", buildData.Version)
		fmt.Printf("Git commit: %s\n", buildData.GitCommit)
		fmt.Printf("Go version: %s\n", buildData.GoVersion)
	} else {
		bot := telegrambot.NewTelegramBot(token, validUsers, chatIds)
		api := api.Initialize(bot)
		api.Run()

		_, err := bot.HandleService()
		if err != nil {
			log.Print("Telegram not starting!")
			log.Fatal(err)
		}

	}
}

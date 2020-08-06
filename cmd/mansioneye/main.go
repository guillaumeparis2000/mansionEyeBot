package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/facebookgo/pidfile"
	"github.com/guillaumeparis2000/mansionEyeBot/internal/pkg/api"
	"github.com/guillaumeparis2000/mansionEyeBot/internal/pkg/telegrambot"
	"github.com/guillaumeparis2000/mansionEyeBot/internal/pkg/yeelight"
	"github.com/guillaumeparis2000/mansionEyeBot/internal/version"
)

func main() {
	versionPtr := flag.Bool("version", false, "MansionEyeBot version")

	token := os.Getenv("TELEGRAM_TOKEN")
	validUsers := strings.Split(os.Getenv("TELEGRAM_VALID_USERS"), ",")
	chatIds := strings.Split(os.Getenv("TELEGRAM_CHAT_IDS"), ",")

	desk := os.Getenv("YEELIGHT_DESK")
	salon := os.Getenv("YEELIGHT_SALON")

	flag.Parse()

	if *versionPtr == true {
		buildData := version.Get()
		fmt.Printf("version: %s\n", buildData.Version)
		fmt.Printf("Git commit: %s\n", buildData.GitCommit)
		fmt.Printf("Go version: %s\n", buildData.GoVersion)
	} else {
		// create pid file
		pidfile.SetPidfilePath("/var/run/mansioneye-bot.pid")

		piderror := pidfile.Write()
		if piderror != nil {
			panic("Could not write pid file")
		}

		yeelights := yeelight.NewYeelights(desk, salon)
		botConfig := telegrambot.NewTelegramBot(token, validUsers, chatIds, yeelights)

		api := api.Initialize(botConfig)
		go http.ListenAndServe(":8001", api.Router)
		log.Print("Rest API started on port 8001")

		log.Print("Starting Telegram Bot...")
		botConfig.HandleService()
		err := botConfig.Bot.Start()
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Telegram Bot successfully started!")
	}
}

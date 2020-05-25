package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guillaumeparis2000/mansionEyeBot/internal/pkg/telegrambot"
	"github.com/guillaumeparis2000/mansionEyeBot/internal/version"
	"github.com/joho/godotenv"
)

func main() {
	versionPtr := flag.Bool("version", false, "MansionEyeBot version")
	actionPtr := flag.String("action", "service", "action name: service, send-picture")
	picturePtr := flag.String("picture", "picture", "picture path")

	token := goDotEnvVariable("TELEGRAM_TOKEN")
	validUsers := strings.Split(goDotEnvVariable("TELEGRAM_VALID_USERS"), ",")
	chatIds := strings.Split(goDotEnvVariable("TELEGRAM_CHAT_IDS"), ",")

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

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {
	var configFile = "/etc/mansionEyeBot/config"
	env := os.Getenv("ENV")
	if "dev" == env {
		configFile = "./config"
	}

	// load .env file
	err := godotenv.Load(configFile)

	if err != nil {
	  log.Fatalf("Error loading %s file", configFile)
	}

	return os.Getenv(key)
  }

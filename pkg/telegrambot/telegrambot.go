package telegrambot

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/guillaumeparis2000/mansionEyeBot/pkg/motioneye"
	"github.com/yanzay/tbot/v2"
)

// Botconfig Telegram bot configuration.
type Botconfig struct {
	token string
	validUsers []string
	chatIds []string
	client *tbot.Client
	bot *tbot.Server
}

// NewTelegramBot telegram bot constructor to initialize the bot
func NewTelegramBot(token string, validUsers []string, chatIds []string) *Botconfig{
	app := &Botconfig{}

	app.token = token
	app.validUsers = validUsers
	app.chatIds = chatIds

	app.bot = tbot.New(app.token)

	// Use validUsers for Auth middleware, allow to interact only with user1 and user2
	app.bot.Use(app.auth)
	app.client = app.bot.Client()

	return app
}

// Middleware to only allow validated user to chat with the bot.
func (bc *Botconfig) auth(h tbot.UpdateHandler) tbot.UpdateHandler {
	return func(u *tbot.Update) {
		if contains(bc.validUsers, u.Message.Chat.Username) {
			h(u)
		} else {
			bc.sendResponse(u.Message, "You are not allowed to use this Bot!")
		}
	}
}

// HandleService start the telegram bot.
func (bc *Botconfig) HandleService() {
	bc.bot.HandleMessage("/status", bc.statusHandler)
	bc.bot.HandleMessage("/pause", bc.pauseHandler)
	bc.bot.HandleMessage("/resume", bc.resumeHandler)
	bc.bot.HandleMessage("/check", bc.checkHandler)
	bc.bot.HandleMessage("/time", bc.timeHandler)
	bc.bot.HandleMessage("/video", bc.videoHandler)
	bc.bot.HandleMessage("/snapshot", bc.snapShotHandler)
	bc.bot.HandleMessage("/get_my_id", bc.getMyIDHandler)
	bc.bot.HandleMessage("/valid_users", bc.validUsersHandler)
}

// HandleSendPicture allow to send a picture with the bot to all chat ids defined.
func (bc *Botconfig) HandleSendPicture(picture string) {
	for _, chatID := range bc.chatIds {
		_, err := bc.client.SendPhotoFile(chatID, picture, tbot.OptCaption("Motion Detected"))
		if err != nil {
			bc.client.SendMessage(chatID, err.Error())
		}
	}
}

func (bc *Botconfig) statusHandler(m *tbot.Message) {
	result := motioneye.Status()
	bc.sendResponse(m, result)
}

func (bc *Botconfig) pauseHandler(m *tbot.Message) {
	result := motioneye.Pause()
	bc.sendResponse(m, result)
}

func (bc *Botconfig) resumeHandler(m *tbot.Message) {
	result := motioneye.Resume()
	bc.sendResponse(m, result)
}

func (bc *Botconfig) checkHandler(m *tbot.Message) {
	result := motioneye.Check()
	bc.sendResponse(m, result)
}

func (bc *Botconfig) timeHandler(m *tbot.Message) {
	result := motioneye.Time()
	bc.sendResponse(m, result)
}

func (bc *Botconfig) videoHandler(m *tbot.Message) {
	bc.sendResponse(m, "This operation can be long! Please wait.")

	videos := motioneye.GetLastVideos()

	if len(videos) > 0 {
		bc.client.SendChatAction(m.Chat.ID, tbot.ActionUploadVideo)
		var lastVideo = filepath.Base(videos[0])
		_, err := bc.client.SendVideoFile(m.Chat.ID, lastVideo, tbot.OptCaption(videos[0]))
		if err != nil {
			bc.client.SendMessage(m.Chat.ID, err.Error())
		}
	} else {
		bc.client.SendMessage(m.Chat.ID, fmt.Sprintf("No video available!"))
	}
}

func (bc *Botconfig) snapShotHandler(m *tbot.Message) {
	result := motioneye.SnapShot()
	bc.sendResponse(m, result)
}

// getMyIDHandler Return the user telegram chat ID
func (bc *Botconfig) getMyIDHandler(m *tbot.Message) {
	bc.sendResponse(m, m.Chat.ID)
}

// validUsersHandler Return the list a allowed users
func (bc *Botconfig) validUsersHandler(m *tbot.Message) {
	bc.sendResponse(m, strings.Join(bc.validUsers, ", "))
}

// Send response messages to the chat
// In case of error, send the error as a text message.
func (bc *Botconfig) sendResponse(m *tbot.Message, response string) {
	arrAction := bc.client.SendChatAction(m.Chat.ID, tbot.ActionTyping)
	if arrAction != nil {
		log.Println(arrAction.Error())
	}

	_, errMsg := bc.client.SendMessage(m.Chat.ID, response)
	if errMsg != nil {
		log.Println(errMsg.Error())
	}
}

// Search if an array contain a given string
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

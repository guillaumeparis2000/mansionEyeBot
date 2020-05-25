package telegrambot

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/guillaumeparis2000/mansionEyeBot/internal/pkg/motioneye"
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
func (bc *Botconfig) HandleService() (bool, error) {
	bc.bot.HandleMessage("/status", bc.statusHandler)
	bc.bot.HandleMessage("/pause", bc.pauseHandler)
	bc.bot.HandleMessage("/resume", bc.resumeHandler)
	bc.bot.HandleMessage("/check", bc.checkHandler)
	bc.bot.HandleMessage("/time", bc.timeHandler)
	bc.bot.HandleMessage("/video", bc.videoHandler)
	bc.bot.HandleMessage("/snapshot", bc.snapShotHandler)
	bc.bot.HandleMessage("/get_my_id", bc.getMyIDHandler)
	bc.bot.HandleMessage("/valid_users", bc.validUsersHandler)

	err := bc.bot.Start()
	if err != nil {
		return false, err
	}

	log.Print("Telegram Bot successfully started!")

	return true, nil
}

// HandleSendPicture allow to send a picture with the bot to all chat ids defined.
func (bc *Botconfig) HandleSendPicture(picture string, name string) {
	for _, chatID := range bc.chatIds {
		_, err := bc.client.SendPhotoFile(chatID, picture, tbot.OptCaption("Motion Detected from " + name))
		if err != nil {
			bc.client.SendMessage(chatID, err.Error())
		}
	}
}

func (bc *Botconfig) statusHandler(m *tbot.Message) {
	result, err := motioneye.Status()
	bc.sendResponseOrError(m, result, err)
}

func (bc *Botconfig) pauseHandler(m *tbot.Message) {
	result, err := motioneye.Pause()
	bc.sendResponseOrError(m, result, err)
}

func (bc *Botconfig) resumeHandler(m *tbot.Message) {
	result, err := motioneye.Resume()
	bc.sendResponseOrError(m, result, err)
}

func (bc *Botconfig) checkHandler(m *tbot.Message) {
	result, err := motioneye.Check()
	bc.sendResponseOrError(m, result, err)
}

func (bc *Botconfig) timeHandler(m *tbot.Message) {
	result := motioneye.Time()
	bc.sendResponse(m, result)
}

func (bc *Botconfig) snapShotHandler(m *tbot.Message) {
	result, err := motioneye.SnapShot()
	bc.sendResponseOrError(m, result, err)
}

func (bc *Botconfig) sendResponseOrError(m *tbot.Message, result string, err error) {
	if err != nil {
		bc.sendResponse(m, err.Error())
	} else {
		bc.sendResponse(m, result)
	}
}

func (bc *Botconfig) videoHandler(m *tbot.Message) {
	bc.sendResponse(m, "This operation can be long! Please wait.")

	videos, err := motioneye.GetLastVideos()
	if err != nil {
		bc.client.SendMessage(m.Chat.ID, err.Error())
	} else if len(videos) > 0 {
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
		bc.client.SendMessage(m.Chat.ID, arrAction.Error())
	}

	_, errMsg := bc.client.SendMessage(m.Chat.ID, response)
	if errMsg != nil {
		bc.client.SendMessage(m.Chat.ID, errMsg.Error())
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

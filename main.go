package main

import (
	"github.com/yanzay/tbot"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type application struct {
	client *tbot.Client
}

var validUsers = []string{"user1", "user2"}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func main() {
	bot := tbot.New(os.Getenv("TELEGRAM_TOKEN"))
	app := &application{}

	// Use validUsers for Auth middleware, allow to interact only with user1 and user2
	bot.Use(app.auth)

	app.client = bot.Client()
	bot.HandleMessage("/status", app.statusHandler)
	bot.HandleMessage("/pause", app.pauseHandler)
	bot.HandleMessage("/resume", app.resumeHandler)
	bot.HandleMessage("/check", app.checkHandler)
	bot.HandleMessage("/time", app.timeHandler)
	bot.HandleMessage("/video", app.videoHandler)
	bot.HandleMessage("/snapshot", app.snapShotHandler)

	err := bot.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func (a *application) auth(h tbot.UpdateHandler) tbot.UpdateHandler {
	return func(u *tbot.Update) {
		if contains(validUsers, u.Message.Chat.Username) {
			h(u)
		} else {
			a.client.SendMessage(u.Message.Chat.ID, "You are not allowed to use this Bot!")
		}
	}
}

func (a *application) statusHandler(m *tbot.Message) {
	result := webControl("detection", "status")
	a.client.SendChatAction(m.Chat.ID, tbot.ActionTyping)
	a.client.SendMessage(m.Chat.ID, result)
}

func (a *application) pauseHandler(m *tbot.Message) {
	result := webControl("detection", "pause")
	a.client.SendChatAction(m.Chat.ID, tbot.ActionTyping)
	a.client.SendMessage(m.Chat.ID, result)
}

func (a *application) resumeHandler(m *tbot.Message) {
	result := webControl("detection", "start")
	a.client.SendChatAction(m.Chat.ID, tbot.ActionTyping)
	a.client.SendMessage(m.Chat.ID, result)
}

func (a *application) checkHandler(m *tbot.Message) {
	result := webControl("detection", "connection")
	a.client.SendChatAction(m.Chat.ID, tbot.ActionTyping)
	a.client.SendMessage(m.Chat.ID, result)
}

func (a *application) snapShotHandler(m *tbot.Message) {
	result := webControl("action", "snapshot")
	a.client.SendChatAction(m.Chat.ID, tbot.ActionUploadPhoto)
	a.client.SendMessage(m.Chat.ID, result)
}

func (a *application) timeHandler(m *tbot.Message) {
	result := time.Now().Format(time.ANSIC)
	a.client.SendChatAction(m.Chat.ID, tbot.ActionTyping)
	a.client.SendMessage(m.Chat.ID, result)
}

func (a *application) videoHandler(m *tbot.Message) {
	a.client.SendMessage(m.Chat.ID, "This operation can be long! Please wait.")
	dir := `/var/lib/motioneye/Camera1/`

	files, _ := ioutil.ReadDir(dir)
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Unix() > files[j].ModTime().Unix()
	})

	var videos []string
	for _, file := range files {
		r := strings.HasSuffix(file.Name(), ".mp4")
		if r {
			videos = append(videos, file.Name())
		}
	}

	if len(videos) > 0 {
		var lastVideo = dir + videos[0]
		a.client.SendMessage(m.Chat.ID, "Uploading Video, Please wait a little more. Thanks")
		a.client.SendChatAction(m.Chat.ID, tbot.ActionUploadVideo)
		_, err := a.client.SendVideoFile(m.Chat.ID, lastVideo)
		if err != nil {
			a.client.SendMessage(m.Chat.ID, err.Error())
		}
	} else {
		a.client.SendMessage(m.Chat.ID, "No video available!")
	}
}

func webControl(cmdType string, cmd string) string {
	// Make a get request
	rs, err := http.Get("http://localhost:7999/0/" + cmdType + "/" + cmd)
	if err != nil {
		return err.Error()
	}
	defer rs.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return err.Error()
	}

	bodyString := string(bodyBytes)

	if len(bodyString) > 0 {
		return bodyString
	}
	return "OK"
}

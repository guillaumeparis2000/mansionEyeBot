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

// Search if an array contain a given string
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

// Middleware to only allow validated user to chat with the bot.
func (a *application) auth(h tbot.UpdateHandler) tbot.UpdateHandler {
	return func(u *tbot.Update) {
		if contains(validUsers, u.Message.Chat.Username) {
			h(u)
		} else {
			a.sendResponse(u.Message, "You are not allowed to use this Bot!")
		}
	}
}

// Return the current status of the camera.
func (a *application) statusHandler(m *tbot.Message) {
	result := webControl("detection", "status")
	a.sendResponse(m, result)
}

// Pause the motion detection of the camera.
func (a *application) pauseHandler(m *tbot.Message) {
	result := webControl("detection", "pause")
	a.sendResponse(m, result)
}

// Resume the motion detection of the camera.
func (a *application) resumeHandler(m *tbot.Message) {
	result := webControl("detection", "start")
	a.sendResponse(m, result)
}

// Return the connection status of the camera.
func (a *application) checkHandler(m *tbot.Message) {
	result := webControl("detection", "connection")
	a.sendResponse(m, result)
}

// Create a snapshot.
func (a *application) snapShotHandler(m *tbot.Message) {
	result := webControl("action", "snapshot")
	a.sendResponse(m, result)
}

// Return the current time in the server.
func (a *application) timeHandler(m *tbot.Message) {
	result := time.Now().Format(time.ANSIC)
	a.sendResponse(m, result)
}

// Send the last recorded video to the chat.
func (a *application) videoHandler(m *tbot.Message) {
	a.sendResponse(m, "This operation can be long! Please wait.")
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
		a.client.SendChatAction(m.Chat.ID, tbot.ActionUploadVideo)
		var lastVideo = dir + videos[0]
		_, err := a.client.SendVideoFile(m.Chat.ID, lastVideo, tbot.OptCaption(videos[0]))
		if err != nil {
			a.client.SendMessage(m.Chat.ID, err.Error())
		}
	} else {
		a.client.SendMessage(m.Chat.ID, "No video available!")
	}
}

// Make the Api request to the motion server.
// Return the result of each command or the error in text format.
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

// Send response messages to the chat
// In case of error, send the error as a text message.
func (a *application) sendResponse(m *tbot.Message, response string) {
	arrAction := a.client.SendChatAction(m.Chat.ID, tbot.ActionTyping)
	if arrAction != nil {
		log.Println(arrAction.Error())
	}

	_, errMsg := a.client.SendMessage(m.Chat.ID, response)
	if errMsg != nil {
		log.Println(errMsg.Error())
	}
}

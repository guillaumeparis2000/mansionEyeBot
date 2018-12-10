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

func main() {
	bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	// Use whitelist for Auth middleware, allow to interact only with user1 and user2
	whitelist := []string{"youUser"}
	bot.AddMiddleware(tbot.NewAuth(whitelist))

	bot.HandleFunc("/status", statusHandler)
	bot.HandleFunc("/pause", pauseHandler)
	bot.HandleFunc("/resume", resumeHandler)
	bot.HandleFunc("/check", checkHandler)
	bot.HandleFunc("/time", timeHandler)
	bot.HandleFunc("/video", videoHandler)
	bot.HandleFunc("/snapshot", snapShotHandler)
	bot.ListenAndServe()
}

func statusHandler(m *tbot.Message) {
	result := webControl("detection", "status")
	m.Reply(result)
}

func pauseHandler(m *tbot.Message) {
	result := webControl("detection", "pause")
	m.Reply(result)
}

func resumeHandler(m *tbot.Message) {
	result := webControl("detection", "start")
	m.Reply(result)
}

func checkHandler(m *tbot.Message) {
	result := webControl("detection", "connection")
	m.Reply(result)
}

func snapShotHandler(m *tbot.Message) {
	result := webControl("action", "snapshot")
	m.Reply(result)
}

func timeHandler(m *tbot.Message) {
	result := time.Now().Format(time.ANSIC)
	m.Reply(result)
}

func videoHandler(m *tbot.Message) {
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
		m.ReplyVideo(lastVideo)
	} else {
		m.Reply("No video available!")
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
		panic(err)
	}

	bodyString := string(bodyBytes)
	return bodyString
}

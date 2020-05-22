package motioneye

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"time"
)
// Status return the current status of the camera.
func Status() string {
	result := webControl("detection", "status")
	return result
}

// Pause the motion detection of the camera.
func Pause() string {
	result := webControl("detection", "pause")
	return result
}

// Resume the motion detection of the camera.
func Resume() string {
	result := webControl("detection", "start")
	return result
}

// Check Return the connection status of the camera.
func Check() string {
	result := webControl("detection", "connection")
	return result
}

// SnapShot Create a snapshot.
func SnapShot() string{
	result := webControl("action", "snapshot")
	return result
}

// Time Return the current time in the server.
func Time() string {
	result := time.Now().Format(time.ANSIC)
	return result
}

// GetLastVideos Return the last recorded videos.
func GetLastVideos() []string {
	rootDir := `/var/lib/motioneye/`

	dirs, err := ioutil.ReadDir(rootDir)
	if err != nil {
		log.Fatal(err)
	}

	var videos []string
	for _, dir := range dirs {
		if dir.IsDir() {
			camDir := filepath.Join(rootDir, dir.Name())
			files, _ := ioutil.ReadDir(camDir)
			sort.Slice(files, func(i, j int) bool {
				return files[i].ModTime().Unix() > files[j].ModTime().Unix()
			})

			for _, file := range files {
				r := strings.HasSuffix(file.Name(), ".mp4")
				if r {
					videos = append(videos, filepath.Join(camDir, file.Name()))
				}
			}
		}
	}

	return videos
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

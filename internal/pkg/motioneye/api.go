package motioneye

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"time"
)
// Status return the current status of the camera.
func Status() (string, error) {
	result, err := restAPI("detection", "status")
	return result, err
}

// Pause the motion detection of the camera.
func Pause() (string, error) {
	result, err := restAPI("detection", "pause")
	return result, err
}

// Resume the motion detection of the camera.
func Resume() (string, error) {
	result, err := restAPI("detection", "start")
	return result, err
}

// Check Return the connection status of the camera.
func Check() (string, error) {
	result, err := restAPI("detection", "connection")
	return result, err
}

// SnapShot Create a snapshot.
func SnapShot() (string, error) {
	result, err := restAPI("action", "snapshot")
	return result, err
}

// Time Return the current time in the server.
func Time() string {
	result := time.Now().Format(time.ANSIC)
	return result
}

// GetLastVideos Return the last recorded videos.
func GetLastVideos() ([]string, error) {
	var videos []string
	rootDir := `/var/lib/motioneye/`

	dirs, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return videos, err
	}

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

	return videos, nil
}

// Make the Api request to the motion server.
// Return the result of each command or the error in text format.
func restAPI(cmdType string, cmd string) (string, error) {
	// Make a get request
	rs, err := http.Get("http://localhost:7999/0/" + cmdType + "/" + cmd)
	if err != nil {
		return "", err
	}
	defer rs.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return "", err
	}

	bodyString := string(bodyBytes)

	if len(bodyString) > 0 {
		return bodyString, nil
	}

	return "OK", nil
}

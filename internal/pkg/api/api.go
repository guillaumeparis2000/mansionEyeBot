package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guillaumeparis2000/mansionEyeBot/internal/pkg/telegrambot"
)

// API configuration.
type API struct {
	Router *mux.Router
	bot *telegrambot.Botconfig
}

//PictureData contain path and name.
type PictureData struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Sent bool `json:"sent"`
}

// Initialize the API.
func Initialize(bot *telegrambot.Botconfig) *API {
	a := &API{}
	a.bot = bot
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/picture", a.handleSendPicture).Methods("POST")

	return a
}

func (a *API) handleSendPicture(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var picture PictureData
	_ = json.NewDecoder(r.Body).Decode(&picture)

	log.Printf("Send picture %s from %s", picture.Path, picture.Name)

	a.bot.HandleSendPicture(picture.Path, picture.Name)
	picture.Sent = true

	json.NewEncoder(w).Encode(picture)
}

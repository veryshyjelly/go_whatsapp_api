package routes

import (
	"github.com/gorilla/mux"
	"go_whatsapp_api/api/pkg/controllers"
	"log"
	"net/http"
)

var RegisterWApiRoutes = func(router *mux.Router) {
	router.HandleFunc("/sendMessage/", func(w http.ResponseWriter, r *http.Request) {
		err := controllers.SendMessage(w, r)
		if err != nil {
			log.Println(err)
		}
	}).Methods("POST")
	router.HandleFunc("/sendPhoto/", func(w http.ResponseWriter, r *http.Request) {
		err := controllers.SendPhoto(w, r)
		if err != nil {
			log.Println(err)
		}
	}).Methods("POST")
	router.HandleFunc("/sendAudio/", func(w http.ResponseWriter, r *http.Request) {
		err := controllers.SendAudio(w, r)
		if err != nil {
			log.Println(err)
		}
	}).Methods("POST")
	router.HandleFunc("/sendDocument/", func(w http.ResponseWriter, r *http.Request) {
		err := controllers.SendDocument(w, r)
		if err != nil {
			log.Println(err)
		}
	}).Methods("POST")
	router.HandleFunc("/sendVideo/", func(w http.ResponseWriter, r *http.Request) {
		err := controllers.SendVideo(w, r)
		if err != nil {
			log.Println(err)
		}
	}).Methods("POST")
	router.HandleFunc("/sendContact/", func(w http.ResponseWriter, r *http.Request) {
		err := controllers.SendContact(w, r)
		if err != nil {
			log.Println(err)
		}
	}).Methods("POST")
	router.HandleFunc("/sendSticker/", func(w http.ResponseWriter, r *http.Request) {
		err := controllers.SendSticker(w, r)
		if err != nil {
			log.Println(err)
		}
	}).Methods("POST")
}

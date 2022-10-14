package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"go_whatsapp_bot/api/pkg/routes"
	"go_whatsapp_bot/app/pkg/models"
	"log"
	"net/http"
)

func APIStart(port string) {
	var r = mux.NewRouter()
	routes.RegisterWApiRoutes(r)
	http.Handle("/", r)
	fmt.Println("server started at http://localhost" + port)
	log.Fatalln(http.ListenAndServe(port, nil))
	models.Client.Disconnect()
}

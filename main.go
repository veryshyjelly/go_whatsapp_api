package main

import (
	"go_whatsapp_bot/app"
	"go_whatsapp_bot/handler"
)

func main() {
	app.Connect(handler.EventHandler, "INFO")
}

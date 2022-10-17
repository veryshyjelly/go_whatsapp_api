package main

import (
	"go_whatsapp_api/api/cmd"
	config2 "go_whatsapp_api/app/cmd/config"
	"go_whatsapp_api/app/pkg/connect"
)

func main() {
	go connect.Connect("INFO")
	cmd.APIStart(config2.WhatsappPort)
}

package main

import (
	"go_whatsapp_api/api/cmd"
	"go_whatsapp_api/app/pkg/config"
)

func main() {
	go config.Connect("INFO")
	cmd.APIStart(":8080")
}

package config

import (
	"context"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go_whatsapp_bot/app/pkg/controllers"
	"go_whatsapp_bot/app/pkg/models"
)

func Connect(minLevel string) {
	fmt.Println(sqlite3.Version())
	dbLog := waLog.Stdout("Database", minLevel, true)

	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", minLevel, true)

	models.Client = whatsmeow.NewClient(deviceStore, clientLog)
	models.Client.AddEventHandler(controllers.EventHandler)

	if models.Client.Store.ID == nil { // No ID stored, new login
		qrChan, _ := models.Client.GetQRChannel(context.Background())
		err = models.Client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrcode.WriteFile(evt.Code, qrcode.Medium, 256, "scan.png")
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else { // Already logged in, just connect
		err = models.Client.Connect()
		if err != nil {
			panic(err)
		}
	}

}

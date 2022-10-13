package app

import (
	"context"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"os"
	"os/signal"
	"syscall"
)

var Client *whatsmeow.Client

func Connect(handler whatsmeow.EventHandler, minLevel string) {
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
	Client = whatsmeow.NewClient(deviceStore, clientLog)
	Client.AddEventHandler(handler)

	if Client.Store.ID == nil { // No ID stored, new login
		qrChan, _ := Client.GetQRChannel(context.Background())
		err = Client.Connect()
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
		err = Client.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	Client.Disconnect()
}

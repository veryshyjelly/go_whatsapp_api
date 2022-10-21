package connect

import (
	"context"
	"encoding/json"
	"fmt"
	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go_whatsapp_api/app/cmd/config"
	"go_whatsapp_api/app/pkg/controllers"
	"go_whatsapp_api/app/pkg/models"
	"io/ioutil"
	"log"
)

func Connect(minLevel string) {
	_, _, _ = sqlite3.Version()
	dbLog := waLog.Stdout("Database", minLevel, true)

	container, err := sqlstore.New("sqlite3", "file:whatsapp.db?_foreign_keys=on", dbLog)
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
				tObj := qrcodeTerminal.New()
				tObj.Get(evt.Code).Print()
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

	bridgeDataFile, err := ioutil.ReadFile("./app/cmd/config/bridge.json")
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(bridgeDataFile, &config.BridgingData)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v", config.BridgingData)
}

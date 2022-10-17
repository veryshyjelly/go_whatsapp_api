package controllers

import (
	"go.mau.fi/whatsmeow/types/events"
	"go_whatsapp_api/app/cmd/config"
	"go_whatsapp_api/app/cmd/routes"
	"go_whatsapp_api/app/pkg/models"
	"go_whatsapp_api/app/pkg/utils"
	"log"
)

func EventHandler(evt interface{}) {
	switch m := evt.(type) {
	case *events.Message:
		go func(m *events.Message) {
			if m == nil {
				return
			}

			if m.Info.IsGroup {
				info, err := models.Client.GetGroupInfo(m.Info.Chat)
				if err != nil {
					log.Println(err)
				}
				utils.PrintMessage(m, info)
			} else {
				utils.PrintMessage(m, nil)
			}

			if brg, ok := config.BridgingData[m.Info.Chat.String()]; ok && brg.PermTG && brg.PermWA {
				err := routes.HandleMessage(m, brg.Destination)
				if err != nil {
					log.Println(err)
				}
			}
		}(m)
	case *events.Receipt:
		//fmt.Println(m.Type.GoString(), "receipt from", m.Sender.User, m.MessageSource)
	}
}

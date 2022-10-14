package controllers

import (
	"go.mau.fi/whatsmeow/types/events"
	"go_whatsapp_api/app/pkg/models"
	"go_whatsapp_api/app/pkg/utils"
	"log"
)

func EventHandler(evt interface{}) {
	switch m := evt.(type) {
	case *events.Message:
		if m.Info.IsGroup {
			info, err := models.Client.GetGroupInfo(m.Info.Chat)
			if err != nil {
				log.Println(err)
			}
			utils.PrintMessage(m, info)
		} else {
			utils.PrintMessage(m, nil)
		}

	case *events.Receipt:
		//fmt.Println(m.Type.GoString(), "receipt from", m.Sender.User, m.MessageSource)
	}
}

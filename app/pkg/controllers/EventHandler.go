package controllers

import (
	"fmt"
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
				//if m.Message.Conversation != nil {
				//	//x, _ := models.Client.GetGroupInfo(m.Info.Chat)
				//	fmt.Println(models.Client.SendMessage(context.Background(), m.Info.Chat, "", &proto.Message{
				//
				//		ExtendedTextMessage: &proto.ExtendedTextMessage{
				//			Text: proto2.String(""),
				//			ContextInfo: &proto.ContextInfo{
				//
				//				Participant: proto2.String("0@s.whatsapp.net"),
				//				QuotedMessage: &proto.Message{
				//					Conversation: proto2.String("bc"),
				//				},
				//			},
				//		},
				//	}))
				//}
				fmt.Printf("%#v\n", m.Info.ID)
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

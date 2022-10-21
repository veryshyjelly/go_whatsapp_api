package routes

import (
	"go.mau.fi/whatsmeow/types/events"
	"go_whatsapp_api/app/cmd/controllers"
)

func HandleMessage(m *events.Message, dest string) error {
	mess, pushName := m.Message, m.Info.PushName

	if mess.Conversation != nil {
		return controllers.HandleText(mess, pushName, dest)
	} else if mess.ImageMessage != nil {
		return controllers.HandlePhoto(mess.ImageMessage, pushName, dest)
	} else if mess.VideoMessage != nil {
		return controllers.HandleVideo(mess.VideoMessage, pushName, dest)
	} else if mess.DocumentMessage != nil {
		return controllers.HandleDocument(mess.DocumentMessage, pushName, dest)
	} else if mess.AudioMessage != nil {
		return controllers.HandleAudio(mess.AudioMessage, pushName, dest)
	} else if mess.StickerMessage != nil {
		//fmt.Printf("%#v\n", m.Message.StickerMessage.ContextInfo)
		return controllers.HandleSticker(mess.StickerMessage, dest)
	} else if mess.GetExtendedTextMessage().GetText() != "" {
		return controllers.HandleText(mess, pushName, dest)
	}

	return nil
}

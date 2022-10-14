package utils

import (
	"fmt"
	"github.com/fatih/color"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"log"
)

func PrintMessage(m *events.Message, info *types.GroupInfo) {
	a := color.New(color.FgBlack).Add(color.BgWhite)
	b := color.New(color.FgBlack).Add(color.BgGreen)
	c := color.New(color.FgBlack).Add(color.BgBlue)
	d := color.New(color.FgWhite).Add(color.BgYellow)

	if m.Message.GetExtendedTextMessage().GetContextInfo().GetIsForwarded() {
		if _, err := a.Print("[FORWARDED MESSAGE] "); err != nil {
			log.Println(err)
		}
	} else {
		if _, err := a.Print("[MESSAGE] "); err != nil {
			log.Println(err)
		}
	}

	if _, err := b.Print(" " + m.Info.Timestamp.Format("02/01/2006 15:04:05") + " "); err != nil {
		log.Println(err)
	}
	if _, err := c.Print(" "+m.Message.GetConversation(), m.Message.ExtendedTextMessage.GetText()+" "); err != nil {
		log.Println(err)
	}
	if _, err := d.Print(m.Info.MediaType); err != nil {
		log.Println(err)
	}
	fmt.Println()

	fmt.Println(color.MagentaString("=> From"), color.GreenString(m.Info.PushName), color.YellowString(m.Info.Sender.String()))

	fmt.Print(color.HiBlueString("=> In "))
	if info != nil {
		fmt.Println(color.GreenString("%v %v", info.Name, m.Info.Chat.String()))
	} else {
		fmt.Println(color.GreenString("Private Chat %v", m.Info.Chat.String()))
	}
}

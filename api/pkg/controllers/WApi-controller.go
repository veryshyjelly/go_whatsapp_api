package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go_whatsapp_bot/api/pkg/models"
	"google.golang.org/protobuf/proto"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	apiModels "go_whatsapp_bot/app/pkg/models"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	data := models.MsgRequest{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusPreconditionFailed)
		b, _ := json.Marshal(err)
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
		return
	}
	var us = strings.Split(data.ChatId, "@")
	jid := types.JID{
		User:   us[0],
		Server: us[1],
	}

	msg := waProto.Message{
		Conversation: &data.Text,
	}
	fmt.Println(apiModels.Client.IsConnected())
	resp, err := apiModels.Client.SendMessage(context.Background(), jid, "", &msg)

	if err != nil {
		fmt.Println("Error sending message")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(resp)
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
	}
}

func SendPhoto(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024 << 20)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	downFile, _, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error retrieving data from form-data")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fileBytes, err := ioutil.ReadAll(downFile)

	resp, err := apiModels.Client.Upload(context.Background(), fileBytes, whatsmeow.MediaImage)
	if err != nil {
		fmt.Println("Error while uploading")
		fmt.Println(err)
		w.WriteHeader(http.StatusNotExtended)
		return
	}

	imageMsg := &waProto.ImageMessage{
		Url:           &resp.URL,
		Mimetype:      proto.String("image/png"), // replace this with the actual mime type
		Caption:       proto.String(r.FormValue("caption")),
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		DirectPath:    &resp.DirectPath,
	}

	var us = strings.Split(r.FormValue("chat_id"), "@")
	jid := types.JID{
		User:   us[0],
		Server: us[1],
	}
	rsp, err := apiModels.Client.SendMessage(context.Background(), jid, "", &waProto.Message{
		ImageMessage: imageMsg,
	})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotImplemented)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(rsp)
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
	}
}

func SendAudio(w http.ResponseWriter, r *http.Request) {

}

func SendDocument(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println(err)
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error retrieving data from form-data")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fileBytes, err := ioutil.ReadAll(file)

	var us = strings.Split(r.FormValue("chat_id"), "@")
	resp, err := apiModels.Client.Upload(context.Background(), fileBytes, whatsmeow.MediaDocument)
	ioutil.WriteFile(handler.Filename, fileBytes, fs.ModeAppend)

	if err != nil {
		fmt.Println("Error uploading file")
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		fmt.Println("Uploaded file:", handler.Filename)
	}

	jid := types.JID{
		User:   us[0],
		Server: us[1],
	}
	fmt.Println(jid)
	var mimeType = http.DetectContentType(fileBytes)
	fmt.Println("Mime type: ", mimeType)

	var msg = waProto.DocumentMessage{
		Url:           &resp.URL,
		Mimetype:      proto.String(mimeType),
		Title:         proto.String(r.FormValue("title")),
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
		MediaKey:      resp.MediaKey,
		FileName:      proto.String(handler.Filename),
		FileEncSha256: resp.FileEncSHA256,
		DirectPath:    &resp.DirectPath,
		Caption:       proto.String(r.FormValue("caption")),
	}

	rsp, err := apiModels.Client.SendMessage(context.Background(), jid, "", &waProto.Message{
		Conversation:    proto.String("xyz"),
		DocumentMessage: &msg,
	})
	if err != nil {
		fmt.Println("Error sending message")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(rsp)
		fmt.Println(rsp.ID)
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
	}
}

func SendVideo(w http.ResponseWriter, r *http.Request) {

}

func SendAnimation(w http.ResponseWriter, r *http.Request) {

}

func SendSticker(w http.ResponseWriter, r *http.Request) {

}

func SendVoice(w http.ResponseWriter, r *http.Request) {

}

func SendContact(w http.ResponseWriter, r *http.Request) {

}

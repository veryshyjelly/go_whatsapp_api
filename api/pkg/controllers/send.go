package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go_whatsapp_bot/api/pkg/models"
	"google.golang.org/protobuf/proto"
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
		w.WriteHeader(http.StatusNotAcceptable)
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
		//ExtendedTextMessage: &waProto.ExtendedTextMessage{
		//	ContextInfo: &waProto.ContextInfo{
		//		MentionedJid: data.MentionID,
		//	},
		//},
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
		Mimetype:      proto.String(http.DetectContentType(fileBytes)), // replace this with the actual mime type
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
		b, _ := json.Marshal(err)
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
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
	err := r.ParseMultipartForm(1024 << 20)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	downFile, _, err := r.FormFile("audio")
	if err != nil {
		fmt.Println("Error retrieving data from form-data")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fileBytes, err := ioutil.ReadAll(downFile)

	resp, err := apiModels.Client.Upload(context.Background(), fileBytes, whatsmeow.MediaAudio)
	if err != nil {
		fmt.Println("Error while uploading")
		fmt.Println(err)
		w.WriteHeader(http.StatusNotExtended)
		return
	}

	audioMsg := &waProto.AudioMessage{
		Url:           &resp.URL,
		Mimetype:      proto.String(http.DetectContentType(fileBytes)), // replace this with the actual mime type
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
		AudioMessage: audioMsg,
	})

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotImplemented)
		b, _ := json.Marshal(err)
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
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

func SendDocument(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024 << 20)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	downFile, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error retrieving data from form-data")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fileBytes, err := ioutil.ReadAll(downFile)

	resp, err := apiModels.Client.Upload(context.Background(), fileBytes, whatsmeow.MediaDocument)
	if err != nil {
		fmt.Println("Error while uploading")
		fmt.Println(err)
		w.WriteHeader(http.StatusNotExtended)
		return
	}

	documentMsg := &waProto.DocumentMessage{
		Url:           &resp.URL,
		Mimetype:      proto.String(http.DetectContentType(fileBytes)), // replace this with the actual mime type
		Caption:       proto.String(r.FormValue("caption")),
		Title:         proto.String(handler.Filename),
		FileName:      proto.String(handler.Filename),
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
		DocumentMessage: documentMsg,
	})

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotImplemented)
		b, _ := json.Marshal(err)
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
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

func SendVideo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024 << 20)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	downFile, _, err := r.FormFile("video")
	if err != nil {
		fmt.Println("Error retrieving data from form-data")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fileBytes, err := ioutil.ReadAll(downFile)

	resp, err := apiModels.Client.Upload(context.Background(), fileBytes, whatsmeow.MediaVideo)
	if err != nil {
		fmt.Println("Error while uploading")
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	videoMsg := &waProto.VideoMessage{
		Url:           &resp.URL,
		Mimetype:      proto.String(http.DetectContentType(fileBytes)), // replace this with the actual mime type
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
		VideoMessage: videoMsg,
	})

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotImplemented)
		b, _ := json.Marshal(err)
		_, err = w.Write(b)
		if err != nil {
			log.Println(err)
		}
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

func SendAnimation(w http.ResponseWriter, r *http.Request) {

}

func SendSticker(w http.ResponseWriter, r *http.Request) {

}

func SendContact(w http.ResponseWriter, r *http.Request) {

}

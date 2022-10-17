package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go_whatsapp_api/api/pkg/utils"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
	"strings"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	apiModels "go_whatsapp_api/app/pkg/models"
)

func SendMessage(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error while parsing form.")
		w.WriteHeader(http.StatusNotAcceptable)
		return err
	}

	var us = strings.Split(r.FormValue("chat_id"), "@")
	jid := types.JID{
		User:   us[0],
		Server: us[1],
	}

	msg := waProto.Message{
		Conversation: proto.String(r.FormValue("text")),
	}

	resp, err := apiModels.Client.SendMessage(context.Background(), jid, "", &msg)

	if err != nil {
		fmt.Println("Error sending message")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(resp)
		_, err = w.Write(b)
	}
	return err
}

func SendPhoto(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(1024 << 20)
	if err != nil {
		fmt.Println("Error while parsing form.")
		w.WriteHeader(http.StatusNotAcceptable)
		return err
	}

	downFile := r.MultipartForm.File["image"][0]
	file, err := downFile.Open()
	defer file.Close()
	if err != nil {
		fmt.Println("Error opening file")
		w.WriteHeader(http.StatusNotAcceptable)
		return err
	}

	fileBytes, err := ioutil.ReadAll(file)

	resp, err := apiModels.Client.Upload(context.Background(), fileBytes, whatsmeow.MediaImage)
	if err != nil {
		fmt.Println("Error while uploading")
		w.WriteHeader(http.StatusInternalServerError)
		return err
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
		fmt.Println("Error while sending message.")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(rsp)
		w.Write(b)
	}
	return err
}

func SendVideo(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(1024 << 20)
	if err != nil {
		fmt.Println("Error while parsing form.")
		w.WriteHeader(http.StatusNotAcceptable)
		return err
	}

	downFile := r.MultipartForm.File["video"][0]
	file, err := downFile.Open()
	if err != nil {
		fmt.Println("Error opening file.")
		w.WriteHeader(http.StatusNotAcceptable)
		return err
	}

	fileBytes, err := ioutil.ReadAll(file)

	resp, err := apiModels.Client.Upload(context.Background(), fileBytes, whatsmeow.MediaVideo)
	if err != nil {
		fmt.Println("Error while uploading")
		w.WriteHeader(http.StatusInternalServerError)
		return err
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
		fmt.Println("Error while sending message")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(rsp)
		w.Write(b)
	}
	return err
}

func SendDocument(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(1024 << 20)
	if err != nil {
		fmt.Println("Error while parsing form.")
		w.WriteHeader(http.StatusNotAcceptable)
		return err
	}

	downFile := r.MultipartForm.File["file"][0]
	file, err := downFile.Open()
	defer file.Close()
	if err != nil {
		fmt.Println("Error opening file.")
		w.WriteHeader(http.StatusNotAcceptable)
		return err
	}

	fileBytes, err := ioutil.ReadAll(file)

	resp, err := apiModels.Client.Upload(context.Background(), fileBytes, whatsmeow.MediaDocument)
	if err != nil {
		fmt.Println("Error while uploading")
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	documentMsg := &waProto.DocumentMessage{
		Url:           &resp.URL,
		Mimetype:      proto.String(http.DetectContentType(fileBytes)), // replace this with the actual mime type
		Caption:       proto.String(r.FormValue("caption")),
		Title:         proto.String(downFile.Filename),
		FileName:      proto.String(downFile.Filename),
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
		fmt.Println("Error while sending message.")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(rsp)
		w.Write(b)
	}
	return err
}

func SendAudio(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(1024 << 20)
	if err != nil {
		fmt.Println("Error while parsing form.")
		w.WriteHeader(http.StatusNotAcceptable)
		return err
	}

	downFile := r.MultipartForm.File["audio"][0]
	file, err := downFile.Open()
	defer file.Close()
	if err != nil {
		fmt.Println("Error opening file.")
		return err
	}

	fileBytes, err := ioutil.ReadAll(file)

	resp, err := apiModels.Client.Upload(context.Background(), fileBytes, whatsmeow.MediaAudio)
	if err != nil {
		fmt.Println("Error while uploading")
		w.WriteHeader(http.StatusNotAcceptable)
		return err
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(rsp)
		w.Write(b)
	}
	return err
}

func SendSticker(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(1024 << 20)
	if err != nil {
		fmt.Println("Error while parsing form.")
		w.WriteHeader(http.StatusNotAcceptable)
		return err
	}

	downFile := r.MultipartForm.File["sticker"][0]
	file, err := downFile.Open()
	if err != nil {
		fmt.Println("Error opening file.")
		return err
	}
	resImage, err := utils.ResizeImage(file)
	fileBytes, err := ioutil.ReadAll(resImage)

	resp, err := apiModels.Client.Upload(context.Background(), fileBytes, whatsmeow.MediaImage)
	if err != nil {
		fmt.Println("Error while uploading")
		w.WriteHeader(http.StatusNotAcceptable)
		return err
	}

	stickerMsg := &waProto.StickerMessage{
		Url:           &resp.URL,
		FileSha256:    resp.FileSHA256,
		FileEncSha256: resp.FileEncSHA256,
		MediaKey:      resp.MediaKey,
		Mimetype:      proto.String(http.DetectContentType(fileBytes)), // replace this with the actual mime type
		DirectPath:    &resp.DirectPath,
		FileLength:    &resp.FileLength,
	}

	var us = strings.Split(r.FormValue("chat_id"), "@")
	jid := types.JID{
		User:   us[0],
		Server: us[1],
	}

	rsp, err := apiModels.Client.SendMessage(context.Background(), jid, "", &waProto.Message{
		StickerMessage: stickerMsg,
	})
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(rsp)
		w.Write(b)
	}
	return err
}

func SendContact(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error while parsing form.")
		w.WriteHeader(http.StatusNotAcceptable)
		return err
	}

	var us = strings.Split(r.FormValue("chat_id"), "@")
	jid := types.JID{
		User:   us[0],
		Server: us[1],
	}

	msg := waProto.ContactMessage{
		DisplayName: proto.String(r.FormValue("first_name") + " " + r.FormValue("last_name")),
		Vcard:       proto.String(r.FormValue("number")),
	}

	resp, err := apiModels.Client.SendMessage(context.Background(), jid, "", &waProto.Message{
		ContactMessage: &msg,
	})
	if err != nil {
		fmt.Println("Error sending message")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(resp)
		w.Write(b)
	}
	return err
}

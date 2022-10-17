package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go_whatsapp_api/api/pkg/models"
	"go_whatsapp_api/api/pkg/utils"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	apiModels "go_whatsapp_api/app/pkg/models"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error while parsing form.")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	var us = strings.Split(r.FormValue("chat_id"), "@")
	jid := types.JID{
		User:   us[0],
		Server: us[1],
	}

	msg := waProto.Message{
		Conversation: proto.String(r.FormValue("text")),
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

	downFile := r.MultipartForm.File["image"][0]
	file, err := downFile.Open()
	defer file.Close()
	if err != nil {
		fmt.Println("Error opening file")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)

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

func SendVideo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024 << 20)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	downFile := r.MultipartForm.File["video"][0]
	file, err := downFile.Open()
	if err != nil {
		fmt.Println("Error opening file.")
		return
	}

	fileBytes, err := ioutil.ReadAll(file)

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

func SendDocument(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024 << 20)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	downFile := r.MultipartForm.File["file"][0]
	file, err := downFile.Open()
	defer file.Close()
	if err != nil {
		fmt.Println("Error opening file.")
		return
	}

	fileBytes, err := ioutil.ReadAll(file)

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

	downFile := r.MultipartForm.File["audio"][0]
	file, err := downFile.Open()
	defer file.Close()
	if err != nil {
		fmt.Println("Error opening file.")
		return
	}

	fileBytes, err := ioutil.ReadAll(file)

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

func SendSticker(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1024 << 20)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	downFile := r.MultipartForm.File["sticker"][0]
	file, err := downFile.Open()
	if err != nil {
		fmt.Println("Error opening file.")
		return
	}
	resImage, err := utils.ResizeImage(file)
	fileBytes, err := ioutil.ReadAll(resImage)

	resp, err := apiModels.Client.Upload(context.Background(), fileBytes, whatsmeow.MediaImage)
	if err != nil {
		fmt.Println("Error while uploading")
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
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

func SendContact(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	data := models.ContactRequest{}
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

	msg := waProto.ContactMessage{
		DisplayName: proto.String(data.FirstName + " " + data.LastName),
		Vcard:       proto.String(data.PhoneNumber),
	}

	fmt.Println(apiModels.Client.IsConnected())
	resp, err := apiModels.Client.SendMessage(context.Background(), jid, "", &waProto.Message{
		ContactMessage: &msg,
	})

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

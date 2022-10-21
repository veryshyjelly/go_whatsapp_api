package controllers

import (
	"bytes"
	"fmt"
	"go.mau.fi/whatsmeow/binary/proto"
	config2 "go_whatsapp_api/app/cmd/config"
	"go_whatsapp_api/app/pkg/models"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HandleText(mess *proto.Message, from string, dest string) error {
	data := url.Values{}
	data.Add("chat_id", dest)
	data.Add("text", from+": "+mess.GetConversation()+mess.ExtendedTextMessage.GetText())

	resp, err := http.PostForm("http://localhost"+config2.TelegramPort+"/sendMessage/", data)
	if err != nil {
		return err
	}

	log.Println(resp)
	return nil
}

func HandlePhoto(mess *proto.ImageMessage, from string, dest string) error {
	client := &http.Client{Timeout: time.Minute * 60}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := models.Client.Download(mess)
	if err != nil {
		log.Println("Error downloading file")
		return err
	}

	var caption string
	if mess.Caption != nil {
		caption = from + ": " + *mess.Caption
	} else {
		caption = from + " sent an image."
	}

	sendQuery := map[string]interface{}{
		"chat_id": dest,
		"caption": caption,
	}

	for k, v := range sendQuery {
		fw, err := writer.CreateFormField(k)
		_, err = io.Copy(fw, strings.NewReader(fmt.Sprint(v)))
		if err != nil {
			return err
		}
	}
	fw, err := writer.CreateFormFile("image", "sendThis"+strings.Split(*mess.Mimetype, "/")[1])
	_, err = io.Copy(fw, bytes.NewReader(file))
	if err != nil {
		fmt.Println("Error while copying data")
		return err
	}
	err = writer.Close()
	if err != nil {
		fmt.Println("Error while closing writer.")
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost"+config2.TelegramPort+"/sendPhoto/", body)
	if err != nil {
		fmt.Println("Error creating request.")
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Println(resp)
	return nil
}

func HandleVideo(mess *proto.VideoMessage, from string, dest string) error {
	client := &http.Client{Timeout: time.Minute * 60}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := models.Client.Download(mess)
	if err != nil {
		log.Println("Error downloading file")
		return err
	}

	var caption string
	if mess.Caption != nil {
		caption = from + ": " + *mess.Caption
	} else {
		caption = from + " sent a video."
	}

	sendQuery := map[string]interface{}{
		"chat_id": dest,
		"caption": caption,
	}

	for k, v := range sendQuery {
		fw, err := writer.CreateFormField(k)
		_, err = io.Copy(fw, strings.NewReader(fmt.Sprint(v)))
		if err != nil {
			return err
		}
	}
	fw, err := writer.CreateFormFile("video", "sendThis"+strings.Split(*mess.Mimetype, "/")[1])
	_, err = io.Copy(fw, bytes.NewReader(file))
	if err != nil {
		fmt.Println("Error while copying data")
		return err
	}
	err = writer.Close()
	if err != nil {
		fmt.Println("Error while closing writer.")
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost"+config2.TelegramPort+"/sendVideo/", body)
	if err != nil {
		fmt.Println("Error creating request.")
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Println(resp)
	return nil
}

func HandleDocument(mess *proto.DocumentMessage, from string, dest string) error {
	client := &http.Client{Timeout: time.Minute * 60}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := models.Client.Download(mess)
	if err != nil {
		log.Println("Error downloading file")
		return err
	}

	var caption string
	if mess.Caption != nil {
		caption = from + ": " + *mess.Caption
	} else {
		caption = from + " sent a document."
	}

	sendQuery := map[string]interface{}{
		"chat_id": dest,
		"caption": caption,
	}

	for k, v := range sendQuery {
		fw, err := writer.CreateFormField(k)
		_, err = io.Copy(fw, strings.NewReader(fmt.Sprint(v)))
		if err != nil {
			return err
		}
	}
	fw, err := writer.CreateFormFile("document", mess.GetTitle()+"."+strings.Split(*mess.Mimetype, "/")[1])
	_, err = io.Copy(fw, bytes.NewReader(file))
	if err != nil {
		fmt.Println("Error while copying data")
		return err
	}
	err = writer.Close()
	if err != nil {
		fmt.Println("Error while closing writer.")
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost"+config2.TelegramPort+"/sendDocument/", body)
	if err != nil {
		fmt.Println("Error creating request.")
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Println(resp)
	return nil
}

func HandleAudio(mess *proto.AudioMessage, from string, dest string) error {
	client := &http.Client{Timeout: time.Minute * 60}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := models.Client.Download(mess)
	if err != nil {
		log.Println("Error downloading file")
		return err
	}

	var caption = from + " sent an audio."

	sendQuery := map[string]interface{}{
		"chat_id": dest,
		"caption": caption,
	}

	for k, v := range sendQuery {
		fw, err := writer.CreateFormField(k)
		_, err = io.Copy(fw, strings.NewReader(fmt.Sprint(v)))
		if err != nil {
			return err
		}
	}
	fw, err := writer.CreateFormFile("audio", "audio."+strings.Split(*mess.Mimetype, "/")[1])
	_, err = io.Copy(fw, bytes.NewReader(file))
	if err != nil {
		fmt.Println("Error while copying data")
		return err
	}
	err = writer.Close()
	if err != nil {
		fmt.Println("Error while closing writer.")
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost"+config2.TelegramPort+"/sendAudio/", body)
	if err != nil {
		fmt.Println("Error creating request.")
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Println(resp)
	return nil
}

func HandleSticker(mess *proto.StickerMessage, dest string) error {
	fmt.Println("handling sticker")
	client := &http.Client{Timeout: time.Minute * 60}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := models.Client.Download(mess)
	if err != nil {
		log.Println("Error downloading file")
		return err
	}

	sendQuery := map[string]interface{}{
		"chat_id": dest,
	}

	for k, v := range sendQuery {
		fw, err := writer.CreateFormField(k)
		_, err = io.Copy(fw, strings.NewReader(fmt.Sprint(v)))
		if err != nil {
			return err
		}
	}
	fw, err := writer.CreateFormFile("sticker", "sendThis"+strings.Split(*mess.Mimetype, "/")[1])
	_, err = io.Copy(fw, bytes.NewReader(file))
	if err != nil {
		fmt.Println("Error while copying data")
		return err
	}
	err = writer.Close()
	if err != nil {
		fmt.Println("Error while closing writer.")
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost"+config2.TelegramPort+"/sendSticker/", body)
	if err != nil {
		fmt.Println("Error creating request.")
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Println(resp)
	return nil
}

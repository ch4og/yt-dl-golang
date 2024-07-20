package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/url"
	"os"

	tgbotapi "github.com/ch4og/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/wader/goutubedl"
)

func main() {
	loadenv()

	// Create telegram API
	telegram_token := os.Getenv("TELEGRAM_API_TOKEN")
	tapi_url := os.Getenv("TELEGRAM_API_URL")
	bot, err := tgbotapi.NewBotAPIWithAPIEndpoint(telegram_token, tapi_url+"/bot%s/%s")
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			message_text := update.Message.Text

			// Check if message is a valid link
			msgid, link := parse_message(message_text, update, bot)
			if link != "" {
				filePath, info, err := youtube_download(link, msgid, update, bot)
				if err != nil {
					edit_message(msgid, "ERROR\n\n"+err.Error(), update, bot)
				} else {
					send_video(filePath, info, msgid, update, bot)
				}
			}
		}
	}
}

func loadenv() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func youtube_download(link string, msgid int, update tgbotapi.Update, bot *tgbotapi.BotAPI) (filePath string, info goutubedl.Info, err error) {
	result, err := goutubedl.New(context.Background(), link, goutubedl.Options{})
	if err != nil {
		log.Println("Invalid URL")
		err = errors.New("this URL is not valid for downloading")
		return
	} else {
		edit_message(msgid, "downloading your video...", update, bot)
	}
	downloadResult, err := result.Download(context.Background(), "best")
	if err != nil {
		log.Println(err)
		err = errors.New("download failed")
		return
	}
	defer downloadResult.Close()
	info = result.Info
	filePath = "vids/" + info.ID + ".mp4"
	f, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		err = errors.New("failed to create file")
		return
	}
	defer f.Close()
	io.Copy(f, downloadResult)
	return
}

func parse_message(message_text string, update tgbotapi.Update, bot *tgbotapi.BotAPI) (msgid int, link string) {
	_, err := url.ParseRequestURI(message_text)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR\n\nthis is not a HTTP or HTTPS link")
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
		return
	} else {
		link = message_text
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "checking your link...")
		msg.ReplyToMessageID = update.Message.MessageID
		answ_message, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		msgid = answ_message.MessageID
		return
	}

}

func send_video(filePath string, info goutubedl.Info, msgid int, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		edit_message(msgid, "ERROR\n\nfile was downloaded, but failed to open it\n\ntry again", update, bot)
		return
	}

	inputFile := tgbotapi.FileReader{Name: "file.mp4", Reader: file}
	footer := "\n\nvia @yt_dl_golang_bot"

	video := tgbotapi.NewVideo(update.Message.Chat.ID, inputFile)
	video.Caption = info.Title + footer
	video.Duration = int(info.Duration)
	video.Width = int(info.Width)
	video.Height = int(info.Height)
	video.SupportsStreaming = true
	video.ReplyToMessageID = update.Message.MessageID

	_, err = bot.Send(video)
	file.Close()
	if err != nil {
		panic(err)
	} else {
		os.Remove(filePath)
		bot.Send(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, msgid))
		log.Println("Sent video to " + update.Message.From.UserName)
	}

}

func edit_message(msgid int, newmsg string, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewEditMessageText(update.Message.Chat.ID, msgid, newmsg)
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

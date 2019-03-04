package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"./lib"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go func() {
			thisUpdate := update

			log.Printf("[%s] %s", thisUpdate.Message.From.UserName, thisUpdate.Message.Text)

			switch thisUpdate.Message.Command() {
			case "today":
				now := time.Now()
				today := strconv.Itoa(int(now.Month())) + "/" + strconv.Itoa(now.Day())
				images, err := lib.RequestImageByDate(today)
				if err != nil {
					msg := tgbotapi.NewMessage(thisUpdate.Message.Chat.ID, "Sorry, I find nothing in that date, maybe you can try other date.")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(thisUpdate.Message.Chat.ID, "///"+today+" Hibiki Pictures ///")
					bot.Send(msg)
					for _, url := range images {
						msg := tgbotapi.NewPhotoUpload(thisUpdate.Message.Chat.ID, nil)
						msg.FileID = url
						msg.UseExisting = true
						bot.Send(msg)
					}
				}
			case "day":
				args := strings.Split(thisUpdate.Message.Text, " ")
				if len(args) < 2 {
					msg := tgbotapi.NewMessage(thisUpdate.Message.Chat.ID, "Please set correct date text, like: 12/08.")
					bot.Send(msg)
				} else {
					images, err := lib.RequestImageByDate(args[1])
					if err != nil {
						msg := tgbotapi.NewMessage(thisUpdate.Message.Chat.ID, "Sorry, I find nothing in that date, maybe you can try other date.")
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(thisUpdate.Message.Chat.ID, "/// "+args[1]+" Hibiki Pictures ///")
						bot.Send(msg)
						for _, url := range images {
							msg := tgbotapi.NewPhotoUpload(thisUpdate.Message.Chat.ID, nil)
							msg.FileID = url
							msg.UseExisting = true
							bot.Send(msg)
						}
					}
				}
			case "yandere":
				//yandere [page] [length]
				args := strings.Split(thisUpdate.Message.Text, " ")
				if len(args) < 3 {
					msg := tgbotapi.NewMessage(thisUpdate.Message.Chat.ID, "Please set correct argument, template: /yandere [page] [length], example: /yandere 1 3")
					bot.Send(msg)
				} else {
					images, err := lib.GetYanderePictures(args[1])
					if err != nil {
						msg := tgbotapi.NewMessage(thisUpdate.Message.Chat.ID, "Sorry, I find nothing, please try other page or length range.")
						bot.Send(msg)
					} else {
						userReqImgLength, err := strconv.Atoi(args[2])
						if err != nil {
							msg := tgbotapi.NewMessage(thisUpdate.Message.Chat.ID, "You're not enter the illegal length, please enter follow this template: /yandere [page] [length], example: /yandere 1 3")
							bot.Send(msg)

						} else {
							rangestart := randomInLength(userReqImgLength, len(images))
							images = images[rangestart : rangestart+userReqImgLength]

							msg := tgbotapi.NewMessage(thisUpdate.Message.Chat.ID, "/// Hibiki Pictures (Yande.re) ///")
							bot.Send(msg)
							for _, url := range images {
								msg := tgbotapi.NewPhotoUpload(thisUpdate.Message.Chat.ID, nil)
								msg.FileID = url
								msg.UseExisting = true
								bot.Send(msg)
							}
						}
					}
				}
			case "sankaku":
				//sankaku [page] [preview] [preview-length]
				args := strings.Split(thisUpdate.Message.Text, " ")
				if len(args) < 3 {
					msg := tgbotapi.NewMessage(thisUpdate.Message.Chat.ID, "Please set correct argument, template: /sankaku [page] [preview(0 or 1)], example: /sankaku 1 0")
					bot.Send(msg)
				} else {
					preview := true //limit preview
					if len(args) >= 2 {
						if args[2] == "0" {
							preview = false
						}
					}

					images, err := lib.GetSankakuPictures(args[1], preview)
					if err != nil {
						msg := tgbotapi.NewMessage(thisUpdate.Message.Chat.ID, "Sorry, I find nothing, please try other page or length range.")
						bot.Send(msg)
					} else {
						if preview && len(args) > 3 {
							userReqImgLength, err := strconv.Atoi(args[3])
							if err != nil {
								userReqImgLength = 1
							}
							rangestart := randomInLength(userReqImgLength, len(images))
							images = images[rangestart : rangestart+userReqImgLength]
						} else {
							images = []string{images[randInt(0, len(images))]} //limit length
						}

						msg := tgbotapi.NewMessage(thisUpdate.Message.Chat.ID, "/// Hibiki Pictures (Sankaku Complex) ///")
						bot.Send(msg)
						for _, url := range images {
							msg := tgbotapi.NewPhotoUpload(thisUpdate.Message.Chat.ID, nil)
							msg.FileID = url
							msg.UseExisting = true
							bot.Send(msg)
						}
					}
				}
			case "nhentai":
				//nhentai [cover only (0 or 1)]
				args := strings.Split(thisUpdate.Message.Text, " ")
				if len(args) < 2 {
					msg := tgbotapi.NewMessage(thisUpdate.Message.Chat.ID, "Please set correct argument, template: /nhentai [cover only (0 or 1)], example: /nhentai 0")
					bot.Send(msg)
				} else {
					book := lib.GetNhentaiBooks()
					if args[1] == "0" {
						msg := tgbotapi.NewPhotoUpload(thisUpdate.Message.Chat.ID, nil)
						msg.Caption = book.Title + "\n" + book.Link
						msg.FileID = book.Cover
						msg.UseExisting = true
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewPhotoUpload(thisUpdate.Message.Chat.ID, nil)
						msg.Caption = book.Title + "\n" + book.Link
						msg.FileID = book.Cover
						msg.UseExisting = true
						bot.Send(msg)

						for _, url := range book.Content {
							msg := tgbotapi.NewPhotoUpload(thisUpdate.Message.Chat.ID, nil)
							msg.FileID = url
							msg.UseExisting = true
							bot.Send(msg)
						}
					}
				}
			default:
				fmt.Println("no thing to do")

			}
		}()
	}
}

func randomInLength(length int, lengthmax int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s1)
	randomIlgRange := r.Int31n(int32(lengthmax))
	for ((randomIlgRange - 1) + int32(length)) > int32(lengthmax) {
		randomIlgRange = r.Int31n(int32(lengthmax))
	}
	return int(randomIlgRange)
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

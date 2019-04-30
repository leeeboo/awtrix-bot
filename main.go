package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/leeeboo/awtrix-bot/utils"
)

var apiBase *string

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("on", "on"),
		tgbotapi.NewInlineKeyboardButtonData("off", "off"),
	),
)

type Config struct {
	Token string `env:"TOKEN"`
	API   string `env:"API"`
}

var cfg Config

func init() {
	cfg = Config{}
}

func main() {

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	if cfg.Token == "" {
		panic("Token empty.")
	}

	if cfg.API == "" {
		panic("API empty.")
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.CallbackQuery != nil {

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))

			switch update.CallbackQuery.Data {
			case "on":
				resp, err := power(true)

				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, err.Error()))
				} else {
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, string(resp)))
				}
				break
			case "off":
				resp, err := power(false)

				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, err.Error()))
				} else {
					bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, string(resp)))
				}
				break
			default:
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Error"))
			}
		}

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = "/power Turn on or off screen power."
			case "power":
				msg.Text = "Choose please:"
				msg.ReplyMarkup = numericKeyboard
			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		}
	}
}

func power(status bool) ([]byte, error) {

	api := fmt.Sprintf("%s/api/v3/basics", cfg.API)
	body, err := utils.HttpPost(api, map[string]interface{}{"power": status})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(string(body))
	return body, nil
}

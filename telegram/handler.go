package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"monitoring/health"
	"monitoring/utils"
)

type Handler struct {
	bot    *tgbotapi.BotAPI
	health health.Handler
}

func NewHandler(bot *tgbotapi.BotAPI, health health.Handler) Handler {
	return Handler{
		bot:    bot,
		health: health,
	}
}

func (h *Handler) InitChannel() (err error) {
	config := tgbotapi.NewUpdate(0)
	config.Timeout = 60
	updates := h.bot.GetUpdatesChan(config)

	for channel := range updates {
		if channel.Message == nil { // ignore any non-Message updates
			continue
		}

		if !channel.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Extract the command from the Message.
		switch channel.Message.Command() {
		case "start":
			err = h.RunStart(channel.Message.Chat.ID)
		case "help":
			err = h.RunHelp(channel.Message.Chat.ID)
		case "status":
			err = h.RunStatus(channel.Message.Chat.ID)
		default:
			err = h.RunUndefined(channel.Message.Chat.ID)
		}

		if err != nil {
			return
		}
	}

	return
}

func (h *Handler) RunStart(chatID int64) (err error) {
	message := tgbotapi.NewMessage(chatID, "Hello! I understand /help and /status.")

	if _, err = h.bot.Send(message); err != nil {
		return
	}

	return
}

func (h *Handler) RunHelp(chatID int64) (err error) {
	message := tgbotapi.NewMessage(chatID, "Hi :)")

	if _, err = h.bot.Send(message); err != nil {
		return
	}

	return
}

func (h *Handler) RunStatus(chatID int64) (err error) {
	body := h.health.Check()
	if err != nil {
		return
	}
	message := tgbotapi.NewMessage(chatID, utils.PrettyString(string(body)))

	if _, err = h.bot.Send(message); err != nil {
		return
	}

	return
}

func (h *Handler) RunUndefined(chatID int64) (err error) {
	message := tgbotapi.NewMessage(chatID, "Undefined command.")

	if _, err = h.bot.Send(message); err != nil {
		return
	}

	return
}

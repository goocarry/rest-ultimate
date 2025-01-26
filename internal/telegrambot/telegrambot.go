package telegrambot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/goocarry/rest-ultimate/internal/storage"
	"log/slog"
)

type TelegramBot struct {
	bot     *tgbotapi.BotAPI
	storage storage.Storage
	logger  *slog.Logger
}

func NewTelegramBot(token string, storage storage.Storage, logger *slog.Logger) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		bot:     bot,
		storage: storage,
		logger:  logger,
	}, nil
}

func (tb *TelegramBot) Start() {
	tb.logger.With(slog.String("tg bot name", tb.bot.Self.UserName)).Info("starting telegram bot")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tb.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			tb.handleMessage(update.Message)
		}
	}
}

func (tb *TelegramBot) handleMessage(msg *tgbotapi.Message) {
	tb.logger.With(
		slog.String("tg user name", msg.From.UserName),
		slog.String("msg text", msg.Text),
	).Info("received tg message")

	switch msg.Text {
	case "/start":
		tb.handleStartCommand(msg)
	default:
		tb.handleUserMessage(msg)
	}
}

func (tb *TelegramBot) handleStartCommand(msg *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, "Добро пожаловать! Отправьте свой email для регистрации.")
	_, err := tb.bot.Send(reply)
	if err != nil {
		tb.logger.Error("failed to send message: %v", err)
	}
}

func (tb *TelegramBot) handleUserMessage(msg *tgbotapi.Message) {
	user := storage.User{
		Email: msg.Text,
	}
	_, err := tb.storage.User().RegisterUser(user)
	if err != nil {
		tb.logger.Error("failed to register user: %v", err)
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Ошибка при регистрации. Попробуйте еще раз.")
		tb.bot.Send(reply)
		return
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "Вы успешно зарегистрированы!")
	tb.bot.Send(reply)
}

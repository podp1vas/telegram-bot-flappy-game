package main

import (
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Panic("BOT_TOKEN не установлен в переменных окружения")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Авторизован как @%s", bot.Self.UserName)

	// Запускаем HTTP-сервер для отдачи веб-страницы
	// Путь к папке с веб-файлами (относительно директории запуска)
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	go func() {
		log.Println("Запуск веб-сервера на :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Ошибка веб-сервера: %s", err)
		}
	}()

	// Сброс webhook, чтобы работать через GetUpdates
	_, err = bot.Request(tgbotapi.DeleteWebhookConfig{})

	if err != nil {
		log.Printf("Ошибка сброса webhook: %s", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("Получено сообщение от %s: %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Text == "/start" {
			urlButton := tgbotapi.NewInlineKeyboardButtonURL("Играть в Flappy Bird", "http://localhost:8080/index.html")
			keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(urlButton))

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажми на кнопку, чтобы открыть игру:")
			msg.ReplyMarkup = keyboard
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка при отправке сообщения: %s", err)
			}
		}
	}
}

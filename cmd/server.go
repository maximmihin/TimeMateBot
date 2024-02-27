package main

import (
	"log"
	"os"
	"timeMate_bot/app/timeMateBot/controllers"
	"timeMate_bot/app/timeMateBot/session"
	"timeMate_bot/app/timeMateBot/tgApp"
	"timeMate_bot/storage"
)

func main() {

	db, err := storage.New(os.Getenv("DB_PATH"))
	if err != nil {
		log.Panic(err)
	}

	alLowedUserName := os.Getenv("TG_ALLOWED_USERNAME")
	if alLowedUserName == "" {
		log.Fatal("не указан TG_ALLOWED_USERNAME")
	}
	sessions := session.NewLocalMemory()

	tgTimeMate, err := tgApp.NewApp(tgApp.NewConfigBuilder().
		BotToken(os.Getenv("TELEGRAM_APITOKEN")).
		DebugMode(true).
		Build())

	updates := tgTimeMate.StartPolling()

	for update := range updates {

		if update.CallbackQuery != nil {
			// todo переписать на нормальный мидлвер
			if update.CallbackQuery.Message.Chat.UserName != alLowedUserName {
				controllers.SendPremisionDenied(tgTimeMate.Bot, update.CallbackQuery.Message.Chat.ID)
			}

			sessions.HistoryMassage.SetChatId(update.CallbackQuery.Message.Chat.ID)
			controllers.CallBackProcessing(tgTimeMate.Bot, sessions, db, update)
		} else if update.Message != nil {
			// todo переписать на нормальный мидлвер
			if update.Message.Chat.UserName != alLowedUserName {
				controllers.SendPremisionDenied(tgTimeMate.Bot, update.Message.Chat.ID)
			}

			sessions.HistoryMassage.SetChatId(update.Message.Chat.ID)
			if update.Message.IsCommand() {
				controllers.CommandProcessing(tgTimeMate.Bot, sessions, db, update)
			} else {
				controllers.AddEventProcessing(tgTimeMate.Bot, sessions, db, update)
			}
		}
	}
}

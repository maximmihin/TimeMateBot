package controllers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"sync"
	"time"
	"timeMate_bot/app/timeMateBot/session"
	"timeMate_bot/storage"
	"timeMate_bot/use_cases"
	"unicode/utf8"
)

var tagsKeyboard = keyboardConfig([]string{"суета", "работа", "цифра", "сон"})

func AddEventProcessing(bot *tgbotapi.BotAPI, session *session.LocalMemory, db *storage.Store, update tgbotapi.Update) {
	_ = session.LastMove.Reset()

	tag, date, comment := ExtractEvent(update.Message.Text)
	if date == "" {
		date = time.Now().Format("2006-01-02 15:04:05")
	}
	err := use_cases.AddNewEvent(db, session.LastMove, tag, date, comment)
	if err != nil {
		// todo написать пользователю, что у данные не валидны
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	t1 := truncateToDay(time.Now())
	t2 := t1.AddDate(0, 0, 1)

	historyToday, err := use_cases.GetHistory(db, t1, t2)
	err = session.HistoryMassage.Update(bot, historyToday)
	if err != nil {
		fmt.Println("не получилось обновить историю")
	}

	tmpDate, _ := time.Parse("2006-01-02 15:04:05", date)
	txt := fmt.Sprintf("Успешно добавлено событие: \n\n %s %s",
		tag, tmpDate.Format("01.02 15:04"))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, txt)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отменить", "undo_add"),
		),
	)

	message, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	session.WgUndo = wg
	wg.Add(1)

	go func() {
		time.Sleep(10 * time.Second)
		wg.Done()
	}()

	go func() {
		wg.Wait()
		dm := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
		resp, err := bot.Request(dm)
		if err != nil || !resp.Ok {
			panic(err)
		}

		dm = tgbotapi.NewDeleteMessage(message.Chat.ID, update.Message.MessageID)
		resp, err = bot.Request(dm)
		if err != nil || !resp.Ok {
			panic(err)
		}
		// todo сделать что-то с этим трешем
		wg.Add(1)
	}()

}

func ExtractEvent(message string) (tag, date, comment string) {

	potentialEvent := strings.Split(message, " ")

	if len(potentialEvent) == 0 {
		return
	}
	tag = potentialEvent[0]

	if len(potentialEvent) >= 2 {

		// если второй аргумент существует, но не является валидным временем, то все последующие аргументы включая текущий будут отнесены к комментарию
		cleanDate, err := time.Parse("15:04", potentialEvent[1])
		if err != nil {
			//log.Println(err)
			comment = strings.Join(potentialEvent[1:], " ")
			return
		}

		// если время всё-таки валидное - приводим его к нужному нам виду - "2006-01-02 15:04:05"
		t := time.Now()
		year, month, day := t.Date()
		// плюсуем к спаршенной date сегодняшнюю дату и секунды; конвертируем в нужном формате строки
		date = cleanDate.
			AddDate(year, int(month)-1, day-1).
			Add(time.Duration(t.Second()) * time.Second).
			Format("2006-01-02 15:04:05")
	}

	if len(potentialEvent) > 2 {
		comment = strings.Join(potentialEvent[2:], " ")
	}

	return
}

func CallBackProcessing(bot *tgbotapi.BotAPI, session *session.LocalMemory, db *storage.Store, update tgbotapi.Update) {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := bot.Request(callback); err != nil {
		panic(err)
	}

	var txt string
	switch update.CallbackQuery.Data {
	case "undo_add":
		txt = "событие удалено"
		err := session.LastMove.Undo(db)
		if err != nil {
			log.Fatal(err)
		}
		session.WgUndo.Done()

		t1 := truncateToDay(time.Now())
		t2 := t1.AddDate(0, 0, 1)

		changes, err := use_cases.GetHistory(db, t1, t2)
		if err != nil {
			log.Fatal(err)
		}
		err = session.HistoryMassage.Update(bot, changes)
		if err != nil {
			log.Fatal(err)
		}

		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, txt)
		message, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			time.Sleep(10 * time.Second)
			dm := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
			resp, err := bot.Request(dm)
			if err != nil || !resp.Ok {
				panic(err)
			}
		}()

	default:
		fmt.Println(update.CallbackQuery.Data)
		txt = "а такой колбек я не знаю..."
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, txt)
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func CommandProcessing(bot *tgbotapi.BotAPI, session *session.LocalMemory, db *storage.Store, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "help":
		msg.Text = fmt.Sprintf("Окей, есть несколько основных команд:\n" +
			"/help - вызов этой памятки\n" +
			"/add - добавить событие\n" +
			"/keyboard - отправить клавиатуру тегов\n",
		//"/stat - вызвать меню статистики\n" +
		//"/export - возвращает json файл со всеми вашими данными\n" +
		//"/import - импортирует json в ваш аккаунт\n"
		)
	case "add":
		msg.Text = "Чтобы добавить событие отправьте обычное сообщение в формате:\n" +
			"<имя тега> | добавит событие с текущим временем\n" +
			"<имя тега> hh:mm | добавит событие на указанное время\n"
	//case "stat":
	//	msg.Text = "Скоро завезу"
	case "keyboard":
		cArgs := update.Message.CommandArguments()
		if cArgs == "" {
			msg.Text = fmt.Sprintf("Клавиатура добавлена\n")
			msg.ReplyMarkup = tagsKeyboard
		} else {
			msg.Text = fmt.Sprintf("Клавиатура обновлена\n")
			msg.ReplyMarkup = keyboardConfig(strings.Split(cArgs, " "))
		}
	case "history":
		t1 := truncateToDay(time.Now())
		t2 := t1.AddDate(0, 0, 1)

		history, err := use_cases.GetHistory(db, t1, t2)
		if err != nil {
			log.Fatal(err)
		}
		err = session.HistoryMassage.Update(bot, history)
		if err != nil {
			fmt.Println(err)
		}

		return

	default:
		msg.Text = "Такого я пока не умею"
	}

	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func keyboardConfig(buttons []string) tgbotapi.ReplyKeyboardMarkup {
	kbb := make([]tgbotapi.KeyboardButton, 0, len(buttons))

	for _, btn := range buttons {
		// todo сделать нормальную обработку ошибки (как минимум уведомить пользователя)
		if !utf8.ValidString(btn) {
			continue
		}
		kbb = append(kbb, tgbotapi.NewKeyboardButton(btn))
	}

	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(kbb...),
	)
}

func SendPermissionDenied(bot *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Я тебя не знаю")
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

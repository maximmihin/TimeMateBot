package session

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"sync"
	"timeMate_bot/app/timeMateBot/last_move"
	"timeMate_bot/use_cases"
)

type HistoryMessage interface {
	Update(*tgbotapi.BotAPI, []string) error
	SetChatId(int64)
}

type NotebookHistory struct {
	messageId int
	chatId    int64
}

func (h *NotebookHistory) SetChatId(chatId int64) {
	h.chatId = chatId
}

func (h *NotebookHistory) Update(bot *tgbotapi.BotAPI, historyArr []string) error {
	if h.chatId == 0 {
		return fmt.Errorf("не указан chat id")
	}

	if h.messageId == 0 {
		//	todo создать новый пост
		historyStr := strings.Join(historyArr, "")
		msg := tgbotapi.NewMessage(h.chatId, historyStr)
		message, err := bot.Send(msg)
		if err != nil {
			fmt.Println(err)
		}
		h.messageId = message.MessageID
	} else {
		historyStr := strings.Join(historyArr, "")
		edit := NewEditMessage(h.messageId, h.chatId, historyStr)
		_, err := bot.Send(edit)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil

}

type LocalMemory struct {
	LastMove       use_cases.StackUserChanges
	HistoryMassage HistoryMessage
	WgUndo         *sync.WaitGroup
}

func NewLocalMemory() *LocalMemory {
	return &LocalMemory{
		LastMove:       &last_move.LastMove{},
		HistoryMassage: &NotebookHistory{},
	}
}

//func (m LocalMemory) GetLastMove() use_cases.StackUserChanges {
//	return m.LastMove
//}
//
//func (m LocalMemory) GetHistoryPost() HistoryMessage {
//	return m.HistoryMassage
//}
//
//func (m LocalMemory) GetWgUndo() *sync.WaitGroup {
//	return m.WgUndo
//}

func NewEditMessage(messageId int, chatId int64, text string) tgbotapi.EditMessageTextConfig {
	return tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			BaseChatMessage: tgbotapi.BaseChatMessage{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: chatId,
				},
				MessageID: messageId,
			},
		},
		Text: text,
	}
}

//
//package session
//
//import (
//"timeMate_bot/app/timeMateBot/last_move"
//"timeMate_bot/use_cases"
//)
//
//type LocalMemory struct {
//	LastMove         use_cases.StackUserChanges
//	historyMassageId *int64
//}
//
//func NewLocalMemory() *LocalMemory {
//	return &LocalMemory{
//		LastMove:         &last_move.LastMove{},
//		historyMassageId: new(int64),
//	}
//}
//
//func (m LocalMemory) GetLastMove() use_cases.StackUserChanges {
//	return m.LastMove
//}
//
//func (m LocalMemory) GetHistoryPost() *int64 {
//	return m.historyMassageId
//}

package services

import tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (t TG) sendMessage(chatId int64, messageId int, replyMsg string) {
	msg := tgbot.NewMessage(chatId, replyMsg)
	msg.ReplyToMessageID = messageId

	t.bot.Send(msg)
}

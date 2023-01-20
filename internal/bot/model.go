package bot

import "gopkg.in/telebot.v3"

const DEFAULT_ID = "_id"

type StoredMessage struct {
	MessageID string `bson:"message_id" json:"message_id"`
	ChatID    int64  `bson:"chat_id" json:"chat_id"`
}
type MediaMessage struct {
	UniqueID  string `bson:"_id" json:"_id"`
	MessageID string `bson:"message_id" json:"message_id"`
	ChatID    int64  `bson:"chat_id" json:"chat_id"`
	FileID    string `bson:"file_id", json:"file_id"`
}

func NewMediaMessage(ctx telebot.Context) *MediaMessage {
	var (
		messageID, chatID = ctx.Message().MessageSig()
		uniqueID          = ctx.Message().Photo.File.UniqueID
		fileID            = ctx.Message().Photo.File.FileID
	)
	return &MediaMessage{
		UniqueID:  uniqueID,
		MessageID: messageID,
		ChatID:    chatID,
		FileID:    fileID,
	}
}
func (x MediaMessage) MessageSig() (string, int64) {
	return x.MessageID, x.ChatID
}
func (x StoredMessage) MessageSig() (string, int64) {
	return x.MessageID, x.ChatID
}

type FeedMessage struct {
	UniqueID  string `bson:"_id" json:"_id"`
	MessageID string `bson:"message_id" json:"message_id"`
	ChatID    int64  `bson:"chat_id" json:"chat_id"`
	Categorys Category
	VideosURL map[string]string
	CreatedAt *string
	ExpireAt  *string

	Tag         *NormalTag
	OptionalTag *TagsStorage
	EventTag    *EventTag
}

func (x FeedMessage) MessageSig() (string, int64) {
	return x.MessageID, x.ChatID
}

//func NewFeedMessage(msg *telebot.Message) *FeedMessage  {
//	return FeedMessage{UniqueID: msg.Photo.UniqueID, MessageID: msg.ID}
//}

func editMessage() {
	mess := StoredMessage{}
	clients.BotClient.Edit(mess, "")
}

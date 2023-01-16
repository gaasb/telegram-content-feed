package bot

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
}

func (x FeedMessage) MessageSig() (string, int64) {
	return x.MessageID, x.ChatID
}

func editMessage() {
	mess := StoredMessage{}
	clients.BotClient.Edit(mess, "")
}

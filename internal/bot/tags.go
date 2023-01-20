package bot

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/telebot.v3"
	"regexp"
	"time"
)

const (
	NORMAL_TYPE     = "normal"
	ADDITIONAL_TYPE = "additional"
	EVENT_TYPE      = "event"
)

var (
	TAG_NORMAL     NormalTag     = "🏷"
	TAG_ADDITIONAL AdditionalTag = "🎴"
	TAG_EVENT      EventTag      = "🏆"
)
var (
	tagTypes = map[string]Tager{
		string(TAG_NORMAL):     &TAG_NORMAL,
		string(TAG_ADDITIONAL): &TAG_ADDITIONAL,
		string(TAG_EVENT):      &TAG_EVENT,
	}
	tagPattenr = regexp.MustCompile(`^[a-zA-Z]+$`)
)

var (
	addMenu = &telebot.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true, ForceReply: true}

	tagNormalBtn     = telebot.Btn{Text: fmt.Sprintf("%s %s tag", string(TAG_NORMAL), NORMAL_TYPE)}
	tagAdditionalBtn = telebot.Btn{Text: fmt.Sprintf("%s %s tag", string(TAG_ADDITIONAL), ADDITIONAL_TYPE)}
	tagEventBtn      = telebot.Btn{Text: fmt.Sprintf("%s %s tag", string(TAG_EVENT), EVENT_TYPE)}
)

type NormalTag string
type AdditionalTag string
type EventTag string

type TagsStorage struct {
	Id          primitive.ObjectID `bson:"_id" json:"_id"`
	CaptionName string             `bson:"caption_name" json:"caption_name"`
	Type        string             `bson:"type" json:"type"`
	ExpireAt    *time.Time         `bson:"expire_at,omitempty" json:"expire_at,omitempty"`
}

type Tager interface {
	GetText() (string, string)
	//Message()
}

func (t *NormalTag) GetText() (string, string)     { return string(*t), NORMAL_TYPE }
func (t *AdditionalTag) GetText() (string, string) { return string(*t), ADDITIONAL_TYPE }
func (t *EventTag) GetText() (string, string)      { return string(*t), EVENT_TYPE }

func onEmoji(emoji string) (Tager, error) {
	var result = tagTypes[emoji]
	if result != nil {
		return result, nil
	}
	return nil, errors.New("emoji not found")
}

func parseEmoji(emoji string) string {
	var runes = []rune(emoji)
	return string(runes[0:1])
}

func DoOnTagEvents(ctx telebot.Context) {
	if msg := ctx.Message(); msg.ReplyTo != nil && msg.IsReply() && len(msg.ReplyTo.Text) > 0 {
		emoji := parseEmoji(msg.ReplyTo.Text)
		if test, ok := onEmoji(emoji); ok == nil {
			if tagPattenr.MatchString(msg.Text) {
				ctx.Send("Added")
			} else {
				ctx.Send(emoji+" Not valid. Write without whitespaces, numbers and symbols. ONLY TEXT", telebot.ForceReply)
			}

			makeText(test.GetText())
		}
	}
}
func makeText(va string, ba string) string {
	return fmt.Sprintf("%s %s tag", va, ba)
}
func createUniqueIndexForCaption() {
	database.Collection("").Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.M{"caption_name": 1}, Options: options.Index().SetUnique(true)})
}

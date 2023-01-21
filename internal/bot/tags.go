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
	"strconv"
	"strings"
	"time"
)

const (
	NORMAL_TYPE     = "normal"
	ADDITIONAL_TYPE = "additional"
	EVENT_TYPE      = "event"
)
const (
	REMOVE_TAG_UNIQUE = "remove_tag"
)

var removeTagButton = telebot.Btn{Unique: REMOVE_TAG_UNIQUE}

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
	tagPattenr = regexp.MustCompile(`^[a-zA-Zа-яА-Я]+$`)
)

var (
	addMenu = &telebot.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true, ForceReply: true}

	tagNormalBtn     = telebot.Btn{Text: fmt.Sprintf("%s %s tag", string(TAG_NORMAL), NORMAL_TYPE)}
	tagAdditionalBtn = telebot.Btn{Text: fmt.Sprintf("%s %s tag", string(TAG_ADDITIONAL), ADDITIONAL_TYPE)}
	tagEventBtn      = telebot.Btn{Text: fmt.Sprintf("%s %s tag", string(TAG_EVENT), EVENT_TYPE)}

	editBtn = telebot.Btn{Text: "Edit"}
)

type NormalTag string
type AdditionalTag string
type EventTag string

type TagsStorage struct {
	Id          *primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	CaptionName string              `bson:"caption_name" json:"caption_name"`
	Type        string              `bson:"type" json:"type"`
	ExpireAt    *time.Time          `bson:"expire_at,omitempty" json:"expire_at,omitempty"`
}

type Tager interface {
	Append(tag string) error
	//Message()
}

func (t *NormalTag) Append(tag string) error {
	return InsertTag(&TagsStorage{Type: NORMAL_TYPE, CaptionName: tag})
}
func (t *AdditionalTag) Append(tag string) error {
	return InsertTag(&TagsStorage{Type: ADDITIONAL_TYPE, CaptionName: tag})
}
func (t *EventTag) Append(tag string) error {
	return InsertTag(&TagsStorage{Type: EVENT_TYPE, CaptionName: tag})
}

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
		if tag, ok := onEmoji(emoji); ok == nil && len(msg.Text) > 0 {
			text := strings.ToLower(msg.Text)
			if tagPattenr.MatchString(text) {
				if err := tag.Append(text); err != nil {
					ctx.Send(emoji+" Tag already in storage", telebot.ForceReply)
					return
				}
			} else {
				ctx.Send(emoji+" Not valid. Write without whitespaces, numbers and symbols. ONLY TEXT", telebot.ForceReply)
				return
			}
			ctx.Send("Added")
		}
	}
}
func GenButtonsForEdit(ctx telebot.Context, data string) {
	values, err := GetTagsByCaptionValue(data)
	if values == nil {
		ctx.Send(fmt.Sprintf("%s is empty", data))
	}
	newReply := ctx.Bot().NewMarkup()
	var buttons []telebot.Btn
	outputText := strings.Builder{}
	outputText.WriteString("⬇️Select tag value⬇️\n\n")
	if err != nil || len(values) <= 0 {
		return
	}
	for i, item := range values {
		buttons = append(buttons, telebot.Btn{Unique: REMOVE_TAG_UNIQUE, Text: strconv.Itoa(i), Data: item.Id.String() + "\t"})
		outputText.WriteString(fmt.Sprintf("%o:️\t%s\n", i, item.CaptionName))
	}
	newReply.Inline(newReply.Split(4, buttons)...)
	ctx.Edit(outputText.String(), newReply)
}
func RemoveTagHandler() (interface{}, telebot.HandlerFunc) {
	return &removeTagButton, func(ctx telebot.Context) error {
		data := strings.Split(ctx.Data(), "\t")
		if data == nil || len(data) < 2 {
			return errors.New("empty data in remove tag button")
		}
		switch data[0] {
		case "cancel":
			err := ctx.Delete()
			return err
		case "/r":
			if err := RemoveTagById(data[1]); err == nil {
				ctx.Edit("Successfully deleted")
				return err
			} else {
				return err
			}
		case "/u":
			//TODO <-----------------------------------------------------
			break
		default:
			reply := ctx.Bot().NewMarkup()
			reply.Inline(reply.Split(2, []telebot.Btn{
				telebot.Btn{Text: "Update", Data: "/u\t" + ctx.Data(), Unique: REMOVE_TAG_UNIQUE},
				telebot.Btn{Text: "Remove", Data: "/r\t" + ctx.Data(), Unique: REMOVE_TAG_UNIQUE},
				telebot.Btn{Text: "Cancel", Data: "cancel", Unique: REMOVE_TAG_UNIQUE}},
			)...)
			ctx.Edit("Select action", reply)
		}
		return nil
	}
}

func createUniqueIndexForCaption() {
	database.Collection("").Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.M{"caption_name": 1}, Options: options.Index().SetUnique(true)})
}

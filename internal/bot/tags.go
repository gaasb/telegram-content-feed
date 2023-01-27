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
	ACTION_BUTTON_UNIQUE = "remove_tag"
	ACTION_CANCEL        = "/c"
	ACTION_UPDATE        = "/u"
	ACTION_REMOVE        = "/r"
	ACTION_EMPTY         = ""
)

var removeTagButton = telebot.Btn{Unique: ACTION_BUTTON_UNIQUE}

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
	tagActions = map[string]func(ctx telebot.Context) error{
		ACTION_CANCEL: actionCancel,
		ACTION_UPDATE: actionUpdate,
		ACTION_REMOVE: actionRemove,
		ACTION_EMPTY:  actionEmpty,
	}
	tagPattern = regexp.MustCompile(`^[a-zA-Zа-яА-Я]+$`)
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
	GetReplyKeyboard(ctx telebot.Context) *telebot.ReplyMarkup
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

func (t *NormalTag) GetReplyKeyboard(ctx telebot.Context) *telebot.ReplyMarkup {
	return replyKeyboard(NORMAL_TYPE, ctx)
}
func (t *AdditionalTag) GetReplyKeyboard(ctx telebot.Context) *telebot.ReplyMarkup {
	return replyKeyboard(ADDITIONAL_TYPE, ctx)
}
func (t *EventTag) GetReplyKeyboard(ctx telebot.Context) *telebot.ReplyMarkup {
	return replyKeyboard(EVENT_TYPE, ctx)
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
			if tagPattern.MatchString(text) {
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
	tagsStorageList, err := GetTagsByTagType(data)
	if tagsStorageList == nil || len(tagsStorageList) <= 0 || err != nil {
		ctx.Send(fmt.Sprintf("%s is empty", data))
		return
	}
	instanceButtonsOfTags(ctx, tagsStorageList)
}
func instanceButtonsOfTags(ctx telebot.Context, list []*TagsStorage) {
	newReply, outputText := genButtons(list)
	ctx.Edit(outputText, newReply)
}
func genButtons(list []*TagsStorage) (*telebot.ReplyMarkup, string) {
	var replyMarkup = new(telebot.ReplyMarkup)
	var buttons []telebot.Btn
	outputText := strings.Builder{}
	outputText.WriteString("⬇️Select tag value⬇️\n\n")
	for i, item := range list {
		buttons = append(buttons,
			telebot.Btn{
				Unique: ACTION_BUTTON_UNIQUE,
				Text:   strconv.Itoa(i),
				Data:   fmt.Sprintf("%s\t", item.Id.Hex()),
			})
		outputText.WriteString(fmt.Sprintf("%o:️\t%s\n", i, item.CaptionName))
	}
	replyMarkup.Inline(replyMarkup.Split(4, buttons)...)
	return replyMarkup, outputText.String()
}

func writeTagCaption(data ...string) string {
	result := strings.Builder{}
	for i, item := range data {
		if i%2 != 0 && item != AcceptMedia {
			result.WriteString("#" + item + "\t")
		} else {
			continue
		}
	}
	return result.String()
}

func onAcceptBtnController(ctx telebot.Context) bool {
	data := strings.Split(ctx.Data(), "\t")
	switch {
	case len(data) == 6 && data[0] == AcceptMedia:
		completeCaption := writeTagCaption(data...)
		ctx.EditCaption(completeCaption)
		return true

	case len(data) == 4 && data[0] == ADDITIONAL_TYPE:
		var replyMarkup = new(telebot.ReplyMarkup)
		backButton := replyMarkup.Data("⬅️ Back", BackBtn, ctx.Data())
		finishBtn := replyMarkup.Data("✅ Accept", AcceptMedia, AcceptMedia+"\t\t"+ctx.Data())
		replyMarkup.Inline(replyMarkup.Row(backButton, finishBtn))
		completeCaption := writeTagCaption(data...)
		ctx.EditCaption(completeCaption)
		ctx.Bot().EditReplyMarkup(ctx.Message(), replyMarkup)
		return false
	case len(data) == 2 && data[0] == NORMAL_TYPE:
		completeCaption := writeTagCaption(data...)
		ctx.EditCaption(completeCaption)
		ctx.Bot().EditReplyMarkup(ctx.Message(), tagTypes[string(TAG_ADDITIONAL)].GetReplyKeyboard(ctx))
		fmt.Println(ctx.Data(), " NORMAL")
		return false
	default:
		return true
	}

}

func replyKeyboard(data string, ctx telebot.Context) *telebot.ReplyMarkup {
	var replyMarkup = new(telebot.ReplyMarkup)
	var buttons []telebot.Btn
	var list []*TagsStorage
	list, _ = GetTagsByTagType(data)
	if len(list) > 0 || list != nil {
		emptyButton := selector.Data("", BackBtn)
		buttons = append(buttons, emptyButton, refreshBtn, dismissBtn)
		for _, item := range list {
			if ctx != nil {
				//emptyButton.Data = ctx.Data()
				buttons = append(buttons, telebot.Btn{Text: item.CaptionName, Unique: AcceptMedia, Data: fmt.Sprintf("%s\t%s\t%s", data, item.CaptionName, ctx.Data())})
			} else {
				buttons = append(buttons, telebot.Btn{Text: item.CaptionName, Unique: AcceptMedia, Data: fmt.Sprintf("%s\t%s", data, item.CaptionName)})
			}
		}
	} else {
		acceptButton := selector.Data("✅ Accept", AcceptMedia)
		buttons = append(buttons, acceptButton, refreshBtn, dismissBtn)
	}
	replyMarkup.Inline(replyMarkup.Split(3, buttons)...)
	return replyMarkup

}
func replyKeyboardWithContext(data string, ctx telebot.Context) *telebot.ReplyMarkup {
	var replyMarkup = new(telebot.ReplyMarkup)
	var buttons []telebot.Btn
	var list []*TagsStorage
	list, _ = GetTagsByTagType(data)
	buttons = append(buttons, acceptBtn, refreshBtn, dismissBtn)
	if list != nil {
		for _, item := range list {
			buttons = append(buttons, telebot.Btn{Text: item.CaptionName, Unique: AcceptMedia, Data: fmt.Sprintf("%s\t%s\t%s", ctx.Data(), data, item.CaptionName)})
		}
	}
	replyMarkup.Inline(replyMarkup.Split(3, buttons)...)
	return replyMarkup

}

func actionCancel(ctx telebot.Context) error { err := ctx.Delete(); return err }
func actionUpdate(ctx telebot.Context) error { return nil }
func actionRemove(ctx telebot.Context) error { return nil }
func actionEmpty(ctx telebot.Context) error  { return nil }

func RemoveTagHandler() (interface{}, telebot.HandlerFunc) {
	return &removeTagButton, func(ctx telebot.Context) error {
		data := strings.Split(ctx.Data(), "\t")
		if data == nil || len(data) < 2 {
			return errors.New("empty data in remove tag button")
		}
		switch data[0] {
		case ACTION_CANCEL:
			err := ctx.Delete()
			return err
		case ACTION_REMOVE:
			if err := RemoveTagById(data[1]); err == nil {
				ctx.Edit("Successfully deleted")
				return err
			} else {
				ctx.Send("Cant delete")
				return err
			}
		case ACTION_UPDATE:
			ctx.Send("Send edited value", telebot.ForceReply)
			//TODO <-----------------------------------------------------
			break
		default:
			reply := ctx.Bot().NewMarkup()
			reply.Inline(reply.Split(2, []telebot.Btn{
				{Text: "Update", Data: splitData(ACTION_UPDATE) + ctx.Data(), Unique: ACTION_BUTTON_UNIQUE},
				{Text: "Remove", Data: splitData(ACTION_REMOVE) + ctx.Data(), Unique: ACTION_BUTTON_UNIQUE},
				{Text: "Cancel", Data: splitData(ACTION_CANCEL), Unique: ACTION_BUTTON_UNIQUE}},
			)...)
			ctx.Edit("Select action", reply)
		}
		return nil
	}
}
func splitData(action string) string {
	return action + "\t"
}
func createUniqueIndexForCaption() {
	database.Collection("").Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.M{"caption_name": 1}, Options: options.Index().SetUnique(true)})
}

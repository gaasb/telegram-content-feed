package bot

import (
	"fmt"
	"gopkg.in/telebot.v3"
)

type Command string
type Query string

var creator *telebot.Chat
var messages []*StoredMessage //MESSAGES FOR REVIEW

const (
	BanCommand         Command = "/ban"
	CreateEventCommand Command = "/create_event"
	Dice               Command = "/dice"
	ReviewMediaContent         = "/start_review"
)
const (
	UpdateTagQuery Query = "update_tag"
	AcceptMedia          = "accept_media|"
	DismissMedia         = "dismiss_media|"
)

var cmd map[Command]func() (interface{}, telebot.HandlerFunc)

var (
	selector = &telebot.ReplyMarkup{}

	refreshBtn = selector.Data("🔃 Refresh", "refresh", "refresh_media")
	acceptBtn  = selector.Data("✔️ Accept", "accept", AcceptMedia)
	dismissBtn = selector.Data("❌ Dismiss", "dismiss", DismissMedia)
)

func OnDice() (interface{}, telebot.HandlerFunc) {
	var f = func(ctx telebot.Context) error {
		var err error

		if len(messages) > 0 {
			_, err = ctx.Bot().Copy(ctx.Sender(), messages[0])
		}
		if err != nil {
			ctx.Send("err")
			return err
		}
		var Cube = &telebot.Dice{Type: "🎲"}
		Cube.Send(ctx.Bot(), ctx.Recipient(), nil)
		return nil
	}
	return string(Dice), f
}

func OnReviewMediaContent() (interface{}, telebot.HandlerFunc) {
	selector.Inline(selector.Row(acceptBtn, refreshBtn, dismissBtn))
	var fun = func(ctx telebot.Context) error {
		if ok := FindAllMedia(); ok != nil && len(ok) > 0 {
			var er error
			for _, i := range ok {
				//acceptBtn.InlineQuery += i.UniqueID
				//acceptBtn.InlineQuery += i.UniqueID
				if _, err := ctx.Bot().Copy(ctx.Sender(), i, selector); err != nil {
					_ = RemoveMediaByID(i.UniqueID) //TODO ERROR HANDLE
					er = err
				}
			}
			return er
		}
		return nil
	}
	return ReviewMediaContent, fun
}

func OnMedia() (interface{}, telebot.HandlerFunc) {
	var fun = func(ctx telebot.Context) error {
		var (
			chatID, messageID = ctx.Message().MessageSig()
			uniqueID          = ctx.Message().Photo.UniqueID
			fileID            = ctx.Message().Photo.File.FileID
			msg               = &MediaMessage{uniqueID, chatID, messageID, fileID}
		)
		if err := AddMedia(msg); err != nil {
			ctx.Send("⛔: File is already in database!")
			return err
		}
		ctx.Send("✅: Ok!")
		return nil
	}
	return MEDIA_TYPE_BY_DEFAULT, fun
}
func OnAcceptMediaButton() (interface{}, telebot.HandlerFunc) {
	var fun = func(ctx telebot.Context) error {
		var (
			chatID, messageID = ctx.Message().MessageSig()
			uniqueID          = ctx.Message().Photo.UniqueID
			fileID            = ctx.Message().Photo.File.FileID
			msg               = &MediaMessage{uniqueID, chatID, messageID, fileID}
		)
		if dbErr := AddMediaToFeed(msg); dbErr != nil {
			emsg := ctx.Message().Photo
			//emsg.File.FileID = ""
			emsg.Caption = "TEST"
			emsg.Height = 0
			emsg.Width = 0
			fmt.Println(ctx.Edit(emsg, selector), emsg.File)
			return dbErr
		}
		_, err := ctx.Bot().Copy(channel, ctx.Message()) //TODO SAVE MESSAGE TO DB FEED COLLECTION
		if err != nil {
			return err
		}
		RemoveMediaByID(ctx.Message().Photo.UniqueID)
		ctx.Delete()
		return nil
	}
	return &acceptBtn, fun
}
func OnDismissMediaButton() (interface{}, telebot.HandlerFunc) {
	var fun = func(ctx telebot.Context) error {
		RemoveMediaByID(ctx.Message().Photo.UniqueID)
		ctx.Delete()
		return nil
	}
	return &dismissBtn, fun
}

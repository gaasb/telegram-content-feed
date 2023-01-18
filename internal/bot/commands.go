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
	AddTag                     = "/add_tag"
)
const (
	UpdateTagQuery Query = "update_tag"
	AcceptMedia          = "accept_media|"
	DismissMedia         = "dismiss_media|"
)

var cmd map[Command]func() (interface{}, telebot.HandlerFunc)

var (
	selector = &telebot.ReplyMarkup{}

	acceptBtn  = selector.Data("✅ Accept", "accept", AcceptMedia)
	refreshBtn = selector.Data("🔃 Refresh", "refresh", "refresh_media")
	dismissBtn = selector.Data("❌ Dismiss", "dismiss", DismissMedia)
)

// TODO -> IF LENGTH < []TAGS send Accept Btn esle Tags Btns				<--------

func init() {
	selector.Inline(selector.Row(acceptBtn, refreshBtn, dismissBtn))
	//selector.Split(3, selector.Row(acceptBtn, refreshBtn, dismissBtn))
	//selector.Split(8, selector.Row(refreshBtn))
}

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
	var fun = func(ctx telebot.Context) error {
		var er error
		if ok := FindAllMedia(); ok != nil && len(ok) > 0 {
			for _, i := range ok {
				dismissBtn.Data = i.UniqueID //TODO <-------------------------
				if _, err := ctx.Bot().Copy(ctx.Sender(), i, selector); err != nil {
					_ = RemoveMediaByID(i.UniqueID) //TODO ERROR HANDLE
					er = err
				} else {
					return nil
				}

			}
		}
		ctx.Send("No Media in Database")
		return er
	}
	return ReviewMediaContent, fun
}

func OnMedia() (interface{}, telebot.HandlerFunc) {
	var fun = func(ctx telebot.Context) error {
		var (
			chatID, messageID = ctx.Message().MessageSig()
			uniqueID          = ctx.Message().Photo.File.UniqueID
			fileID            = ctx.Message().Photo.File.FileID
			msg               = &MediaMessage{uniqueID, chatID, messageID, fileID}
		)
		replyErr := ctx.Reply("✅: Ok!")
		if replyErr != nil {
			ctx.Send("🔴: Not Found, file have been removed!")
			return replyErr
		}
		if err := AddMedia(msg); err != nil {
			ctx.Send("⛔: File is already in database!")

			return err
		}
		return nil
	}
	return MEDIA_TYPE_BY_DEFAULT, fun
}
func OnAcceptMediaButton() (interface{}, telebot.HandlerFunc) {
	var fun = func(ctx telebot.Context) error {
		var (
			chatID, messageID = ctx.Message().MessageSig()
			uniqueID          = ctx.Message().Photo.File.UniqueID
			fileID            = ctx.Message().Photo.File.FileID
			msg               = &MediaMessage{uniqueID, chatID, messageID, fileID}
		)
		if dbErr := FindMediaById(msg.UniqueID); dbErr != nil {
			updateInvalidMediaPost(ctx)
			fmt.Println("Already in Feed")
			return nil
		}
		_, err := ctx.Bot().Copy(channel, msg)
		if err != nil {
			return err
		}
		RemoveMediaByID(msg.UniqueID)
		AddMediaToFeed(msg) //TODO CHANGE type TO FeedMessage
		updateInvalidMediaPost(ctx)

		//ctx.Delete()
		return nil
	}
	return &acceptBtn, fun
}
func OnDismissMediaButton() (interface{}, telebot.HandlerFunc) {
	var fun = func(ctx telebot.Context) error {
		//RemoveMediaByID(ctx.Message().Photo.File.UniqueID)
		//ctx.Delete()
		updateInvalidMediaPost(ctx)
		return nil
	}
	return &dismissBtn, fun
}
func OnAddTagCommand() (interface{}, telebot.HandlerFunc) {
	return AddTag, func(context telebot.Context) error {

		return nil
	}
}
func OnRefreshButton() (interface{}, telebot.HandlerFunc) {
	return &refreshBtn, func(ctx telebot.Context) error {
		_, err := ctx.Bot().EditMedia(ctx.Message(), &telebot.Photo{
			File: telebot.File{FileID: "AgACAgIAAxkBAAICRmPGxVNHd5bgb0Q_M5kxw5x07tvnAAKjxDEbw74wSkJVI1oGX0l2AQADAgADeQADLQQ",
				UniqueID: "AQADo8QxG8O-MEp-"},
			Width:   0,
			Height:  0,
			Caption: "test",
		})
		fmt.Println(err)
		return nil
	}
}

func updateInvalidMediaPost(ctx telebot.Context) {
	RemoveMediaByID(ctx.Message().Photo.File.UniqueID)
	refreshMedia := FindFirstMedia()

	if refreshMedia == nil {
		ctx.Delete()
		ctx.Send("🤷 Видимо на этом всё...🚩\nПопробуй обновить!\t" + ReviewMediaContent)
		return
	}
	// <- TODO get first media from review
	for refreshMedia != nil {

		_, err := ctx.Bot().EditMedia(ctx.Message(), &telebot.Photo{
			Caption: "PRESSED",
			File: telebot.File{
				FileID:   refreshMedia.FileID,
				UniqueID: refreshMedia.UniqueID,
			}}, selector)
		if err != nil {
			RemoveMediaByID(refreshMedia.UniqueID)
		} else {
			return
		}
		refreshMedia = FindFirstMedia()
	}
	ctx.Delete()

}

package bot

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"strconv"
)

type Command string
type Query string

var creator *telebot.Chat

const (
	StartCommand               = "/start"
	BanCommand         Command = "/ban"
	CreateEventCommand Command = "/create_event"
	Dice               Command = "/dice"
	ReviewMediaContent         = "/start_review"
	AddTag                     = "/add_tag"
	AddTagKeyboard             = "⚙ AddTag"
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
		var Cube = &telebot.Dice{Type: "🎲"}
		Cube.Send(ctx.Bot(), ctx.Recipient(), nil)
		return nil
	}
	return string(Dice), f
}

func OnReviewMediaContent() (interface{}, telebot.HandlerFunc) {
	return ReviewMediaContent, func(ctx telebot.Context) error {
		var er error
		if ok := FindAllMedia(); ok != nil && len(ok) > 0 {
			for _, i := range ok {
				dismissBtn.Data = i.UniqueID //TODO <-------------------------
				//dismissBtn.Inline().
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
}

func OnMedia() (interface{}, telebot.HandlerFunc) {
	return MEDIA_TYPE_BY_DEFAULT, func(ctx telebot.Context) error {
		var (
			chatID, messageID = ctx.Message().MessageSig()
			uniqueID          = ctx.Message().Photo.File.UniqueID
			fileID            = ctx.Message().Photo.File.FileID
			msg               = &MediaMessage{uniqueID, chatID, messageID, fileID}
		)
		replyErr := ctx.Reply("✅: Ok!")
		if replyErr != nil {
			ctx.Send("🔴: Not Found, probably file have been removed!")
			return replyErr
		}
		if err := AddMedia(msg); err != nil {
			ctx.Send("⛔: File is already in database!")

			return err
		}
		return nil
	}

}
func OnAcceptMediaButton() (interface{}, telebot.HandlerFunc) {
	return &acceptBtn, func(ctx telebot.Context) error {
		mediaMessage := NewMediaMessage(ctx)
		if dbErr := FindMediaById(mediaMessage.UniqueID); dbErr != nil {
			updateInvalidMediaPost(ctx)
			fmt.Println("Already in Feed")
			return nil
		}
		channelMsg, err := ctx.Bot().Copy(channel, mediaMessage)
		if err != nil {
			return err
		}

		RemoveMediaByID(mediaMessage.UniqueID)
		mediaMessage.ChatID = channel.ID
		mediaMessage.MessageID = strconv.Itoa(channelMsg.ID)
		AddMediaToFeed(mediaMessage) //TODO CHANGE type TO FeedMessage
		updateInvalidMediaPost(ctx)
		return nil
	}
}
func OnDismissMediaButton() (interface{}, telebot.HandlerFunc) {
	return &dismissBtn, func(ctx telebot.Context) error {
		//RemoveMediaByID(ctx.Message().Photo.File.UniqueID)
		//ctx.Delete()
		updateInvalidMediaPost(ctx)
		return nil
	}

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
		// ctx.Bot().EditReplyMarkup(msg) TODO <=====================================================================
		fmt.Println(err)
		return err
	}
}
func OnAddTag() (interface{}, telebot.HandlerFunc) {
	return AddTag, func(ctx telebot.Context) error {
		//menu := ctx.Bot().NewMarkup()
		//menu.Inline(menu.Row(btn))
		addMenu.ForceReply = true
		addMenu.OneTimeKeyboard = true
		addMenu.Reply(addMenu.Row(tagNormalBtn, tagAdditionalBtn, tagEventBtn))
		ctx.Send("➕\tSelect what you want to add", addMenu)
		return nil
	}
}

func TagNormalButton() (interface{}, telebot.HandlerFunc) {
	return &tagNormalBtn, func(ctx telebot.Context) error {
		ctx.Send("🏷 Send text without whitespaces", telebot.ForceReply)
		return nil
	}
}

func TagAdditionalButton() (interface{}, telebot.HandlerFunc) {
	return &tagAdditionalBtn, func(ctx telebot.Context) error {
		ctx.Send("🎴 Send text without whitespaces", telebot.ForceReply)
		return nil
	}
}

func TagEventButton() (interface{}, telebot.HandlerFunc) {
	return &tagEventBtn, func(ctx telebot.Context) error {
		ctx.Send("Send text without whitespaces", telebot.ForceReply)
		return nil
	}
}

func HandleText() (interface{}, telebot.HandlerFunc) {
	return telebot.OnText, func(ctx telebot.Context) error {
		DoOnTagEvents(ctx)
		return nil
	}
}
func updateInvalidMediaPost(ctx telebot.Context) {
	RemoveMediaByID(ctx.Message().Photo.File.UniqueID)
	refreshMedia := FindFirstMedia()

	// <- TODO get first media from review
	for refreshMedia != nil {

		_, err := ctx.Bot().EditMedia(ctx.Message(), &telebot.Photo{
			Caption: "от @asd\n#PRESSED\t#week",
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
	ctx.Send("🤷 Видимо на этом всё...🚩\nПопробуй обновить!\t" + ReviewMediaContent)
}

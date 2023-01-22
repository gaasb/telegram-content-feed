package bot

import (
	"github.com/gaasb/telegram-content-feed/pkg/utils"
	"gopkg.in/telebot.v3"
	"log"
)

// For any media
//
// MEDIA_TYPE_BY_DEFAULT = telebot.OnMedia
const (
	CHANNEL_ID            = "CHANNEL_ID"
	MEDIA_TYPE_BY_DEFAULT = telebot.OnPhoto //telebot.OnMedia
)

var (
	channel *telebot.Chat
	tags    []string
)

func Setup() {
	build()
	setChannel()
	instanceDatabaseCollections()
	_ = clients.BotClient.SetCommands("/start start")
	creator, _ = clients.BotClient.ChatByUsername("@xgaax")
	//fmt.Println(clients.BotClient.AdminsOf(creator))

	clients.BotClient.Handle(OnDice())
	clients.BotClient.Handle(OnMedia())
	clients.BotClient.Handle(OnAddTag())
	clients.BotClient.Handle(HandleText())
	clients.BotClient.Handle(OnReviewMediaContent())
	clients.BotClient.Handle(OnAcceptMediaButton())
	clients.BotClient.Handle(OnDismissMediaButton())
	clients.BotClient.Handle(OnRefreshButton())

	clients.BotClient.Handle(OnEditAction())
	clients.BotClient.Handle(RemoveTagHandler())
	clients.BotClient.Handle(OnEditButton())
	clients.BotClient.Handle(TagNormalButton())
	clients.BotClient.Handle(TagAdditionalButton())
	clients.BotClient.Handle(TagEventButton())
	clients.BotClient.Start()
}

func setChannel() {
	channelID := utils.GetEnv(CHANNEL_ID)
	databaseName := utils.GetEnv(DATABASE_NAME)
	db := clients.DatabaseClient.Database(databaseName)
	chat, err := clients.BotClient.ChatByUsername(channelID)
	if err != nil {
		log.Fatalln("Channel not found")
		return
	}
	channel = chat
	database = db
}
func OnStart() (interface{}, telebot.HandlerFunc) {
	return StartCommand, func(ctx telebot.Context) error {
		keyboad := &telebot.ReplyMarkup{ReplyKeyboard: [][]telebot.ReplyButton{}}
		ctx.Bot().EditReplyMarkup(ctx.Message(), keyboad)
		return nil
	}
}

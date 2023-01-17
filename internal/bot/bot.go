package bot

import (
	"context"
	"fmt"
	"github.com/gaasb/telegram-content-feed/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/telebot.v3"
	"log"
)

// For any media
//
// MEDIA_TYPE_BY_DEFAULT = telebot.OnMedia
const (
	CHANNEL_ID            = "CHANNEL_ID"
	MEDIA_TYPE_BY_DEFAULT = telebot.OnPhoto //telebot.OnMedia <- For any media
)

var (
	channel *telebot.Chat
	tags    []string
)

func Setup() {
	build()
	setChannel()
	instanceDatabaseCollections()
	creator, _ = clients.BotClient.ChatByUsername("xgaax")

	fmt.Println(clients.DatabaseClient.Ping(context.TODO(), readpref.Primary()))
	clients.BotClient.Handle("/start", func(c telebot.Context) error {
		return c.Send("Hello")
	})

	clients.BotClient.Handle(OnDice())
	clients.BotClient.Handle(OnMedia())
	clients.BotClient.Handle(OnReviewMediaContent())
	clients.BotClient.Handle(OnAcceptMediaButton())
	clients.BotClient.Handle(OnDismissMediaButton())
	clients.BotClient.Handle(OnRefreshButton())

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

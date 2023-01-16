package bot

import (
	"context"
	"github.com/gaasb/telegram-content-feed/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/telebot.v3"
	"log"
	"time"
)

var (
	clients *utils.Utilities[mongo.Client, telebot.Bot]
)

type MongoImpl struct {
	utils.Client[mongo.Client]
}
type TelegramImpl struct {
	utils.Client[telebot.Bot]
}

func (t *MongoImpl) Set(envValue string) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(envValue))
	if err != nil {
		log.Fatalln(err)
		return
	}
	t.Value = client
}
func (t *TelegramImpl) Set(envValue string) {
	client, err := telebot.NewBot(telebot.Settings{Token: envValue, Poller: &telebot.LongPoller{Timeout: time.Second * 5}})
	if err != nil {
		log.Fatalln(err)
		return
	}
	t.Value = client
}

func build() {
	a := &MongoImpl{}
	b := &TelegramImpl{}
	a.Client.Setter = utils.Setter(a)
	b.Client.Setter = utils.Setter(b)
	clients = utils.NewUtilities(&a.Client, &b.Client)
}

func CloseConnections() {
	clients.DatabaseClient.Disconnect(context.TODO())
	clients.BotClient.Close()
}

package main

import "github.com/gaasb/telegram-content-feed/internal/bot"

func main() {
	bot.Setup()
	defer func() {
		bot.CloseConnections()
	}()
}

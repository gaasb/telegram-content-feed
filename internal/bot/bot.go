package bot

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Setup() {
	build()
	fmt.Println(clients.DatabaseClient.Ping(context.TODO(), readpref.Primary()))
}

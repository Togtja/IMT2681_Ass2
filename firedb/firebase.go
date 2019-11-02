package firedb

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

//Ctx the connect of the firebase client
var Ctx context.Context

//Client the firestore client
var Client *firestore.Client

func init() {
	Ctx = context.Background()
	opt := option.WithCredentialsFile("firedb/servicekey.json")
	//fmt.Println(opt)
	app, err := firebase.NewApp(Ctx, nil, opt)
	if err != nil {
		fmt.Println("error initializing app:", err)
		return
	}
	//fmt.Println(app)
	Client, err = app.Firestore(Ctx)
	if err != nil {
		fmt.Println("error initializing app:", err)
		return
	}
}

//InitDataBase does
func InitDataBase() {
	fmt.Println("Database Init finished")
}

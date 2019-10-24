package firedb

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func init() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("/home/tomas/Documents/Programming/2019NTNU_Sem5/GoTime/IMT2681_Ass2/firedb/servicekey.json")
	//fmt.Println(opt)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
		fmt.Println(err)
		return
	}
	//fmt.Println(app)
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(client)

	defer client.Close()
}
func Test() {
	fmt.Println("Test")
}

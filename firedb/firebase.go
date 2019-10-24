package firedb
import (
	"fmt"
	"google.golang.org/api/option"
	"firebase.google.com/go"  
	"context"

  )
  


func init() {
	opt := option.WithCredentialsFile("servicekey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
	}
	fmt.Println(app)

}
func Test(){
	fmt.Println("Test")
}
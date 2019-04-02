package tibber

import (
	"context"

	"github.com/machinebox/graphql"
	log "github.com/sirupsen/logrus"
	"github.com/tskaard/tibber-golang/model"
)

// SendPushNotification from tibber app
func (t *Client) SendPushNotification(title, msg string) (int, error) {
	req := graphql.NewRequest(`
		mutation sendPushNotification($input: PushNotificationInput!){
			sendPushNotification(input: $input){
		  		successful
		  		pushedToNumberOfDevices
			}
	  }`)
	input := model.PushInput{
		Title:        title,
		Message:      msg,
		ScreenToOpen: "CONSUMPTION",
	}
	req.Var("input", input)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+t.Token)
	ctx := context.Background()
	//ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	var result model.PushResponse
	if err := t.gqlClient.Run(ctx, req, &result); err != nil {
		log.Error(err)
		return 0, err
	}
	return result.SendPushNotification.PushedToNumberOfDevices, nil
}

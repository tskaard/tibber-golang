package tibber

import (
	"context"

	"github.com/machinebox/graphql"
	log "github.com/sirupsen/logrus"
)

// PushResponse respons from notification api
type PushResponse struct {
	SendPushNotification SendPushNotification `json:"sendPushNotification"`
}

// SendPushNotification data in push response
type SendPushNotification struct {
	Successful              bool `json:"successful"`
	PushedToNumberOfDevices int  `json:"pushedToNumberOfDevices"`
}

// PushInput push message
type PushInput struct {
	Title        string `json:"title"`
	Message      string `json:"message"`
	ScreenToOpen string `json:"screenToOpen"`
}

// SendPushNotification from tibber app
func (t *Client) SendPushNotification(title, msg string) (int, error) {
	req := graphql.NewRequest(`
		mutation sendPushNotification($input: PushNotificationInput!){
			sendPushNotification(input: $input){
		  		successful
		  		pushedToNumberOfDevices
			}
	  }`)
	input := PushInput{
		Title:        title,
		Message:      msg,
		ScreenToOpen: "CONSUMPTION",
	}
	req.Var("input", input)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+t.Token)
	ctx := context.Background()
	//ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	var result PushResponse
	if err := t.gqlClient.Run(ctx, req, &result); err != nil {
		log.Error(err)
		return 0, err
	}
	return result.SendPushNotification.PushedToNumberOfDevices, nil
}

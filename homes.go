package tibber

import (
	"context"
	"time"

	"github.com/machinebox/graphql"
	log "github.com/sirupsen/logrus"
)

type PushResponse struct {
	SendPushNotification SendPushNotification `json:"sendPushNotification"`
}

type SendPushNotification struct {
	Successful              bool `json:"successful"`
	PushedToNumberOfDevices int  `json:"pushedToNumberOfDevices"`
}

type PushInput struct {
	Title        string `json:"title"`
	Message      string `json:"message"`
	ScreenToOpen string `json:"screenToOpen"`
}

type HomesResponse struct {
	Viewer HomeViewer `json:"viewer"`
}

type HomeViewer struct {
	Homes []Home `json:"homes"`
}

type Home struct {
	ID                string            `json:"id"`
	AppNickname       string            `json:"appNickname"`
	MeteringPointData MeteringPointData `json:"meteringPointData"`
	Features          Features          `json:"features"`
}

type MeteringPointData struct {
	ConsumptionEan string `json:"consumptionEan"`
}

type Features struct {
	RealTimeConsumptionEnabled bool `json:"realTimeConsumptionEnabled"`
}

func (t *TibberClient) GetHomes() ([]Home, error) {
	req := graphql.NewRequest(`
		query {
			viewer {
				homes {
					id
					appNickname
      				meteringPointData{
        				consumptionEan
      				}
					features {
						realTimeConsumptionEnabled
					}
					
				}
			}
		}`)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+t.Key)
	ctx := context.Background()
	var result HomesResponse
	if err := t.gqlClient.Run(ctx, req, &result); err != nil {
		log.Error(err)
		return nil, err
	}
	return result.Viewer.Homes, nil
}

func (t *TibberClient) SendPushNotification(title, msg string) (int, error) {
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
	req.Header.Set("Authorization", "Bearer "+t.Key)
	//ctx := context.Background()
	//
	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	var result PushResponse
	if err := t.gqlClient.Run(ctx, req, &result); err != nil {
		log.Error(err)
		return 0, err
	}
	return result.SendPushNotification.PushedToNumberOfDevices, nil
}

// CurrentSub CurrentSubscription `json:"currentSubscription"`

// type CurrentSubscription struct {
// 	PriceInfo PriceInfo `json:"priceInfo"`
// }

// type PriceInfo struct {
// 	Current struct {
// 		Total    float64 `json:"total"`
// 		Energy   float64 `json:"energy"`
// 		Tax      float64 `json:"tax"`
// 		Currency string  `json:"currency"`
// 	} `json:"current"`
// }

// req := graphql.NewRequest(`
// 		query {
// 			viewer {
// 				homes {
// 					id
// 					features {
// 						realTimeConsumptionEnabled
// 					}
// 					currentSubscription{
// 						priceInfo{
// 							current{
// 								total
// 								energy
// 								tax
// 								currency
// 						  	}
// 						}
// 					}
// 				}
// 			}
// 		}`)

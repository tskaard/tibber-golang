package tibber

import (
	"context"

	"github.com/machinebox/graphql"
	log "github.com/sirupsen/logrus"
)

// HomesResponse response from homes
type HomesResponse struct {
	Viewer HomeViewer `json:"viewer"`
}

// HomeViewer list of homes
type HomeViewer struct {
	Homes []Home `json:"homes"`
}

// Home structure
type Home struct {
	ID                string            `json:"id"`
	AppNickname       string            `json:"appNickname"`
	MeteringPointData MeteringPointData `json:"meteringPointData"`
	Features          Features          `json:"features"`
}

// MeteringPointData - meter number
type MeteringPointData struct {
	ConsumptionEan string `json:"consumptionEan"`
}

// Features - tibber pulse connected
type Features struct {
	RealTimeConsumptionEnabled bool `json:"realTimeConsumptionEnabled"`
}

// GetHomes get a list of homes with information
func (t *Client) GetHomes() ([]Home, error) {
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
	req.Header.Set("Authorization", "Bearer "+t.Token)
	ctx := context.Background()
	var result HomesResponse
	if err := t.gqlClient.Run(ctx, req, &result); err != nil {
		log.Error(err)
		return nil, err
	}
	return result.Viewer.Homes, nil
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

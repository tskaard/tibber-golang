package tibber

import (
	"context"

	"github.com/machinebox/graphql"
	log "github.com/sirupsen/logrus"
	"github.com/tskaard/tibber-golang/model"
)

// GetHomes get a list of homes with information
func (t *Client) GetHomes() ([]model.Home, error) {
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
	var result model.HomesResponse
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

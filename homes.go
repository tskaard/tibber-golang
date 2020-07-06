package tibber

import (
	"context"
	"fmt"
	"time"

	"github.com/machinebox/graphql"
	log "github.com/sirupsen/logrus"
)

// HomesResponse response from homes
type HomesResponse struct {
	Viewer HomesViewer `json:"viewer"`
}

// HomeViewer list of homes
type HomesViewer struct {
	Homes []Home `json:"homes"`
}

type HomeResponse struct {
	Viewer HomeViewer `json:"viewer"`
}

type HomeViewer struct {
	Home Home `json:"home"`
}

type PreviousMeterData struct {
	Power           float64 `json:"power"`
	PowerProduction float64 `json:"powerProduction"`
}

// Home structure
type Home struct {
	ID                   string              `json:"id"`
	AppNickname          string              `json:"appNickname"`
	MeteringPointData    MeteringPointData   `json:"meteringPointData"`
	Features             Features            `json:"features"`
	Address              Address             `json:"address"`
	Size                 int                 `json:"size"`
	MainFuseSize         int                 `json:"mainFuseSize"`
	NumberOfResidents    int                 `json:"numberOfResidents"`
	PrimaryHeatingSource string              `json:"primaryHeatingSource"`
	HasVentilationSystem bool                `json:"hasVentilationSystem"`
	CurrentSubscription  CurrentSubscription `json:"currentSubscription"`
	PreviousMeterData    PreviousMeterData   `json:"previousMeterData"`
}

type Address struct {
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	Address3   string `json:"address3"`
	PostalCode string `json:"postalCode"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
}

// MeteringPointData - meter number
type MeteringPointData struct {
	ConsumptionEan string `json:"consumptionEan"`
}

// Features - tibber pulse connected
type Features struct {
	RealTimeConsumptionEnabled bool `json:"realTimeConsumptionEnabled"`
}

type CurrentSubscription struct {
	PriceInfo PriceInfo `json:"priceInfo"`
}

type PriceInfo struct {
	CurrentPriceInfo CurrentPriceInfo `json:"current"`
}

type CurrentPriceInfo struct {
	Level    string    `json:"level"`
	Total    float64   `json:"total"`
	Energy   float64   `json:"energy"`
	Tax      float64   `json:"tax"`
	Currency string    `json:"currency"`
	StartsAt time.Time `json:"startsAt"`
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
					address {
						address1
						address2
						address3
						postalCode
						city
						country
						latitude
						longitude
					}
					size
					mainFuseSize
					numberOfResidents
					primaryHeatingSource
					hasVentilationSystem
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

// GetHomeById get a home with information
func (t *Client) GetHomeById(homeId string) (Home, error) {
	req := graphql.NewRequest(fmt.Sprintf(`
		query {
			viewer {
				home(id:"%s") {
					id
					appNickname
      				meteringPointData{
        				consumptionEan
      				}
					features {
						realTimeConsumptionEnabled
					}
					address {
						address1
						address2
						address3
						postalCode
						city
						country
						latitude
						longitude
					}
					size
					mainFuseSize
					numberOfResidents
					primaryHeatingSource
					hasVentilationSystem
				}
			}
		}`, homeId))
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+t.Token)
	ctx := context.Background()
	var result HomeResponse
	if err := t.gqlClient.Run(ctx, req, &result); err != nil {
		log.Error(err)
		return Home{}, err
	}
	return result.Viewer.Home, nil
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

func (t *Client) GetCurrentPrice(homeId string) (CurrentPriceInfo, error) {
	req := graphql.NewRequest(fmt.Sprintf(`
		query {
			viewer {
				home(id:"%s")  {
					currentSubscription {
						priceInfo {
							current {
								level
								total
								energy
								tax
								currency
								startsAt
							}
						}
					}
				}
			}
		}`, homeId))
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+t.Token)
	ctx := context.Background()
	var result HomeResponse
	if err := t.gqlClient.Run(ctx, req, &result); err != nil {
		log.Error(err)
		return CurrentPriceInfo{}, err
	}

	return result.Viewer.Home.CurrentSubscription.PriceInfo.CurrentPriceInfo, nil
}

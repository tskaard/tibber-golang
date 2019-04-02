package model

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

// MsgChan for reciving messages
type MsgChan chan *TibberMsg

// TibberMsg for streams
type TibberMsg struct {
	HomeID  string  `json:"homeId"`
	Type    string  `json:"type"`
	ID      int     `json:"id"`
	Payload Payload `json:"payload"`
}

// Payload in TibberMsg
type Payload struct {
	Data Data `json:"data"`
}

// Data in Payload
type Data struct {
	LiveMeasurement LiveMeasurement `json:"liveMeasurement"`
}

// LiveMeasurement in data payload
type LiveMeasurement struct {
	Timestamp              string  `json:"timestamp"`
	Power                  int     `json:"power"`
	AccumulatedConsumption float64 `json:"accumulatedConsumption"`
	AccumulatedCost        float64 `json:"accumulatedCost"`
	Currency               string  `json:"currency"`
	MinPower               int     `json:"minPower"`
	AveragePower           float64 `json:"averagePower"`
	MaxPower               int     `json:"maxPower"`
}

package tibber

type TibberMsg struct {
	Type    string  `json:"type"`
	ID      int     `json:"id"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	Data Data `json:"data"`
}

type Data struct {
	LiveMeasurement LiveMeasurement `json:"liveMeasurement"`
}

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

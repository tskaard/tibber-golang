package tibber

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const subscriptionEndpoint = "v1-beta/gql/subscriptions"
const tibberHost = "api.tibber.com"

// MsgChan for reciving messages
type MsgChan chan *StreamMsg

// StreamMsg for streams
type StreamMsg struct {
	HomeID  string  `json:"homeId"`
	Type    string  `json:"type"`
	ID      int     `json:"id"`
	Payload Payload `json:"payload"`
}

// Payload in StreamMsg
type Payload struct {
	Data Data `json:"data"`
}

// Data in Payload
type Data struct {
	LiveMeasurement LiveMeasurement `json:"liveMeasurement"`
}

// LiveMeasurement in data payload
type LiveMeasurement struct {
	Timestamp              time.Time `json:"timestamp"`
	Power                  float64   `json:"power"`
	LastMeterConsumption   float64   `json:"lastMeterConsumption"`
	LastMeterProduction    float64   `json:"lastMeterProduction"`
	AccumulatedConsumption float64   `json:"accumulatedConsumption"`
	AccumulatedCost        float64   `json:"accumulatedCost"`
	AccumulatedProduction  float64   `json:"accumulatedProduction"`
	AccumulatedReward      float64   `json:"accumulatedReward"`
	MinPower               float64   `json:"minPower"`
	AveragePower           float64   `json:"averagePower"`
	MaxPower               float64   `json:"maxPower"`
	PowerProduction        float64   `json:"powerProduction"`
	MinPowerProduction     float64   `json:"minPowerProduction"`
	MaxPowerProduction     float64   `json:"maxPowerProduction"`
	VoltagePhase1          float64   `json:"voltagePhase1"`
	VoltagePhase2          float64   `json:"voltagePhase2"`
	VoltagePhase3          float64   `json:"voltagePhase3"`
	CurrentPhase1          float64   `json:"currentPhase1"`
	CurrentPhase2          float64   `json:"currentPhase2"`
	CurrentPhase3          float64   `json:"currentPhase3"`
}

// IsExtended returns whether the report is normal or extended.
// In an extended report we would have at least one phase information
func (m *LiveMeasurement) IsExtended() bool {
	return m.CurrentPhase1 > 0 || m.CurrentPhase2 > 0 || m.CurrentPhase3 > 0
}

// AsFloatMap returns the LiveMeasurement struct as a float map
func (m *LiveMeasurement) AsFloatMap() map[string]float64 {
	return map[string]float64{
		"p_import":      m.Power,
		"e_import":      m.LastMeterConsumption,
		"e_export":      m.LastMeterProduction,
		"last_e_import": m.AccumulatedConsumption,
		"last_e_export": m.AccumulatedProduction,
		"p_import_min":  m.MinPower,
		"p_import_avg":  m.AveragePower,
		"p_import_max":  m.MaxPower,
		"p_export":      m.PowerProduction,
		"p_export_min":  m.MinPowerProduction,
		"p_export_max":  m.MaxPowerProduction,
		"u1":            m.VoltagePhase1,
		"u2":            m.VoltagePhase2,
		"u3":            m.VoltagePhase3,
		"i1":            m.CurrentPhase1,
		"i2":            m.CurrentPhase2,
		"i3":            m.CurrentPhase3,
	}
}

// Stream for subscribing to Tibber pulse
type Stream struct {
	Token       string
	ID          string
	isRunning   bool
	initialized bool
	client      *websocket.Conn
}

// NewStream with id and token
func NewStream(id, token string) *Stream {
	ts := Stream{
		ID:          id,
		Token:       token,
		isRunning:   true,
		initialized: false,
	}
	return &ts
}

// StartSubscription init connection and subscibes to home id
func (ts *Stream) StartSubscription(outputChan MsgChan) chan error {
	// Connect
	errChan := make(chan error)
	for {
		err := ts.connect()
		if err != nil {
			log.WithError(err).Error("<TibberStream> Could not connect to websocket")
			errChan <- err
			time.Sleep(time.Second * 7) // trying to repair the connection
		} else {
			ts.initialized = false
			break // connection was made
		}
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("<TibberStream> Process CRASHED with error: ", r)
				ts.StartSubscription(outputChan)
			}
		}()
		defer ts.client.Close()
		for {
			if !ts.initialized {
				ts.sendInitMsg()
			}
			tm := StreamMsg{}
			err := ts.client.ReadJSON(&tm)
			if err != nil {
				if websocket.IsCloseError(err,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
					websocket.CloseNormalClosure) {
					log.WithError(err).Error("<TibberStream> CloseError, Reconnecting after 7 seconds")
					errChan <- err
					time.Sleep(time.Second * 7) // trying to repair the connection
					ts.initialized = false
					err = ts.connect()
					if err != nil {
						log.WithError(err).Error("<TibberStream> Could not connect to websocket")
						errChan <- err
					}
					continue
				} else {
					log.WithError(err).Error()
					errChan <- err
					continue
				}
			} else {
				switch tm.Type {
				case "init_success":
					log.Info("<TibberStream> Init success")
					ts.initialized = true
					ts.sendSubMsg()

				case "subscription_success":
					log.Info("<TibberStream> Subscription success")

				case "subscription_data":
					tm.HomeID = ts.ID
					outputChan <- &tm

				case "subscription_fail":
					err := fmt.Errorf("Subscription failed")
					log.WithError(err).Error()
					errChan <- err

				default:
					err := fmt.Errorf("Unexpected message: %s", tm.Type)
					log.WithError(err).Error()
					errChan <- err
				}
			}
			if !ts.isRunning {
				log.Debug("<TibberStream> Stopping")
				break
			}
		}
	}()
	return errChan
}

func (ts *Stream) connect() error {
	defer func() {
		if r := recover(); r != nil {
			log.Error("<TibberStream> ID: ", ts.ID, " - Process CRASHED with error : ", r)
		}
	}()
	u := url.URL{Scheme: "wss", Host: tibberHost, Path: subscriptionEndpoint}
	log.Infof("<TibberStream> Connecting to %s", u.String())
	var err error
	for {
		reqHeader := make(http.Header)
		reqHeader.Add("Sec-WebSocket-Protocol", "graphql-subscriptions")
		ts.client, _, err = websocket.DefaultDialer.Dial(u.String(), reqHeader)

		if err != nil {
			log.Error("<TibberStream> Dial error", err)
			time.Sleep(time.Second * 2)
		} else {
			log.Info("<TibberStream> WS Client is connected - ID: ", ts.ID, " error: ", err)
			ts.isRunning = true
			return nil
		}
	}
}

// Stop stops stream
func (ts *Stream) Stop() {
	log.Debug("<TibberWsClient> setting isRunning to false")
	ts.isRunning = false
}

func (ts *Stream) sendInitMsg() {
	init := `{"type":"init","payload":"token=` + ts.Token + `"}`
	ts.client.WriteMessage(websocket.TextMessage, []byte(init))
}

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	// Trim the beginning and trailing " character
	return string(b[1 : len(b)-1])
}

func (ts *Stream) sendSubMsg() {
	homeID := ts.ID
	var subscriptionQuery = fmt.Sprintf(`
	subscription {
		liveMeasurement(homeId:"%s") {
			timestamp
			power
			lastMeterConsumption
			lastMeterProduction
			accumulatedConsumption
			accumulatedCost
			accumulatedProduction
			accumulatedReward
			minPower
			averagePower
			maxPower
			powerProduction
			minPowerProduction
			maxPowerProduction
			voltagePhase1
			voltagePhase2
			voltagePhase3
			currentPhase1
			currentPhase2
			currentPhase3
		}
	}`,
		homeID)

	sub := fmt.Sprintf(`
	{
		"query": "%s",
		"variables":null,
		"type":"subscription_start",
		"id":0
	}`, jsonEscape(subscriptionQuery))

	log.Debug("Subscribe with query", sub)
	ts.client.WriteMessage(websocket.TextMessage, []byte(sub))
}

# tibber-golang  
Limited implementation of the Tibber API in golang  
[developer.tibber.com](https://developer.tibber.com)  

**Possibilities:**
* Get a list of homes with id, name, meterID and features
* Send push notification from the Tibber App
* Subscribe to data from Tibber Pulse


## Usage

```go
package main

import (
	"fmt"
	"strconv"

	tibber "github.com/tskaard/tibber-golang"
)
const token = "<Tibber token>"

type Handler struct {
	tibber     *tibber.Client
	streams    map[string]*tibber.Stream
	msgChannal tibber.MsgChan
}

func NewHandler() *Handler {
	h := &Handler{}
	h.tibber = tibber.NewClient("")
	h.streams = make(map[string]*tibber.Stream)
	h.msgChannal = make(tibber.MsgChan)
	return h
}

func main() {
    h := NewHandler()
	h.tibber.Token = token
	homes, err := h.tibber.GetHomes()
	if err != nil {
		panic("Can not get homes from Tibber")
	}
	for _, home := range homes {
		fmt.Println(home.ID)
		if home.Features.RealTimeConsumptionEnabled {
			stream := tibber.NewStream(home.ID, h.tibber.Token)
			stream.StartSubscription(h.msgChannal)
			h.streams[home.ID] = stream
		}
	}
	_, err = h.tibber.SendPushNotification("Tibber-Golang", "Message from GO")
	if err != nil {
		panic("Push failed")
	}
	go func(msgChan tibber.MsgChan) {
		for {
			select {
			case msg := <-msgChan:
				h.handleStreams(msg)
			}
		}
	}(h.msgChannal)

	for {
	}

}

func (h *Handler) handleStreams(newMsg *tibber.StreamMsg) {
	fmt.Println(newMsg.Payload.Data.LiveMeasurement.Timestamp + " :: " + strconv.Itoa(newMsg.Payload.Data.LiveMeasurement.Power) + " Watt")
}
```
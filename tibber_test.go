package tibber

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

func TestGetHomes(t *testing.T) {
	token := string(helperLoadBytes(t, "token.txt"))
	tc := NewClient(token)
	homes, _ := tc.GetHomes()
	if homes == nil {
		t.Fatalf("GetHomes: %v", homes)
	}
}

func TestGetHomeById(t *testing.T) {
	token := string(helperLoadBytes(t, "token.txt"))
	tc := NewClient(token)
	homeID := string(helperLoadBytes(t, "homeId.txt"))
	home, _ := tc.GetHomeById(homeID)
	if home.ID == "" {
		t.Fatalf("GetHomeById: %s %v", homeID, home)
	}
}

func TestPush(t *testing.T) {
	token := string(helperLoadBytes(t, "token.txt"))
	tc := NewClient(token)
	_, err := tc.SendPushNotification("Golang Test", "Running golang test")
	if err != nil {
		t.Fatalf("Push: %v", err)
	}
}

func TestStreams(t *testing.T) {
	var msgCh MsgChan
	token := string(helperLoadBytes(t, "token.txt"))
	homeID := string(helperLoadBytes(t, "homeId.txt"))
	stream := NewStream(homeID, token)
	err := stream.StartSubscription(msgCh)
	if err != nil {
		t.Fatalf("Push: %v", err)
	}
	select {
	case msg := <-msgCh:
		t.Log(msg)
	case <-time.After(time.Second * 7):
		break
	}
	stream.Stop()
}

func TestGetCurrentPrice(t *testing.T) {
	token := string(helperLoadBytes(t, "token.txt"))
	tc := NewClient(token)
	homeID := string(helperLoadBytes(t, "homeId.txt"))
	priceInfo, _ := tc.GetCurrentPrice(homeID)
	if priceInfo.Level == "" {
		t.Fatalf("GetCurrentPrice: %v", priceInfo)
	}
}

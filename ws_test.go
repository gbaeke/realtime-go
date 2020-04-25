package main

import (
	"net/http"
	"testing"

	"github.com/posener/wstest"
)

type join struct {
	channel string
}

func TestWSMessage(t *testing.T) {
	ws := getWSServer("channel")
	d := wstest.NewDialer(ws)
	c, resp, err := d.Dial("ws://"+"127.0.0.1"+"/?EIO=3&transport=websocket", nil)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := resp.StatusCode, http.StatusSwitchingProtocols; got != want {
		t.Errorf("resp.StatusCode = %q, want %q", got, want)
	}

	msg := join{channel: "device01"}
	err = c.WriteJSON(&msg)
	if err != nil {
		t.Fatal(err)
	}

}

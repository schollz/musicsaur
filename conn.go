// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	syncPeriod = 500 * time.Millisecond

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

type Ntp struct {
	Origin string `json:"origin"`
	T0     int64  `json:"t0"`
	T1     int64  `json:"t1"`
	T2     int64  `json:"t2"`
	T3     int64  `json:"t3"`
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		h.broadcast <- message
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	syncer := time.NewTicker(syncPeriod)
	defer func() {
		ticker.Stop()
		syncer.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			tRecieved := time.Now().UnixNano() / 1000000
			if !ok {
				log.Println(ok)
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			var ntpTransfer Ntp
			err := json.Unmarshal(message, &ntpTransfer)
			if err != nil {
				panic(err)
			}
			if ntpTransfer.Origin == "server" {
				ntpTransfer.T3 = tRecieved
				fmt.Println(ntpTransfer.T3 - ntpTransfer.T0)
			} else if ntpTransfer.Origin == "client" {
				ntpTransfer.T1 = tRecieved
				ntpTransfer.T2 = time.Now().UnixNano() / 1000000
				ntpJson, _ := json.Marshal(ntpTransfer)
				c.write(websocket.TextMessage, ntpJson)
			}
		case <-ticker.C:
			err := c.write(websocket.PingMessage, []byte{})
			if err != nil {
				// log.Println(err)
				// return
			}
		case <-syncer.C:
			ntpJson, _ := json.Marshal(Ntp{Origin: "server", T0: time.Now().UnixNano() / 1000000})
			err := c.write(websocket.TextMessage, ntpJson)
			if err != nil {
				// log.Println(err)
				// return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writePump()
	c.readPump()
}

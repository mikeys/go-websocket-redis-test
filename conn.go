package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/garyburd/redigo/redis"
)

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan string
}

type response struct {
	// The connection to whom we send the response.
	conn *connection

	// The actual value being sent in the response.
	value string
}

func (conn *connection) redisFetchAsync(key string) {
	redisConn := redisPool.Get()
	defer redisConn.Close()
	reply, err := redis.String(redisConn.Do("GET", key))
	
	if err != nil {
		reply = ""
	}

	// 3. Send response back to client (async) through the hub.
	h.respond <- &response{conn: conn, value: reply}
}

func (conn *connection) reader() {
	for {
		var key string

		// 1. Get request from the client
		err := websocket.Message.Receive(conn.ws, &key)
		if err != nil {
			break
		}

		// 2. Fetch value from redis (async)
		go conn.redisFetchAsync(key)
	}

	conn.ws.Close()
}

func (conn *connection) writer() {
	for message := range conn.send {
		err := websocket.Message.Send(conn.ws, message)
		if err != nil {
			break
		}
	}

	conn.ws.Close()
}

func wsHandler(ws *websocket.Conn) {
	conn := &connection{send: make(chan string, 256), ws: ws}
	h.register <- conn
	defer func() { h.unregister <- conn }()
	go conn.writer()
	conn.reader()
}
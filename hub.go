package main

type hub struct {
	// Register response to the connections.
	respond chan *response

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection

	// Registered connections.
	connections map[*connection]bool
}

var h = hub{
	respond: make(chan *response),
	register: make(chan *connection),
	unregister: make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
			case conn := <-h.register:
				h.connections[conn] = true
			case conn := <-h.unregister:
				delete(h.connections, conn)
				close(conn.send)
			case response := <-h.respond:
				conn := response.conn
				select {
					case conn.send <- response.value:
					default:
						delete(h.connections, conn)
						close(conn.send)
						go conn.ws.Close()
				}
			}
	}
}

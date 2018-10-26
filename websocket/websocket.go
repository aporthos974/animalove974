package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"net/url"
)

const WS_URL = "http://localhost:8080"

var connections = make(map[*websocket.Conn]bool)

var wsHeaders = http.Header{
	"Origin":                   {WS_URL},
	"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
}

func sendAll(msg []byte) {
	for conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			delete(connections, conn)
			conn.Close()
		}
	}
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Panicln(err)
	}
	connections[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(connections, conn)
			conn.Close()
			return
		}
		sendAll(msg)
	}
}

func SendNotification(writer http.ResponseWriter, request *http.Request, object interface{}) {
	wsUrl, err := url.Parse(WS_URL + "/socket")

	wsConn, resp, err := websocket.NewClient(getRawUrl(wsUrl), wsUrl, wsHeaders, 1024, 1024)
	if err != nil {
		log.Panicf("websocket.NewClient Error: %s\nResp:%+v", err, resp)
	}
	defer wsConn.Close()

	wsConn.WriteJSON(object)
}

func getRawUrl(wsUrl *url.URL) net.Conn {
	rawConn, err := net.Dial("tcp", wsUrl.Host)
	if err != nil {
		log.Panicln(err)
	}
	return rawConn
}

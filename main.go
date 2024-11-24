package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func newServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

// Handle WebSocket connections
func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New incoming connection from client:", ws.RemoteAddr())

	// Register the connection
	s.conns[ws] = true
	defer func() {
		delete(s.conns, ws) // Clean up on disconnect
		ws.Close()
		fmt.Println("Connection closed:", ws.RemoteAddr())
	}()

	// Start reading messages
	s.readLoop(ws)
}

// Read messages from a WebSocket connection
func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Read error:", err)
			return
		}
		msg := buf[:n]

		s.broadcast(msg)
		// fmt.Println("Received message:", string(msg))

		// // Echo response back to client
		// _, writeErr := ws.Write([]byte("Thank you for the message"))
		// if writeErr != nil {
		// 	fmt.Println("Write error:", writeErr)
		// 	return
		// }
	}
}
//broad cast 
func (s *Server) broadcast (b []byte){
	for ws:= range s.conns{
		go func(ws *websocket.Conn){
			if _,err:= ws.Write(b); err!=nil{
				fmt.Println("write error:",err)
			}
		}(ws)
	}
}

func main() {
	server := newServer()

	// Define WebSocket endpoint
	http.Handle("/ws", websocket.Handler(server.handleWS))

	// Start the HTTP server
	port := ":3000"
	fmt.Println("WebSocket server listening on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

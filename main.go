package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	conns    map[*websocket.Conn]bool
	upgrader websocket.Upgrader
	mu sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *Server) handleWS(c echo.Context) error {
	ws, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	fmt.Println("new client: ", ws.RemoteAddr().String())

	s.addClient(ws)

	go s.readLoop(ws)

	return nil
}

func (s *Server) readLoop(ws *websocket.Conn) {
	defer func() {
        fmt.Println("client disconnected: ", ws.RemoteAddr().String())
		s.removeClient(ws)
    }()
	
	for {
		// Read
		var req Request
		var resp Response
		err := ws.ReadJSON(&req)
		if err != nil {
			fmt.Println("read error: ", err)
			break
		}
		fmt.Printf("%v\n", req)

		resp.DestinationAmount = req.SourceAmount*455

		err = ws.WriteJSON(resp)
		if err != nil {
			fmt.Println("write error", err)
			break
		}
	}
}

func (s *Server) addClient(ws *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conns[ws] = true
}

func (s *Server) removeClient(ws *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if Client exists, then delete it
	if _, ok := s.conns[ws]; ok {
		// close connection
		ws.Close()
		// remove
		delete(s.conns, ws)
	}
}

func main() {
	server := NewServer()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default to port 8080 if PORT is not set
	}

	// Create a new Echo instance
	e := echo.New()

	// Use middleware to log all requests and recover from panics
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define a route for the root path
	e.GET("/", func(c echo.Context) error {
		message := c.QueryParam("message")
		return c.String(http.StatusOK, fmt.Sprintf("Echo: %s", message))
	})
	e.GET("/ws", server.handleWS)

	// Start the server and listen on the specified port
	e.Logger.Fatal(e.Start(":" + port))
}

type Request struct {
	SourceAmount float64 `json:"source_amount"`
}

type Response struct {
	DestinationAmount float64 `json:"destination_amount"`
}
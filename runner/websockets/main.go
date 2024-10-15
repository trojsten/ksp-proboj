package websockets

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"iter"
	"log"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{}

var connections map[string]*websocket.Conn

func StartWebSocketServer() error {
	connections = make(map[string]*websocket.Conn)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		fmt.Println("Client connected, acquiring name")
		c.WriteMessage(websocket.TextMessage, []byte("Send name"))
		_, name, err := c.ReadMessage()
		fmt.Println("Got name: " + string(name))
		if err != nil {
			log.Print("Error acquiring name:", err)
		}
		connections[string(name)] = c
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "example.html")
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	return nil
}

func SendMessage(player string, payload string) error {
	if connections[player] == nil {
		return errors.New("connection for player " + player + " not found")
	}
	return connections[player].WriteMessage(websocket.TextMessage, []byte(payload))
}

func waitForConnection(player string) {
	fmt.Println("Waiting for player " + player + " to connect")
	defer fmt.Println("Player " + player + " has connected")
	for {
		time.Sleep(300 * time.Millisecond)
		fmt.Println("waiting")
		if _, ok := connections[player]; ok {
			fmt.Println("Return")
			return
		}
	}
}

func WaitForPlayers(players iter.Seq[string]) {
	var wg sync.WaitGroup
	//ctx, cancel := context.WithCancel(context.Background())
	for player := range players {
		wg.Add(1)

		go func(player string) {
			defer wg.Done()
			waitForConnection(player)
		}(player)

		fmt.Println("Waiting for player " + player)
	}

	wg.Wait()
	return
}

func Shutdown() {
	for _, conn := range connections {
		conn.Close()
	}
}

func ReceiveMessage(player string) chan string {
	ret := make(chan string)
	connections[player].WriteMessage(websocket.TextMessage, []byte("Next move"))
	for {
		_, message, err := connections[player].ReadMessage()
		fmt.Println("Got message: " + string(message))
		if err != nil {
			return nil
		}
		ret <- string(message)
		if string(message) == "." {
			return ret
		}
	}
}

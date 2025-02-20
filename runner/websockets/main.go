package websockets

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"iter"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{}

var connections map[string]*websocket.Conn

func StartWebSocketServer(port int, sourceRoot string) error {
	connections = make(map[string]*websocket.Conn)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		fmt.Println("Client connected, acquiring name")
		c.WriteMessage(websocket.TextMessage, []byte("GET NAME"))
		c.SetReadDeadline(time.Now().Add(100 * time.Hour))
		name, err := getMessageFromConnection(c)
		fmt.Println("Got name: " + string(name))
		if err != nil {
			log.Print("Error acquiring name:", err)
		}
		connections[name] = c
		c.SetCloseHandler(func(code int, text string) error {
			fmt.Println("Player " + name + " has disconnected")
			connections[name] = nil
			return nil
		})
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir(sourceRoot)).ServeHTTP(w, r)
	})

	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	return nil
}

func getMessageFromConnection(conn *websocket.Conn) (string, error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	var msg string
	done := make(chan bool)
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			msg += string(message)
			lines := strings.Split(msg, "\n")
			if len(lines) == 0 {
				continue
			}

			if lines[len(lines)-1] == "." {
				msg = strings.Join(lines[:len(lines)-1], "\n")
				done <- true
				return
			}

			if err != nil {
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			return msg, nil
		case <-sigChan:
			return "", errors.New("interrupted")
		default:
			time.Sleep(30 * time.Millisecond)
		}
	}
}

func SendMessage(player string, payload string) error {
	fmt.Println("Sending data to player " + player)
	defer fmt.Println("Player " + player + " received sent data successfully")
	if connections[player] == nil {
		return errors.New("connection for player " + player + " not found")
	}

	err := connections[player].WriteMessage(websocket.TextMessage, []byte("GAMEDATA\n"+payload))
	if err != nil {
		return err
	}
	for {
		res, err := getMessageFromConnection(connections[player])
		lines := strings.Split(res, "\n")
		if err != nil {
			return err
		}

		if lines[0] == "ERROR" {
			return errors.New(lines[1])
		}

		if lines[0] != "GAMEDATA LOADED" {
			continue
		}

		break
	}

	return nil
}

func waitForConnection(player string) bool {
	fmt.Println("Waiting for player " + player + " to connect")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	done := make(chan bool)
	go func() {
		for {
			time.Sleep(300 * time.Millisecond)
			if _, ok := connections[player]; ok {
				done <- true
			}
		}
	}()

	for {
		select {
		case <-done:
			return true
		case <-sigChan:
			return false
		default:
			time.Sleep(300 * time.Millisecond)
		}
	}
	return false
}

func WaitForPlayers(players iter.Seq[string]) {
	var wg sync.WaitGroup
	for player := range players {
		wg.Add(1)

		go func(player string) {
			defer wg.Done()
			if waitForConnection(player) {
				fmt.Println("Player " + player + " has connected")
			}
		}(player)
	}

	wg.Wait()
	return
}

func Shutdown() {
	for _, conn := range connections {
		conn.Close()
	}
}

func ReceiveMessage(player string) (string, error) {
	defer fmt.Println("Received from player " + player)
	fmt.Println("Receiving input from player " + player)
	if connections[player] == nil {
		return "", errors.New("connection for player " + player + " not found")
	}

	connections[player].WriteMessage(websocket.TextMessage, []byte("NEXT TURN\n"))
	for {
		res, err := getMessageFromConnection(connections[player])

		lines := strings.Split(res, "\n")

		if err != nil {
			return "", err
		}

		if lines[0] == "TURN" {
			return res, nil
		}

		fmt.Println("Unexpected message: " + res)
	}
}

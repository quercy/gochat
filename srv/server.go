// @todo implement logger
// @todo robusticise users
// make stateless, to handle multiple connections?3

package gochat

import(
	"net"
	"fmt"
	"io/ioutil"
	"os"
	"bufio"
	"strings"
)

type user struct {
	username string
	userid string
	online bool
	connection net.Conn
}

// Adds a new user to the specified user collection - returns that user
func  (UC *userCollection)  AddUser (name string, address string, cn net.Conn) {
	// @todo log: adding a user
	u := user{username:name, userid:address, online:true, connection:cn}
	UC._users = append(UC._users, &u)
}

type userCollection struct {
	_users []*user
}

func RunServer() {
	// create logger
	InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	// listen
	ln, errors := net.Listen("tcp", ":8080")
	// create the user collection
	uc := userCollection{_users:make([]*user, 0)}

	if errors != nil {
		panic(errors) // ???
	}
//	ch := make(chan []byte, 5)
//	go broadcast(ch, uc)
	for {
		connection, error := ln.Accept()
		if errors != nil {
			fmt.Println(error)
		} else {
			// start chatting
			go handleConnection (connection, &uc)
		}
	}
}

func handleConnection (connection net.Conn, uc *userCollection) {
	messages := 0
	connection.Write([]byte("Enter your username: "))
	readUsername := false
	username := ""
	for {
		input, is_command := readInput(connection)
		if(readUsername == false) { // only in first iteration
			uc.AddUser(input, connection.RemoteAddr().String(), connection)
			username = input
			welcomeMessage := fmt.Sprintf("%v has entered the chat.", username)
			uc.broadcast(welcomeMessage, "")
//			ch <- []byte(welcomeMessage)
			readUsername = true
		} else if(!is_command) {
			text := fmt.Sprintf("%v: %s", username, input)
			uc.broadcast(text, username)
//			ch <- []byte(text) // send to the broadcast stream
		} else {
//			 user has entered a command
			if(input == "/quit") {
				connection.Write([]byte("Goodbye.\n"))
				defer connection.Close()
				break
			}
		}

		messages++
		if(messages == 5) {
			connection.Write([]byte("You have exceeded your message quota! Goodbye.\n"))
			defer connection.Close()
			break
		}
	}
}

func (uc *userCollection) broadcast (msg string, usr string) {
	for i := 0; i < len(uc._users); i++ {
		if uc._users[i].username != usr {
			uc._users[i].connection.Write([]byte(msg + "\n")) // getting there
		}
	}
}

func readInput(cn net.Conn) (string, bool) {
	buf := bufio.NewReader(cn)
	output, _ := buf.ReadString('\n')
	return strings.TrimSpace(output), false // hacky way to remove newline
}
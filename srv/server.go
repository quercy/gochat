// @todo implement logger
// @todo robusticise users
// make stateless, to handle multiple connections?3

package gochat

import(
	"net"
	"fmt"
	"strconv"
	"bytes"
	"strings"
	"io/ioutil"
	"os"
)

type user struct {
	username string
	userid string
	online bool
}

// Adds a new user to the specified user collection - returns that user
func  (UC userCollection)  AddUser (name string, address string) user {
	// @todo log: adding a user
	// @todo use a list instead
	u := user{username:name + strconv.Itoa(UC._users_count), userid:address, online:true}
	UC._users[UC._users_count] = u
	UC._users_count++
	return u
}

type userCollection struct {
	_users [10]user //@todo use pointers for users?
	_users_count int
}

//var broadcast = make(chan string, 2)

func RunServer() {
	// create logger
	InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	// listen
	ln, errors := net.Listen("tcp", ":8080")
	// create the user collection
	uc := userCollection{_users_count:0}

	if errors != nil {
		panic(errors) // ???
	}
	ch := make(chan string, 2)
	go channelTest(ch)
	for {
		connection, error := ln.Accept()
		if errors != nil {
			fmt.Println(error)
		} else {
			// start chatting
			go doSomething(connection,uc, ch)
		}
	}
}

func channelTest(ch chan string) {
	for {
		fmt.Print(<- ch)
	}
}

func doSomething (connection net.Conn, uc userCollection, ch chan string) {
	u := uc.AddUser("User", connection.RemoteAddr().String())
	//@todo prompt for or assign username
	welcomeMessage := "Hello, " + u.username + ".\n"
	connection.Write([]byte(welcomeMessage))
	messages := 0
	for {
		input := make([]byte, 256)
		//@todo convert to use channels? each user will send to a universal receiving channel that will broadcast the messages
		connection.Read(input)
		input_processed, is_command := processInput(input)
		if(!is_command) {
//			fmt.Printf("%v: %s", u.username, input_processed) // @todo broadcast to channel
			text := fmt.Sprintf("%v: %s", u.username, input_processed)
			ch <- text
		} else {
			// user has entered a command
			if(input_processed == "/quit") {
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

func processInput (input []byte) (string, bool) {
	output := string(bytes.Trim(input, "\x00"))
	return output, strings.HasPrefix(output, "/")
}


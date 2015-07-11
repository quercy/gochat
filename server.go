// @todo implement logger
// @todo robusticise users

package main

import(
	"net"
	"fmt"
	"strconv"
	"bytes"
	"strings"
)

type User struct {
	username string
	userid string
	online bool
}

func  (UC UserCollection)  AddUser (name string, address string) User {
	// log: adding a user
	//@todo use a list instead
	u := User{username:name + strconv.Itoa(UC._users_count), userid:address, online:true}
	UC._users[UC._users_count] = u
	UC._users_count++
	return u
}

type UserCollection struct {
	_users [10]User
	_users_count int
}

func main() {
	// listen
	ln, errors := net.Listen("tcp", ":8080")
	// create the user collection
	uc := UserCollection{_users_count:0}
	if errors != nil {
		panic(errors)
	}
	for {
		connection, error := ln.Accept()
		if errors != nil {
			fmt.Println(error)
		} else {
			// start chatting
			go doSomething(connection,uc)
		}
	}
}

func doSomething (connection net.Conn, uc UserCollection) {
	u := uc.AddUser("User", connection.RemoteAddr().String())
	//@todo prompt for or assign username
	welcomeMessage := "Hello, " + u.username + ".\n"
	connection.Write([]byte(welcomeMessage))
	messages := 0
	for {
		input := make([]byte, 256)
		//@todo convert to use channels? each user will send to a universal receiving channel that will broadcast the messages
		connection.Read(input)
		input_processed, is_command := ProcessInput(input)
		if(!is_command) {
			fmt.Printf("%v: %s", u.username, input_processed)
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

func ProcessInput (input []byte) (string, bool) {
	output := string(bytes.Trim(input, "\x00"))
	return output, strings.HasPrefix(output, "/")
}

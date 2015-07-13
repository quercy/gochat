package gochat

import(
	"net"
	"fmt"
	"os"
	"bufio"
	"strings"
)

func init() {
	// create logger
	InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	Trace.Println("Logger initialized")
}

func RunServer() {
	Trace.Println("Listening on port 8080")
	ln, errors := net.Listen("tcp", ":8080")
	if errors != nil {
		Error.Println("Errors while listening: ",errors)
		os.Exit(1)
	}
	Trace.Println("Creating userCollection with _users:make([]*user, 0)")
	uc := userCollection{_users:make([]*user, 0)}
//	ch := make(chan []byte, 5)
//	go broadcast(ch, uc)
	for {
		connection, err := ln.Accept()
		if err != nil {
			Error.Println("Errors while connecting: ",error)
			os.Exit(1)
		} else {
			// start chatting
			Info.Println("New connection established")
			go handleConnection (connection, &uc)
		}
	}
}


type userCollection struct {
	_users []*user
}

type user struct {
	username string
	user_network_location string
	online bool
	connection net.Conn
}

// Adds a new user to the specified user collection - returns that user
func  (UC *userCollection)  AddUser (name string, address string, cn net.Conn) *user {
	u := user{username:name, user_network_location:address, online:true, connection:cn}
	Trace.Println("Adding user ", u, " to userCollection")
	UC._users = append(UC._users, &u)
	return &u
}

func handleConnection (connection net.Conn, uc *userCollection) {

	var usr *user;

	initial_loop: for {
		connection.Write([]byte("Enter your username: "))
		input, _ := readInput(connection)
		usr = uc.AddUser(input, connection.RemoteAddr().String(), connection)
		welcomeMessage := fmt.Sprintf("%v has entered the chat.", input)
		uc.broadcast(welcomeMessage, "")
		//	ch <- []byte(welcomeMessage)
		break initial_loop
	}

	connection_loop: for {
		if !usr.online {
			break connection_loop
		}

		input, is_command := readInput(connection)

		if(!is_command) {
			text := fmt.Sprintf("%v: %s", usr.username, input)
			uc.broadcast(text, usr.username)
//			ch <- []byte(text) // send to the broadcast stream
		} else {
			command := strings.Split(input, " ")
			switch command[0] {
				case "/quit":
					connection.Write([]byte("Goodbye.\n"))
					uc.broadcast("User " + usr.username + " has left the channel.", usr.username)
					defer connection.Close()
					break connection_loop
				case "/list":
					connection.Write([]byte(uc.getUserListString()))
				case "/kick":
					if len(command) > 1 {
						result := uc.kickUser(command[1])
						if result {
							connection.Write([]byte("You kicked " + command[1] + ".\n"))
						} else {
							connection.Write([]byte("No such user is online.\n"))
						}
					}
			}
		}
	}
}

// Broadcasts a message to all connected users : takes the message as a string and the username of the sender
func (uc *userCollection) broadcast (msg string, usr string) {
	for i := 0; i < len(uc._users); i++ {
		if uc._users[i].username != usr { // when called, no need to broadcast to the sending user (usually)
			uc._users[i].connection.Write([]byte(msg + "\n")) // getting there
		}
	}
}

func readInput(cn net.Conn) (string, bool) {
	buf := bufio.NewReader(cn)
	output, _ := buf.ReadString('\n')
	output = strings.TrimSpace(output)
	return output, strings.HasPrefix(output,"/")
}

func (uc *userCollection) getUserListString() string {
	output := "Username :: Online :: Network Location"
	for i := 0; i < len(uc._users); i++ {
		line := fmt.Sprintf("%v :: %v :: %v ", uc._users[i].username, uc._users[i].online, uc._users[i].user_network_location)
		output = strings.Join([]string{output,line}, "\n")
	}
	return strings.Join([]string{output,""}, "\n")
}

func (uc *userCollection) kickUser (userString string) bool {
	for _, u := range uc._users {
		if u.username == userString && u.online == true {
			u.connection.Write([]byte("You have been kicked.\n"))
			u.online = false
			u.connection.Close()
			return true;
		}
	}
	return false;
}
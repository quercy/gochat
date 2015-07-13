//package main
//import (
//	"net"
//	"bufio"
//	"os"
//	"strings"
//	"bytes"
//	"fmt"
//)
//
//func main() {
//	connection, error := net.Dial("tcp", "localhost:8080")
//	if error != nil {
//		panic("") // oh no!
//	}
//	input := bufio.NewReader(os.Stdin)
////	fmt.Println("Enter your username: ")
//	username, _ := input.ReadString('\n') // throw out the error
//	username = strings.Trim(username,"\n")
//	connection.Write([]byte(username)) // tell the server the username
//	channel := make(chan string)
//	buf := new(bytes.Buffer)
//	go do(channel)
//	for {
//		buf.ReadFrom(connection)
//		channel <- buf.String()
//	}
//}
//
//func do (channe chan string) {
//	fmt.Println(channe)
//}
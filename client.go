package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
)

/**
 * Listen for server responses
 */
func serverRes(conn net.Conn) {
	buff := make([]byte, 32);
	for {
		n, err := conn.Read(buff);
		if err != nil {
			fmt.Println("Error reading from connection:", err);
			continue;
		}
		msg := string(buff[:n]);
		fmt.Println(msg);
	}
}

/**
 * Attempt to send a message to the server
 */
func sendMsg(conn net.Conn, msg string) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
}

func readMsg(conn net.Conn, buff []byte) string {
	n, err := conn.Read(buff)
	if err != nil {
		return ""
	}
	return string(buff[:n])
}

func login(conn net.Conn) {
	for {
		fmt.Println("Prompt user name:")
		var uname string
		fmt.Scanln(&uname)

		sendMsg(conn, "l" + uname);
		
		buff := make([]byte, 128)
		res := readMsg(conn, buff)
		if res != "" {
			if res[0:1] == "1" {
				fmt.Println("Error: Invalid user name or this user name already exists")
			} else {
				break
			}
		} else {
			fmt.Println("Error reading response")
		}
	}
}

func connect() {
	conn, err := net.Dial("tcp", "localhost:4200")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1);
	} else {
		fmt.Println("Connection established with the server")
	}
	defer conn.Close();

	login(conn);
	go serverRes(conn);

	var msg string;
	scanner := bufio.NewScanner(os.Stdin);
	for {
		if scanner.Scan() {
			msg = scanner.Text();
			sendMsg(conn, msg);
		}
	}
}

func main() {
	connect();
}
package main

import (
	"fmt"
	"net"
	"strings"
)

/**
 * Reverse the content of a string
 */
func reverse(str string) (tmp string) {
	for _, v := range str {
		tmp = string(v) + tmp;
	}
	return;
}

/**
 * Listen for server responses, then reverse the order of character of the 
 * message received and send it back
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
		fmt.Println("msg:", msg);
		at := strings.Index(msg, "@");
		colon := strings.Index(msg, ":");
		uname := msg[at:colon];
		fmt.Println("uname:", uname);
		msg = `\msg ` + uname + ` ` + reverse(msg[colon + 1:]);
		fmt.Println("msg:", msg);
		sendMsg(conn, msg);
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
		fmt.Println("Prompt bot name:")
		var bname string;
		fmt.Scanln(&bname);

		sendMsg(conn, "b" + bname);
		
		buff := make([]byte, 128);
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
		return;
	} else {
		fmt.Println("Connection established with the server");
	}
	defer conn.Close();

	login(conn);
	serverRes(conn);
}

func main() {
	connect();
}
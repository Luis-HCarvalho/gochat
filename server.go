package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const CHANNEL_LEN = 2;

var channels = make(map[string](chan string));

/**
 * evaluate if there is already a channel vinculated to a given user name
 */
func evalUName(uname string) bool {
	if channels[uname] != nil {
		return false
	}
	return true
}

/**
 * filter array based on the return of the predicate function
 */
func filter(arr []string, predicate func(string) bool) (tmp []string) {
	for _, el := range arr {
		if predicate(el) {
			tmp = append(tmp, el);
		}
	}
	return;
}

/**
 * Concat all values of an array of string into a single string, with the 
 * values separated by a single white space char
 */
func concat(arr []string) (s string) {
	for _, v := range arr {
		s += v + " ";
	}
	return;
}

/**
 * Awaits for client side responses 
 */
func clientRes(conn net.Conn, read chan string) {
	buff := make([]byte, 32);
	n, err := conn.Read(buff);
	if err != nil {
		fmt.Println("Error reading from connection:", err);
		return;
	}
	msg := string(buff[:n]);
	read <- msg;
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var uname string;
	isbot := false;
	buff := make([]byte, 32);
	for {		// 
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		msg := string(buff[:n])
		fmt.Println("Received:", msg)

		if msg[0:1] == "l" && evalUName(msg[1:]) {
			channels[msg[1:]] = make(chan string, CHANNEL_LEN)
			uname = msg[1:]
			fmt.Println("UNAME:", uname)
			conn.Write([]byte("0"))
			break
		} else if msg[0:1] == "b" && evalUName(msg[1:]) {
			isbot = true;
			channels[msg[1:]] = make(chan string, CHANNEL_LEN)
			uname = msg[1:];
			fmt.Println("BNAME:", uname)
			conn.Write([]byte("0"))
			break
		} else {
			conn.Write([]byte("1"))
		}
	}
	channels["broadcast"] <- "@" + uname + " arrive!"

	read := make(chan string);
	go clientRes(conn, read);

	for {
		select {
		case msgPriv := <- channels[uname]:
			conn.Write([]byte(msgPriv))
		case msgBroad := <- channels["broadcast"]:
			if isbot == false {
				conn.Write([]byte(msgBroad));
			}
		case msgClient := <- read:
			splitMsg := strings.Split(msgClient, " ");
			splitMsg = filter(splitMsg, func(s string) bool { return len(s) > 0 })

			if splitMsg[0] == `\msg` {
				if strings.HasPrefix(splitMsg[1], `@`) {
					if evalUName(splitMsg[1][1:]) == false {
						msg := "[direct] @" + uname + ":" + concat(splitMsg[2:]);
						channels[splitMsg[1][1:]] <- msg;
					}
				} else {
					msg := "[all] @" + uname + ":" + concat(splitMsg[1:]);
					channels["broadcast"] <- msg;
				}
			} else if splitMsg[0] == `\changenick` {
				if len(splitMsg) > 1 && evalUName(splitMsg[1]) {
					channels[splitMsg[1]] = make(chan string, CHANNEL_LEN);
					close(channels[uname]);
					delete(channels, uname);
					channels["broadcast"] <- "User @" + uname + " now is: @" + splitMsg[1];
					uname = splitMsg[1];
				}
			} else if splitMsg[0] == `\exit` {
				close(channels[uname]);
				delete(channels, uname);
				channels["broadcast"] <- "User @" + uname + " logout";
				return;
			} else {
				conn.Write([]byte(`Invalid command`));
			}
			go clientRes(conn, read);
		}
	}
}

func listen() {
	listener, err := net.Listen("tcp", "localhost:4200")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Server is listening on localhost:4200")

	// init broadcast channel
	channels["broadcast"] = make(chan string, CHANNEL_LEN);

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn) // Handle each connection in a new goroutine
	}
}

func main() {
	listen();
}
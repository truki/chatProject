package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
)

type user struct {
	name string
	connection net.Conn
}

type users []user

func main() {
	// creating listener
	listener, err := net.Listen("tcp", ":8887")
	if err != nil {
		log.Fatalf("unable to start server: %s", err)
	}
	defer listener.Close()

	log.Printf("Chat server started and listening on :8887")

	var chatUsers users

	for {
		// listening
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err)
			continue
		}

		// Readig join message
		data, err := bufio.NewReader(conn).ReadString('\n')
		command := strings.Split(data, "|")[0]
		name := strings.Split(data, "|")[1]
		name = strings.TrimSuffix(name, "\n")
		log.Println("User connected --> ", name)
		if command == "join" {
			io.WriteString(conn, "server| #### Wellcome to GDG Marbella TCP Chat Server, Enjoy! ####" + "\n")
			chatUsers = append(chatUsers, user{name, conn})
			userList := ""
			for _, u := range chatUsers {
				if userList == "" {
					userList = u.name
				} else {
					userList = userList + ":" + u.name
				}
			}
			log.Println(userList)
			// Sendind to all users the new user (we send the whole list of users)
			for _, u := range chatUsers {
				io.WriteString(u.connection, "newuser|server|"+name+"|"+userList+"\n")
				log.Printf("Notified new user: %s", name)
			}
		}



		go func() {
			for {
				data, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					log.Println(err)
					continue
				}
				command := strings.Split(data, "|")[0]

				if command == "message" {
					user := strings.Split(data, "|")[1]
					msg := strings.Split(data,"|")[2]
					for _, u := range chatUsers {
						if u.connection.RemoteAddr().String() != conn.RemoteAddr().String() {
							io.WriteString(u.connection, command+"|"+user+"|"+msg+"\n")
							log.Printf("Message %s sent to %s from %s", command+"|"+user+"|"+msg, u.name, user)
						}
					}
				}

				if command == "privmessage" {
					user := strings.Split(data, "|")[1]
					userTo := strings.Split(data,"|")[2]
					msg := strings.Split(data,"|")[3]
					msg = strings.Join(msg, " ")
					msg = msg[1:]
					msg = strings.Join(msg, " ")
					for _, u := range chatUsers {
						if u.name == userTo {
							io.WriteString(u.connection, command+"|"+user+"|"+userTo+"|"+msg+"\n")
							log.Printf("Private Message %s sent to %s from %s", command+"|"+user+"|"+userTo+"|"+msg, u.name, user)
						}
					}
				}

			}
		}()
	}
}

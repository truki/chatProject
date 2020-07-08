package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/marcusolsson/tui-go"
)

func main() {

	usersConnected := make(map[string]bool)

	my_name := os.Args[1]
	usersConnected[my_name] = true
	conn, err := net.Dial("tcp", "localhost:8887")
	if err != nil {
		log.Println("Error connecting to chat")
		os.Exit(1)
	}
	defer conn.Close()

	// sending my name to server
	_, err = io.WriteString(conn,"join|"+my_name+"\n")
	if err != nil {
		log.Println("Error sending username to chat server")
		os.Exit(1)
	}
	// right side bar
	sidebar := tui.NewVBox(
		tui.NewLabel("USERS"),
		tui.NewLabel(my_name),
		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	// Main box where you can see all users messages
	history := tui.NewVBox()

	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err)
	}

	wellcome := strings.Split(data,"|")[1]

	history.Append(tui.NewHBox(
		tui.NewLabel(time.Now().Format("15:04")),
		tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", "server"))),
		tui.NewLabel(wellcome),
		tui.NewSpacer(),
	))


	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)



	// listening messages submitted
	input.OnSubmit(func(e *tui.Entry) {
		_, err := io.WriteString(conn, "message|"+my_name+"|"+e.Text()+"\n")
		if err != nil {
			log.Println("Error sending username to chat server")
			os.Exit(1)
		}
		history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("15:04")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", my_name))),
			tui.NewLabel(e.Text()),
			tui.NewSpacer(),
		))
		input.SetText("")
	})

	root := tui.NewHBox(sidebar, chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	go func() {
		for {
			data, _ := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				log.Println(err)
				continue
			}



			ui.Update(func(){
				command := strings.Split(data, "|")[0]
				if command == "message" {
					user := strings.Split(data, "|")[1]
					msg := strings.Split(data, "|")[2]
					msg = strings.TrimSuffix(msg,"\n")

					history.Append(tui.NewHBox(
						tui.NewLabel(time.Now().Format("15:04")),
						tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", user))),
						tui.NewLabel(msg),
						tui.NewSpacer(),
					))
				}

				if command == "newuser" {
					user := strings.Split(data, "|")[1]
					newuser := strings.Split(data, "|")[2]
					msg := strings.Split(data, "|")[3]
					msg = strings.TrimSuffix(msg,"\n")
					userList := strings.Split(msg, ":")
					if newuser != my_name {
						history.Append(tui.NewHBox(
							tui.NewLabel(time.Now().Format("15:04")),
							tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", user))),
							tui.NewLabel("Hi " + my_name + ", " + newuser + " has joined GDG Marbella chat."),
							tui.NewSpacer(),
						))
					}

					for _, u := range userList {
						_, found := usersConnected[u]
						if found == false {
							sidebar.Insert(sidebar.Length()-1, tui.NewLabel(u))
							usersConnected[u] = true
						}
					}
				}

			})
		}
	}()

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}

}

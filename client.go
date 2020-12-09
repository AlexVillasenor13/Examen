package main

import (
	"bufio"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

var name string
var Port int
var Chat []string

type ClientServer struct {
	Chat []string
}

func client() {
	c, err := rpc.Dial("tcp", "127.0.0.1:1306")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Nickname: ")
	fmt.Scanln(&name)
	err = c.Call("Server.Connect", name, &Port)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Connected Successfully c:")
	}
	go server()

	var op int64
	for {
		fmt.Println("1) Send message")
		fmt.Println("2) Send file")
		fmt.Println("3) Show chat")
		fmt.Println("0) Exit")
		fmt.Scanln(&op)

		switch op {
		case 1:
			var result string
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Message: ")
			scanner.Scan()
			mssg := scanner.Text()
			err = c.Call("Server.SendMssg", []string{name, mssg}, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
		case 2:
			//files
		case 3:
			printChat()
		case 0:
			var result string
			err = c.Call("Server.Exit", name, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
				return
			}
		}
	}
}

func server() {
	new_server := new(ClientServer)
	rpc.Register(new_server)
	port_str := ":" + strconv.Itoa(Port)
	ln, err := net.Listen("tcp", port_str)
	if err != nil {
		fmt.Println(err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}

func (this *ClientServer) SetChat(chat []string, reply *string) error {
	this.Chat = chat
	Chat = chat
	this.printChat()
	return nil
}
func (this *ClientServer) printChat() {
	fmt.Println("----------------Chat------------------")
	for _, mssg := range this.Chat {
		fmt.Println(mssg)
	}
}

func printChat() {
	fmt.Println("----------------Chat------------------")
	for _, mssg := range Chat {
		fmt.Println(mssg)
	}
}

func main() {
	client()
}

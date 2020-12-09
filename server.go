package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"strconv"
)

type Server struct {
	Users     []string
	UsersChat []*[]string
	UsersPort []int
	Files     []string
	Chat      []string
}

var port int

func (this *Server) findUser(name string) int {
	for i, n := range this.Users {
		if name == n {
			return i
		}
	}
	return -1
}

func (this *Server) Connect(name string, new_port *int) error {
	if name != "" {
		port += 1
		this.Users = append(this.Users, name)
		var new_chat []string
		this.UsersChat = append(this.UsersChat, &new_chat)
		this.UsersPort = append(this.UsersPort, port)
		*new_port = port
		str := name + " Connected c:"
		fmt.Println(str)
		return nil
	} else {
		str := name + "NOT CONNECTED!!!"
		return errors.New(str)
	}
}

func (this *Server) Exit(name string, reply *string) error {
	i := this.findUser(name)

	if i != -1 {
		copy(this.Users[i:], this.Users[i+1:])
		this.Users[len(this.Users)-1] = ""
		this.Users = this.Users[:len(this.Users)-1]
		copy(this.UsersChat[i:], this.UsersChat[i+1:])
		this.UsersChat[len(this.UsersChat)-1] = nil
		this.UsersChat = this.UsersChat[:len(this.UsersChat)-1]
		*reply = name + " Disconnected :c"
		fmt.Println(*reply)
		fmt.Println(this.Users)
		fmt.Println(this.UsersChat)
		*reply = "Disconnected Successfully :c"
		return nil
	} else {
		str := name + "NOT DISCONNECTED!!!"
		return errors.New(str)
	}
}

func (this *Server) SendMssg(data []string, reply *string) error {
	if data[1] != "" {
		this.addChat(data)
		this.printChat()
		*reply = "Message received"
		return nil
	} else {
		return errors.New("Message NOT received!!")
	}
}
func (this *Server) addChat(data []string) {
	this.Chat = append(this.Chat, data[0]+": "+data[1])
	i_user := this.findUser(data[0])

	for i, chat := range this.UsersChat {
		if i == i_user {
			*chat = append(*chat, "You: "+data[1])
		} else {
			*chat = append(*chat, data[0]+": "+data[1])
			port_str := "127.0.0.1:" + strconv.Itoa(this.UsersPort[i])
			c, err := rpc.Dial("tcp", port_str)
			if err != nil {
				fmt.Println(err)
				return
			}
			result := ""
			err = c.Call("ClientServer.SetChat", *chat, &result)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}

func (this *Server) printChat() {
	fmt.Println("----------------Chat------------------")
	for _, mssg := range this.Chat {
		fmt.Println(mssg)
	}
	fmt.Println("----------------Chats------------------")
	for _, chat := range this.UsersChat {
		fmt.Println("----------------------------------")
		for _, mssg := range *chat {
			fmt.Println(mssg)
		}
	}
}

func server() {
	new_server := new(Server)
	rpc.Register(new_server)
	port_str := ":" + strconv.Itoa(port)
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

func main() {
	port = 1306
	go server()
	var input string
	fmt.Scanln(&input)
}

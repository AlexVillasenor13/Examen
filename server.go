package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

type Server struct {
	Users      []string
	UsersChat  []*[]string
	UsersPort  []int
	NamesFiles []string
	Files      [][]byte
	Chat       []string
}

var port int
var local_chat []string
var local_files []string

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
func (this *Server) SendFile(data []byte, reply *string) error {
	if string(data) != "" {
		this.NamesFiles = append(this.NamesFiles, *reply)
		this.Files = append(this.Files, data)
		local_files = this.NamesFiles
		*reply = "File received"
		return nil
	} else {
		return errors.New("File NOT received!!")
	}
}

func (this *Server) addChat(data []string) {
	this.Chat = append(this.Chat, data[0]+": "+data[1])
	local_chat = this.Chat
	i_user := this.findUser(data[0])

	for i, chat := range this.UsersChat {
		if i == i_user {
			*chat = append(*chat, "You: "+data[1])
		} else {
			*chat = append(*chat, data[0]+": "+data[1])
		}
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
		c.Close()

	}
}

func (this *Server) printChat() {
	fmt.Println("----------------Chat------------------")
	for _, mssg := range this.Chat {
		fmt.Println(mssg)
	}
}
func printLocalChat() {
	fmt.Println("----------------Chat------------------")
	for _, mssg := range local_chat {
		fmt.Println(mssg)
	}
}
func printLocalFiles() {
	fmt.Println("----------------Files------------------")
	for _, mssg := range local_files {
		fmt.Println(mssg)
	}
}

func writeFile(content, file_name string) {
	file, err := os.Create(file_name + ".txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	file.WriteString(content)
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
	var op int64
	for {
		fmt.Println("1) Show Chat")
		fmt.Println("2) Show Files")
		fmt.Println("3) Save Chat")
		fmt.Println("4) Save Files")
		fmt.Println("0) Exit")
		fmt.Scanln(&op)

		switch op {
		case 1:
			printLocalChat()
		case 2:
			printLocalFiles()
		case 3:
			writeFile(strings.Join(local_chat, "\n"), "Chat")
			fmt.Println("Chat saved")
		case 4:
			writeFile(strings.Join(local_files, "\n"), "Files")
			fmt.Println("Files saved")
		case 0:
			return
		}
	}
}

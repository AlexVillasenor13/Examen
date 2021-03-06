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
			var result string
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Name: ")
			scanner.Scan()
			file_name := scanner.Text()
			fmt.Print("Location: ")
			scanner.Scan()
			loc := scanner.Text()
			file, err := os.Open(loc + file_name)

			if err != nil {
				fmt.Println(err)
				return
			}

			defer file.Close()

			stat, err := file.Stat()
			if err != nil {
				fmt.Println(err)
				return
			}

			total := stat.Size()

			bs := make([]byte, total)
			count, err := file.Read(bs)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(bs, "bytes:", count)
			err = c.Call("Server.SendFile", [][]byte{[]byte(loc + file_name), bs, []byte(name)}, &result)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
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
func (this *ClientServer) SetFile(data [][]byte, reply *string) error {
	*reply = "File received"
	fmt.Print("File Received from ")
	fmt.Println(string(data[2]))
	fmt.Println(data[1])
	fmt.Println(string(data[1]))
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

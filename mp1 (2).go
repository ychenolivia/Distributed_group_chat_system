package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var login_name string
var reader = bufio.NewReader(os.Stdin)
var client_map = make(map[string]*net.TCPConn)
var pmap = make(map[int]string)
var pmap1 = make(map[int]string)
var pmap2 = make(map[int]*net.TCPConn)
var mmap = make(map[string]string)
var nmap = make(map[string]string)

func listen_client(ip_port string, num int) { //listen to the port
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ip_port)
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Selesct another port. This port has been used.")
		os.Exit(1)
	}
	flag := true
	for {
		client_con, _ := tcpListener.AcceptTCP()
		kv := strings.Split(client_con.RemoteAddr().String(), ":")
		for key := range client_map {
			kk := strings.Split(key, ":")
			if kv[0] == kk[0] {
				flag = false
			}
		}
		if flag {
			client_map[client_con.RemoteAddr().String()] = client_con
			go add_receiver(client_con)
		}

	}

}

func add_receiver(current_connect *net.TCPConn) { //receive massage from connect
	for {
		byte_msg := make([]byte, 2048)
		len, err := current_connect.Read(byte_msg) //read messsage from connect
		if err != nil {
			return
		}
		remsg := string(byte_msg[:len])
		name := strings.Split(remsg, ":")
		ki := strings.Split(current_connect.RemoteAddr().String(), ":")
		nmap[ki[0]] = name[0]
		msg_broadcast(byte_msg[:len], current_connect.RemoteAddr().String()) //immediately broadcast message when receving message ensure R-multicast
	}
}

func msg_broadcast(byte_msg []byte, key string) { //send message to all exist connect
	for _, con := range client_map {
		_, err := con.Write(byte_msg)
		if err != nil {
			fmt.Println("Write error")
		}
	}
}
func msg_receiver(self_connect *net.TCPConn) { // receive message from broadcast
	buff := make([]byte, 2048)
	for {
		len, err := self_connect.Read(buff)
		if err != nil { //check if someone left
			for key := range pmap2 {
				if pmap2[key] == self_connect {
					fmt.Println(nmap[pmap1[key]] + "left")
				}
			}
			return
		}
		remsg := string(buff[:len])
		sremsg := strings.Split(remsg, "/")
		flag := true
		for key := range mmap { //check if the message have been received, the repetitive message will be ignored
			if key == sremsg[0] && mmap[key] == sremsg[1] {
				flag = false
				break
			}
		}
		if flag {
			mmap[sremsg[0]] = sremsg[1]
			fmt.Println(sremsg[0])
		}

	}
}

func msg_sender() { //send message
	for {
		read_line_msg, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("input error")
			continue
		}
		t := time.Now().Nanosecond() //add a timestamp for message
		ct := strconv.Itoa(t)
		read_line_msg = []byte(login_name + " : " + string(read_line_msg) + "/" + ct)
		for key := range pmap2 {
			_, err2 := pmap2[key].Write(read_line_msg)
			if err2 != nil {
				fmt.Println("Write error")
			}
		}

	}
}
func connect(num int, port string) { //Constantly connect to other machines
	for {
		for key := range pmap {
			addr := pmap[key] + ":" + port
			tcp_addr, _ := net.ResolveTCPAddr("tcp", addr)
			con, err := net.DialTCP("tcp", nil, tcp_addr)
			if err != nil {
				continue
			} else {
				pmap2[key] = con
				pmap[key] = ""

			}

		}
		if len(pmap2) == num { //check if all of the people are here
			fmt.Println("Ready")
			break
		}

	}
	go msg_sender()
	for key := range pmap2 {
		go msg_receiver(pmap2[key])
	}

}

func main() { //run this file as form:$go run mp1.go name port num
	pmap[1] = "172.22.94.28"
	pmap[2] = "172.22.156.29"
	pmap[3] = "172.22.158.29"
	pmap[4] = "172.22.94.29"
	pmap[5] = "172.22.156.30"
	pmap[6] = "172.22.158.30"
	pmap[7] = "172.22.94.30"
	pmap[8] = "172.22.156.31"
	pmap[9] = "172.22.158.31"
	pmap[10] = "172.22.94.31"
	mmap["1"] = "1"
	for key := range pmap {
		pmap1[key] = pmap[key]
	}
	if len(os.Args) == 4 {
		login_name = os.Args[1]
		Port := os.Args[2]
		var port string = Port
		n := os.Args[3]
		num, _ := strconv.Atoi(n)
		if num < 2 || num > 8 {
			fmt.Println("Only 2-8 people are allowed")
			os.Exit(1)
		}
		myip, _ := net.InterfaceAddrs()
		myIp := myip[1].String()
		ip := strings.Split(myIp, "/")
		myipaddr := ip[0] + ":" + port
		go listen_client(myipaddr, num)
		time.Sleep(1 * time.Second)
		go connect(num, port)
	} else {
		fmt.Println("Please enter the correct information.")
	}
	select {}
}

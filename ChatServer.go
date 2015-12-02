// ChatServer
package ChatServer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type Client struct {
	AccountID int
	Name      string
	Icon      int
	Level     int
	Area      int
	Viplevel  int
	RecMsg    chan<- []byte
	AckMsg    <-chan []byte
	//ChatChannlID  int
	Client_Socket net.Conn
	UserRelation
}

type UserRelation struct {
	FriendInfo map[uint64]string
	IgnorID    map[uint64]string
}

type ChatServer struct {
	ID       int                  //服务器ID
	Port     int                  //监听Client端口号
	Host     string               //host地址
	UserList map[net.Conn]*Client //用户列表
	//ChannlList       map[int]chan byte    //通道列表
	//ChannlList_Count []int                //统计Channlist_count里面各个Channl的用户数
	//DataChannl       chan byte            //中转接受数据
	ChatConfigData map[string]string
}

func (this *ChatServer) ProcessConf(args []string) {
	if len(args) != 2 {
		return
	} else {
		this.ChatConfigData[args[0]] = args[1]
	}
	fmt.Println(args, this.ChatConfigData)
}

func (this *ChatServer) ReadConf(name string, handler func([]string)) (err error) {
	this.InitChatServerData()
	f, err := os.Open(name)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		buf := strings.Split(line, " = ")

		handler(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func (this *ChatServer) InitChatServerData() {
	this.ID = 0
	this.Port = 0
	this.Host = ""
	//this.DataChannl = make(chan byte, 1024)
	this.ChatConfigData = make(map[string]string, 0)
	this.UserList = make(map[net.Conn]*Client, 100)
	//this.ChannlList = make(map[int]chan byte, 0)
}

func (this *ChatServer) InitServer() error {
	serverhost := this.ChatConfigData["ChatHost"] + ":" + this.ChatConfigData["ChatPort"]
	server, err := net.Listen("TCP4", serverhost)
	defer server.Close()

	if err != nil {
		return err
	}

	//监听连接
	for {
		client_socket, err := server.Accept()
		if err != nil {
			return err
		}
		this.NewClient(client_socket)
	}

	return nil
}

func (this *ChatServer) NewClient(n net.Conn) {
	client := &Client{Client_Socket: n}
	client.RecMsg = make(chan []byte, 1024)
	client.AckMsg = make(chan []byte, 1024)
	this.UserList[n] = client

	//创建接受发送线程

	go func() {
		for {
			var buf []byte
			n, err := client.Client_Socket.Read(buf)
			defer client.Client_Socket.Close()

			if err != nil {
				if err != io.EOF {
					fmt.Println("Read Msg:", err)
					return
				}
			}

			if n != 0 {
				client.RecMsg <- buf
				client.ClientMsgProcess()
			}
		}
	}()

	go func() {
		for {
			buf := <-client.AckMsg
			n, err := client.Client_Socket.Write(buf)
			defer client.Client_Socket.Close()

			if err != nil {
				if err != io.EOF {
					fmt.Println("Write Msg:", err)
					return
				}
			}

			if n != 0 {
			}
		}
	}()

}

func (this *Client) ClientMsgProcess() {
	select {
	case <-this.RecMsg:
		{

		}
	default:
		{

		}
	}
}

func (this *ChatServer) InitDB(addr string, host string, DBname string) {}

func (this *ChatServer) ServerMsgProcess() {}

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
	AccountID     int
	Name          string
	Icon          int
	Level         int
	Area          int
	Viplevel      int
	IsLive        bool
	RecMsg        chan string
	AckMsg        chan string
	ChatServer    *ChatServer
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
	DataChannl     chan byte //中转接受数据
	ChatConfigData map[string]string
	RecDataSize    uint64
	AckDataSize    uint64
}

func (this *ChatServer) processConf(args []string) {
	if len(args) != 2 {
		return
	} else {
		this.ChatConfigData[args[0]] = args[1]
	}
	fmt.Println(args, this.ChatConfigData)
}

func (this *ChatServer) ReadConf(name string) (err error) {
	this.initChatServerData()
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

		this.processConf(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func (this *ChatServer) initChatServerData() {
	this.ID = 0
	this.Port = 0
	this.Host = ""
	this.DataChannl = make(chan byte, 1024)
	this.ChatConfigData = make(map[string]string, 0)
	this.UserList = make(map[net.Conn]*Client, 100)
	//this.ChannlList = make(map[int]chan byte, 0)
}

func (this *ChatServer) InitServer() error {
	serverhost := ":" + this.ChatConfigData["ChatPort"]

	server, err := net.Listen("tcp4", serverhost)

	defer server.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}
	//return nil
	//监听连接
	for {
		fmt.Println("Wait Accept.")
		client_socket, err := server.Accept()
		if err != nil {
			return err
		}
		this.newClient(client_socket)
		fmt.Println("Accept success")
	}
	fmt.Println("End......")

	return nil
}

func (this *ChatServer) newClient(n net.Conn) {
	client := &Client{IsLive: true, ChatServer: this, Client_Socket: n}
	client.RecMsg = make(chan string, 1024)
	client.AckMsg = make(chan string, 1024)

	this.UserList[n] = client
	fmt.Println(client)
	fmt.Println("UserList: ", len(this.UserList))
	//return
	//创建接受发送线程

	go func() {
		var buf []byte
		for {
			fmt.Println("Socket Waiting Read.")
			buf = make([]byte, 1024)
			n, err := client.Client_Socket.Read(buf)
			defer client.Client_Socket.Close()

			if err != nil {
				if err != io.EOF {
					fmt.Println("Read Msg:", err)
					client.IsLive = false
					break
				}
			}
			str := string(buf)
			fmt.Println(str)
			fmt.Println("Rec len: ", n)
			if n != 0 {
				client.RecMsg <- str
				client.ChatServer.RecDataSize += uint64(n)
				client.ClientMsgProcess()
			}
		}
	}()

	//	go func() {
	//		for {
	//			fmt.Println("ClientMsgProcess.")
	//			client.ClientMsgProcess()
	//			defer func() {
	//				close(client.AckMsg)
	//				close(client.RecMsg)
	//				client.Client_Socket.Close()
	//			}()
	//			if !client.IsLive {
	//				break
	//			}
	//		}
	//	}()

}

func (this *Client) ClientMsgProcess() {
	fmt.Println("In ClientMsgProcess...")
	select {
	case buf := <-this.RecMsg:
		{
			fmt.Println("In <-this.RecMsg...")
			//buf := <-this.RecMsg
			fmt.Println(buf)
			fmt.Println("Out <-this.RecMsg...")
		}
	case <-this.AckMsg:
		{
			buf := []byte(<-this.AckMsg)
			n, err := this.Client_Socket.Write(buf)
			if err != nil {
				if err != io.EOF {
					return
				}
			}

			if n != 0 {
				this.ChatServer.AckDataSize += uint64(n)
			}

		}
	}
}

func (this *Client) SendToClient(buf string) {
	this.AckMsg <- buf
	this.ClientMsgProcess()
}

func (this *ChatServer) InitDB(addr string, host string, DBname string) {}

func (this *ChatServer) ServerMsgProcess() {}

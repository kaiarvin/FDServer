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
	RecMsg        chan string
	ChatChannlID  int
	Client_Socket net.Conn
	UserRelation
}

type UserRelation struct {
	FriendInfo map[uint64]string
	IgnorID    map[uint64]string
}

type ChatServer struct {
	ID               int                 //服务器ID
	Port             int                 //监听Client端口号
	Host             string              //host地址
	UserList         map[net.Conn]Client //用户列表
	ChannlList       map[int]chan byte   //通道列表
	ChannlList_Count []int               //统计Channlist_count里面各个Channl的用户数
	DataChannl       chan byte           //中转接受数据
	ChatConfigData   map[string]string
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
	this.DataChannl = make(chan byte, 1024)
	this.ChatConfigData = make(map[string]string, 0)
	this.UserList = make(map[net.Conn]Client, 100)
	this.ChannlList = make(map[int]chan byte, 20)
}

func (this *ChatServer) InitServer() error {
	server, err := net.Listen("TCP4", "localhost:8085")
	defer server.Close()

	if err != nil {
		return err
	}

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
	this.UserList[n] = client

	//分配channl统计channl count

}

func (this *ChatServer) InitDB(addr string, host string, DBname string) {}

func (this *ChatServer) ClientMsgProcess() {}

func (this *ChatServer) ServerMsgProcess() {}

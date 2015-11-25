// ChatServer
package ChatServer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Client struct {
	AccountID    int
	Name         string
	Icon         int
	Level        int
	Area         int
	Viplevel     int
	RecMsg       chan string
	ChatChannlID int
	UserRelation
}

type UserRelation struct {
	FriendInfo map[uint64]string
	IgnorID    map[uint64]string
}

type ChatServer struct {
	ID             int               //服务器ID
	Port           int               //监听Client端口号
	Host           string            //host地址
	UserList       map[uint64]Client //用户列表
	ChannlList     map[int]chan byte //通道列表
	DataChannl     chan byte         //中转接受数据
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

func (this *ChatServer) initChatServerData() {
	this.ID = 0
	this.Port = 0
	this.Host = ""
	this.DataChannl = make(chan byte, 1024)
	this.ChatConfigData = make(map[string]string, 0)
	this.UserList = make(map[uint64]Client)
	this.ChannlList = make(map[int]chan byte)
}

func (this *ChatServer) InitServer() {

}

func (this *ChatServer) InitDB(addr string, host string, DBname string) {}

func (this *ChatServer) ClientMsgProcess() {}

func (this *ChatServer) ServerMsgProcess() {}

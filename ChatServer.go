// ChatServer
package ChatServer

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"qbs"
	"strings"
)

type ChatServer struct {
	ID             int                  //服务器ID
	Port           int                  //监听Client端口号
	Host           string               //host地址
	UserList       map[net.Conn]*Client //用户列表
	AccountList    map[uint64]net.Conn  //account 列表
	NameList       map[string]net.Conn  //name列表
	DataChannl     chan byte            //中转接受数据
	ChatConfigData map[string]string
	RecDataSize    uint64
	AckDataSize    uint64
	DB             *qbs.Qbs //character数据库
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

	DBType := this.ChatConfigData["ChatDBType"]
	DBUser := this.ChatConfigData["ChatDBUser"]
	DBPw := this.ChatConfigData["ChatDBPw"]
	DBHost := this.ChatConfigData["ChatDBHost"]
	DBPort := this.ChatConfigData["ChatDBPort"]
	DBName := this.ChatConfigData["ChatDBName"]
	DBParam := this.ChatConfigData["ChatDBParam"]
	dberr := this.InitQbs(DBType, DBUser, DBPw, DBHost, DBPort, DBName, DBParam)
	if dberr != nil {
		fmt.Println(dberr)
		return err
	}
	defer this.GetServerQbs().Close()

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

			if n != 0 {
				data := string(buf[:n])
				fmt.Println(data)
				fmt.Println("Rec len: ", n)

				client.RecMsg <- data
				client.ChatServer.RecDataSize += uint64(n)
				client.ClientMsgProcess()
			}
		}
	}()
}



func (this *ChatServer) InitQbs(dbtype, dbuser, pw, dbhost, dbport, dbname, param string) error {
	dsn := dbuser + ":" + pw + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?" + param
	fmt.Println(dsn)

	q, err := DBInit(dbtype, dsn, dbname)

	if nil != err {
		fmt.Println(err)
		return err
	}
	this.DB = q
	return nil
}

func (this *ChatServer) ServerMsgProcess() {}

func (this *ChatServer) GetServerQbs() *qbs.Qbs {
	return this.DB
}
